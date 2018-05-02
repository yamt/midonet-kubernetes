# MidoNet Kubernetes Integration

## How to build

### Prequisite

- docker
- dep

### Build docker images

<pre>
	% dep ensure
	% docker build -f Dockerfile -t yamt/midonet-kube-controllers .
	% docker build -f Dockerfile-node -t yamt/midonet-kube-node .
</pre>

## How to install

0. Install MidoNet as usual. Every Kubernetes nodes including the master
   node should run MidoNet agent.
1. Create a MidoNet logical router. Record its UUID for the later use.
2. "kubeadm init" with Node IPAM enabled.
<pre>
	% kubeadm init --pod-network-cidr=10.1.0.0/16
</pre>
3. Remove kube-proxy.
<pre>
	% kubectl -n kube-system delete ds kube-proxy
</pre>
4. Copy and edit midonet-kube-config.template.yaml to match your deployment.
   Use the above mentioned MidoNet router UUID here.
   The modified file will be called midonet-kube-config.yaml hereafter.
5. Apply manifests.
<pre>
	% kubectl apply -f midonet-kube-config.yaml
	% kubectl apply -f midonet-kube-controllers.yaml
	% kubectl apply -f midonet-kube-node.yaml
</pre>
6. "kubeadm join" as usual.

## Limitations

* MidoNet API authentication is not supported.  You should use Mock auth.
  [MNA-1273][MNA-1273]
* Only ClusterIP Service Type is implemented.
* Even if a Service has multiple Endpoints, only one endpoint which happens
  to be first is always used.  I.e. no load-balancing. [MNA-1264][MNA-1264]

[MNA-1273]: https://midonet.atlassian.net/browse/MNA-1273
[MNA-1264]: https://midonet.atlassian.net/browse/MNA-1264
