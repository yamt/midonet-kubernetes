Resource mapping between Kubernetes and MidoNet
===============================================

| Kubernetes | MidoNet     |
|:-----------|:------------|
| Node       | Bridge      |
| Pod        | Bridge Port |

Kubernetes Node
---------------

For a Kubernetes Node, the controller would create the following MidoNet
REST API objects.

- A Bridge
- A Bridge Port on the bridge
- A Router Port on the cluster router
- A Port Link to link the above two ports
- A local Route on the cluster router for the subnet (PodCIDR of the node)

Kubernetes Pod
--------------

For a Kubernetes Pod, the controller would create the following MidoNet
REST API objects.

- A Bridge Port on the Node Bridge
