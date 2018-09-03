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
