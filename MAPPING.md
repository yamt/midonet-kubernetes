Resource mapping between Kubernetes and MidoNet
===============================================

| Kubernetes | MidoNet     |
|:-----------|:------------|
| Node       | Bridge      |
| Pod        | Bridge Port |

Kubernetes Node
---------------

When the controller notices a Kubernetes Node, it would create
the following MidoNet REST API objects.

- A Bridge
- A Bridge Port on the bridge
- A Router Port on the cluster router
- A Port Link to link the above two ports
- A local Route on the cluster router for the subnet (PodCIDR of the node)

Kubernetes Pod
--------------

When the controller notices a Kubernetes Pod, it would create
the following MidoNet REST API objects.

- A Bridge Port on the Node Bridge
