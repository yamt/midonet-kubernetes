# midonet-kube-node and CNI

<pre>
                    ............
                    :
                    :   Per Node components
                    :
+------------+      :   +-------------------+
|            |      :   |                   |
|            |      :   | midonet-kube-node |
|            |      :   |                   |
|            |      :   |        +----------|
|            |      :   |        | shared   |
| Kubernetes |      :   |        | plumbing +----------------------------+
| API server |      :   |        | code     | Node plumbing              |
|            |      :   |        +----------|                            |
|            |      :   |                   |                            |
|            |      :   | +-------------+   |                            |
|    Pod <----------------+ gRPC server |   |                            |
|            | Annotate | +-------------+   |                            |
|            |  IP  :   |       ^           |                            |
+------------+  MAC :   |       |           |                            |
                etc :   |       |           |                            |
                    :   +-------|-----------+                            |
                    :           |                                        |
 ....................           | gRPC call                              |
 :                              | over Unix domain socket                |
 :                              |                                        |
 :  +---------+                 |                                        |
    |         |          +------+-----------+         +------------+     |
    | kubelet +--------> |                  +-------> | host-local |     |
    |         | exec     | midonet-kube-cni | exec    | IPAM CNI   |     |
    +---------+          |                  |         +------------+     |
                         |       +----------|                            v
                         |       | shared   |            +-----------------+
                         |       | plumbing +----------> | kernel          |
                         |       | code     | Pod        |                 |
                         +------------------+ plumbing   |  routing table  |
                                                         |  veth pair      |
                                                         |                 |
                                                         +-----------------+
</pre>

## midonet-kube-cni

midonet-kube-cni is the CNI plugin for this integration.
It connects the Pod to the cluster network by setting up
a veth pair and IP routing.
It reports the generated MAC address for the Pod via local
midonet-kube-node instance.
Note: CNIs don't have API credentials for Kubernetes or MidoNet.

## midonet-kube-node

midone-kube-node connects the Node (Linux root netns of the host)
to the cluster network.

It also provides a gRPC service over a unix domain socket
for local midonet-kube-cni instances.

# Node connectivity

We connect Nodes to the cluster network in a similar way as Pods.
That is, to create a veth pair and connect one side to the network.
The following is a summary of differences between Pods and Nodes.

|              | Pod (CNI)                | Node                        |
|:-------------|:-------------------------|:----------------------------|
| IPAM         | host-local               | 2nd IP of podCIDR           |
| contNS       | args.Netns (from docker) | host namespace              |
| contVethName | args.IfName ("eth0")     | fixed ("midokube-node")     |
| hostVethName | generated from NS/Pod    | fixed ("midokube-mido")     |
| MidoNet port | pod.idForKey(podKey)     | node.portIDForKey(nodename) |
