package midonet

import (
	"crypto/rand"
	"net"

	log "github.com/sirupsen/logrus"
)

func RandomMac() net.HardwareAddr {
	// AC-CA-BA  Midokura Co., Ltd.
	addr := [6]byte{ 0xac, 0xca, 0xba, }
	_, err := rand.Read(addr[3:])
	if err != nil {
		log.Fatal("rand")
	}
	return net.HardwareAddr(addr[:])
}
