package utils

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/types/current"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/yamt/midonet-kubernetes/pkg/cni/types"
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
)

// DoNetworking performs the networking for the given config and IPAM result
func DoNetworking(args *skel.CmdArgs, conf types.NetConf, result *current.Result, logger *logrus.Entry, desiredVethName string) (hostVethName, contVethMAC string, err error) {
	// Select the first 11 characters of the containerID for the host veth.
	hostVethName = "cali" + args.ContainerID[:Min(11, len(args.ContainerID))]
	contVethName := args.IfName
	var hasIPv4, hasIPv6 bool

	// If a desired veth name was passed in, use that instead.
	if desiredVethName != "" {
		hostVethName = desiredVethName
	}

	logger.Infof("Setting the host side veth name to %s", hostVethName)

	// Clean up if hostVeth exists.
	if oldHostVeth, err := netlink.LinkByName(hostVethName); err == nil {
		if err = netlink.LinkDel(oldHostVeth); err != nil {
			return "", "", fmt.Errorf("failed to delete old hostVeth %v: %v", hostVethName, err)
		}
		logger.Infof("Cleaning old hostVeth: %v", hostVethName)
	}

	err = ns.WithNetNSPath(args.Netns, func(hostNS ns.NetNS) error {
		veth := &netlink.Veth{
			LinkAttrs: netlink.LinkAttrs{
				Name:  contVethName,
				Flags: net.FlagUp,
				MTU:   conf.MTU,
			},
			PeerName: hostVethName,
		}

		if err := netlink.LinkAdd(veth); err != nil {
			logger.Errorf("Error adding veth %+v: %s", veth, err)
			return err
		}

		hostVeth, err := netlink.LinkByName(hostVethName)
		if err != nil {
			err = fmt.Errorf("failed to lookup %q: %v", hostVethName, err)
			return err
		}

		if mac, err := net.ParseMAC("EE:EE:EE:EE:EE:EE"); err != nil {
			logger.Infof("failed to parse MAC Address: %v. Using kernel generated MAC.", err)
		} else {
			// Set the MAC address on the host side interface so the kernel does not
			// have to generate a persistent address which fails some times.
			if err = netlink.LinkSetHardwareAddr(hostVeth, mac); err != nil {
				logger.Warnf("failed to Set MAC of %q: %v. Using kernel generated MAC.", hostVethName, err)
			}
		}

		// Explicitly set the veth to UP state, because netlink doesn't always do that on all the platforms with net.FlagUp.
		// veth won't get a link local address unless it's set to UP state.
		if err = netlink.LinkSetUp(hostVeth); err != nil {
			return fmt.Errorf("failed to set %q up: %v", hostVethName, err)
		}

		contVeth, err := netlink.LinkByName(contVethName)
		if err != nil {
			err = fmt.Errorf("failed to lookup %q: %v", contVethName, err)
			return err
		}

		// Fetch the MAC from the container Veth. This is needed by Calico.
		contVethMAC = contVeth.Attrs().HardwareAddr.String()
		logger.WithField("MAC", contVethMAC).Debug("Found MAC for container veth")

		// At this point, the virtual ethernet pair has been created, and both ends have the right names.
		// Both ends of the veth are still in the container's network namespace.

		for _, addr := range result.IPs {

			// Before returning, create the routes inside the namespace, first for IPv4 then IPv6.
			if addr.Version == "4" {
				// Add a connected route to a dummy next hop so that a default route can be set
				gw := net.IPv4(169, 254, 1, 1)
				gwNet := &net.IPNet{IP: gw, Mask: net.CIDRMask(32, 32)}
				err := netlink.RouteAdd(
					&netlink.Route{
						LinkIndex: contVeth.Attrs().Index,
						Scope:     netlink.SCOPE_LINK,
						Dst:       gwNet,
					},
				)

				if err != nil {
					return fmt.Errorf("failed to add route inside the container: %v", err)
				}

				if err = ip.AddDefaultRoute(gw, contVeth); err != nil {
					return fmt.Errorf("failed to add the default route inside the container: %v", err)
				}

				if err = netlink.AddrAdd(contVeth, &netlink.Addr{IPNet: &addr.Address}); err != nil {
					return fmt.Errorf("failed to add IP addr to %q: %v", contVethName, err)
				}
				// Set hasIPv4 to true so sysctls for IPv4 can be programmed when the host side of
				// the veth finishes moving to the host namespace.
				hasIPv4 = true
			}
		}

		if err = configureContainerSysctls(hasIPv4, hasIPv6); err != nil {
			return fmt.Errorf("error configuring sysctls for the container netns, error: %s", err)
		}

		// Now that the everything has been successfully set up in the container, move the "host" end of the
		// veth into the host namespace.
		if err = netlink.LinkSetNsFd(hostVeth, int(hostNS.Fd())); err != nil {
			return fmt.Errorf("failed to move veth to host netns: %v", err)
		}

		return nil
	})

	if err != nil {
		logger.Errorf("Error creating veth: %s", err)
		return "", "", err
	}

	// Moving a veth between namespaces always leaves it in the "DOWN" state. Set it back to "UP" now that we're
	// back in the host namespace.
	hostVeth, err := netlink.LinkByName(hostVethName)
	if err != nil {
		return "", "", fmt.Errorf("failed to lookup %q: %v", hostVethName, err)
	}

	if err = netlink.LinkSetUp(hostVeth); err != nil {
		return "", "", fmt.Errorf("failed to set %q up: %v", hostVethName, err)
	}

	return hostVethName, contVethMAC, err
}

// configureContainerSysctls configures necessary sysctls required inside the container netns.
func configureContainerSysctls(hasIPv4, hasIPv6 bool) error {
	var err error

	// Globally disable IP forwarding of packets inside the container netns.
	// Generally, we don't expect containers to be routing anything.

	if hasIPv4 {
		if err = writeProcSys("/proc/sys/net/ipv4/ip_forward", "0"); err != nil {
			return err
		}
	}

	if hasIPv6 {
		if err = writeProcSys("/proc/sys/net/ipv6/conf/all/forwarding", "0"); err != nil {
			return err
		}
	}

	return nil
}

// writeProcSys takes the sysctl path and a string value to set i.e. "0" or "1" and sets the sysctl.
func writeProcSys(path, value string) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	n, err := f.Write([]byte(value))
	if err == nil && n < len(value) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}