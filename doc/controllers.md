# Controllers

midonet-kube-controllers continuously watches changes in Kubernetes
Resources like Pods and Nodes, and makes the necessary changes
on the MidoNet API.

The executable contains several controllers.
You can choose which controllers to enable by the ENABLED_CONTROLLER
environment variable.  By default all controllers are enabled.

By design, those controllers are independent each other and can be
run in separate processes.  Such a setup is not extensively tested
though.

midonet-kube-controllers can run anywhere, as far as it has
L3 connectivity to the Kubernetes API server.
Some of its embedded controllers needs the connectivity to
the MidoNet API too. (See the following diagram.)

<pre>
...................         ...................................
:                 :         :                                 :
:  K8s API server :         :  midonet-kube-controllers       :
:                 :         :                                 :
:                 : Watch   :                                 :
:                 : Annotate:                                 :
:  +-----------+ <-----+    :    +-------------+              :
:  |Node       |  :    |    :    |nodeannotator|              :  Query
:  +-----------+ +--+  +-------> |controller   | +--------------------+
:                 : |       :    +-------------+              :       |
:  +-----------+  : |       :                                 :       |
:  |Pod        | +--+       :                                 :       v
:  +-----------+  : |       :   +--------------------+        :
:                 : |       :   |pod controller      |        :   +-----------+
:  +-----------+  : |       :   +--------------------+        :   |           |
:  |Service    | +--+ Watch :   |node controller     |        :   |MidoNet API|
:  +-----------+  : +---------> +--------------------+ +-+    :   |           |
:                 : |       :   |service controller  |   |    :   +-----------+
:  +-----------+  : |       :   +--------------------+   |    :
:  |Endpoints  | +--+       :   |endpoints controller|   |    :       ^
:  +-----------+  :         :   +--------------------+   |    :       |
:                 :  Update :                            |    :       |
:                 :   +----------------------------------+    :       |
:                 :   |     :                                 :       |
:  +-----------+ <----+     :                                 :       |
:  |Translation|  :         :   +------------------+          :       |
:  +-----------+ +------------> |pusher controller | +----------------+
:                 :  Watch  :   +------------------+          :  Update
:                 :         :                                 :
:                 :         :                                 :
...................         ...................................
</pre>

## pod, node, service, endpoints

These controllers watch the corresponding Kubernetes resources
and create/update/delete Translation custom resources accordingly.

## pusher

This controller watches Translation custom resources and
create/update/delete MidoNet API resources accordingly.

## nodeannotator

This controller adds "midonet.org/host-id" annotation to Kubernetes
Node resources, by querying MidoNet API with the assumption that
MidoNet Host name and Kubernetes Node name on a node match.

This controller also adds "midonet.org/tunnel-zone-id" and
"midonet.org/tunnel-endpoint-ip" annotations.

The annotation is used by pod and node controllers.
