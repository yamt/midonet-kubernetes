## Translation Custom Resource

A Translation is a Kubernetes Custom Resource, which represents
an ordered list of backend resources for the correspoinding
Kubernetes resource.

It uses the ownership mechanism to simplify the deletion handling.

It has a finanizer to ensure the deletion of the corresponding resources
on the backend.

<pre>
Pod -------- Translation
              owner: Pod
              finalizers:
                deleter.midonet.org
              resources:
                Port
                HostInterfacePort
</pre>

## Intermediate Custom Resources

We want to make the number of backend resources for a kubernetes
resource fixed to keep the backend interaction simple.
However, for some of complex Kubernetes resources like Service and
Endpoints, it isn't straightforward. To mitigate the complexity,
we introduce intermediate Custom Resources. ("servicePort" in
the below figure)

<pre>
Service ----+--- servicePort 1 --- Translation
            |     owner: Service    owner: servicePort 1
            |
            +--- servicePort 2 --- Translation
            |     owner: Service    owner: servicePort 2
            |
            +--------------------- Translation
                                    owner: Service
</pre>

### concern
consider updating a Service twice.
the first update adds a servicePort and the second update deletes it.
when the controller processes the second update, it's possible that
its informer have not seen the servicePort addtion from the first
update yet. in that case, it might fail to delete the servicePort.
is it a real problem? if so, what can we do?
is there a way to wait for our own updates being propagated?
