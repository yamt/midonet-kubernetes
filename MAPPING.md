Resource mapping between Kubernetes and MidoNet
===============================================

| Kubernetes | MidoNet     |
|:-----------|:------------|
| Node       | Bridge      |
| Pod        | Bridge Port |
| Service    | Chain/Rules |
| Endpoint   | Chain/Rules |

Prerequisite MidoNet resources for a deployemnt
-----------------------------------------------

For a given Kubernetes deployment:

- A Router (we call this the cluster router)
- "SERVICES" Chain
- REV_SNAT Rule
- REV_DNAT Rule

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

Kubernetes Service
------------------

- A Chain for the Service
- In the global "SERVICES" Chain which is shared by all Node bridges:
	- Rules to redirect the service traffic to the above per-Service Chain

Kubernetes Endpoint
------------------

- A Chain for the Endpoint
- In the Service Chain:
	- A Jump rule to the Endpoint Chain
	  (In future, we might implement probability match for this rule)
- In the Endpoint Chain:
	- A Rule to SNAT if the source IP matches the Endpoint IP
	- A Rule to DNAT to the endpoint IP
