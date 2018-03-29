Resource mapping between Kubernetes and MidoNet
===============================================

| Kubernetes | MidoNet     |
|:-----------|:------------|
| Node       | Bridge      |
| Pod        | Bridge Port |

Kubernetes Node
---------------

For a Kubernetes Node, the controller would create the following MidoNet
virtual devices.

- A Bridge
- A Bridge Port
- A Router Port on the cluster router
- A local Route on the cluster router for the subnet (PodCIDR of the node)

Kubernetes Pod
--------------

For a Kubernetes Pod, the controller would create the following MidoNet
virtual devices.

- A Bridge Port on the Node Bridge
