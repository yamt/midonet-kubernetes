## How to deploy

The following procedure assumes [kubeadm][kubeadm] based deployment.
But something similar should apply to other deployment methods as well.

[kubeadm]: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/

0. Deploy MidoNet as usual.

   * Every Kubernetes nodes including the master node should run MidoNet agent.
   * The following instruction assumes that MidoNet Host names and
     Kubernetes Node names are same for each nodes. It's usually the case
     because both of them are inferred from the hostname.
   * No need to create a tunnel zone.  The integration will create its own
	 automatically.
   * Some of these assumptions can be overridden with manual [annotations][annotations].

1. "kubeadm init" with Node IPAM enabled.
   (This integration relies on Node's spec.PodCIDR.)
<pre>
	% kubeadm init --pod-network-cidr=10.1.0.0/16
</pre>
2. Remove kube-proxy.
   (It isn't necessary or compatible with this integration.
   Unfortunately, kubeadm unconditionally sets it up.
   cf. [kubeadm issue 776][kubeadm-776])
<pre>
	% kubectl -n kube-system delete ds kube-proxy
</pre>
3. After stopping kube-proxy, you might need to remove iptables rules
   installed by kube-proxy manually.
   Note: the following commands would remove many of relevant rules but
   leave some of rules and chains installed by kube-proxy. The simplest
   way to get a more clean state is to reboot the system.
<pre>
	% sudo iptables -t nat -F KUBE-SERVICES
	% sudo iptables -F KUBE-SERVICES
</pre>
4. Look at [manifests][manifests] directory in this repository.
   Copy and edit midonet-kube-config.template.yaml to match your deployment.
   The modified file will be called midonet-kube-config.yaml hereafter.
5. Apply manifests.
<pre>
	% kubectl apply -f midonet-kube-config.yaml
	% kubectl apply -f midonet-kube.yaml
</pre>
6. "Untaint" the master node if you want.
7. If you have workers, do "kubeadm join" as usual.

[annotations]: ./doc/labels-annotations.md#annotations
[kubeadm-776]: https://github.com/kubernetes/kubeadm/issues/776
[manifests]: ./manifests
