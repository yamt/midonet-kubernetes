# Controllers

midonet-kube-controllers executable contains several controllers.

You can choose which controllers to enable by the ENABLED_CONTROLLER
environment variable.  By default all controllers are enabled.

## pod, node, service, endpoints

These controllers watch the corresponding Kubernetes resources
and create/update/delete Translation custom resources accordingly.

## pusher

This controller watches Translation custom resources and
create/update/delete MidoNet API resources accordingly.

## hostresolver

This controller adds "midonet.org/host-id" annotation to Kubernetes
Node resources, by querying MidoNet API with the assumption that
MidoNet Host name and Kubernetes Node name on a node match.

The annotation is used by pod and node controllers.
