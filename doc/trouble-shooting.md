## Trouble shooting

### Translation custom resources

midonet-kube-controllers use a CRD "Translation" to store the
intermediate objects converted from Kubernetes resources like
Pods and Services.
You can investigate them with kubectl command.
"tr" is a short name for the CRD.

#### Examples

- Get a list of Translation resources for a Service.
<pre>
k% kubectl -n kube-system get tr -l midonet.org/owner-uid=$(kubectl -n kube-system get svc kube-dns -o jsonpath='{.metadata.uid}')
NAME                                                                                                   AGE
service-port.3.kube-system-kube-dns-dns-10.96.0.10-17-53-a3982c897cb53bcf4ec5dd4895d5e191e28628a5      9d
service-port.3.kube-system-kube-dns-dns-tcp-10.96.0.10-6-53-4b9da0761bde315827a2f3a49560b52b3c7a53a5   9d
service.3.kube-dns                                                                                     9d
</pre>

- The same for an Endpoints.
<pre>
k% kubectl -n kube-system get tr -l midonet.org/owner-uid=$(kubectl -n kube-system get ep kube-dns -o jsonpath='{.metadata.uid}')
NAME                                                                                                      AGE
endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.21-53-udp-54b27350879608cd686eb61cedb7273386a367c2        11s
endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.22-53-udp-3adaf57ddc83531c51e037c7489c27ab6f671afa        10s
endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.23-53-udp-1277c2fba4ef5db3a8a406311d1c1a534b09e171        10s
endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.24-53-udp-9562d46816ee2c652d4d0f8a880489a13ec1e073        8s
endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.1.145-53-udp-8c74d17e9a85139c771b0eb8321c000c88907c94       1h
endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.1.177-53-udp-7afaaf3a2ba96d55bf2c9192a0cb4a902c8271f6       22sendpoints-port.3.kube-dns-dns-tcp-10.96.0.10-10.1.0.21-53-tcp-6245126dfb5cd4878ccce4b468dd229e73ee7f6c    10s
endpoints-port.3.kube-dns-dns-tcp-10.96.0.10-10.1.0.22-53-tcp-1400a542aaf5e799316f80cb6daa0dc141f0902f    9s
endpoints-port.3.kube-dns-dns-tcp-10.96.0.10-10.1.0.23-53-tcp-7ebfb0755608d0cefb2939d92b18edd654fcb1ca    9s
endpoints-port.3.kube-dns-dns-tcp-10.96.0.10-10.1.0.24-53-tcp-f624318b925446dfa7016c1d2e1adfb4917f3fa9    12sendpoints-port.3.kube-dns-dns-tcp-10.96.0.10-10.1.1.145-53-tcp-3302c499a8532147a976c74a96165bea98dc7034   1h
endpoints-port.3.kube-dns-dns-tcp-10.96.0.10-10.1.1.177-53-tcp-3ae2c8ecc12fb3c76ac73e6260b16bb35d8795ed   20s
</pre>

- Investigate one of those Translations.
<pre>
k% kubectl -n kube-system describe tr endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.24-53-udp-9562d46816ee2c652d4d0f8a880489a1
3ec1e073
Name:         endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.24-53-udp-9562d46816ee2c652d4d0f8a880489a13ec1e073
Namespace:    kube-system
Labels:       midonet.org/owner-uid=a55ca29d-2330-11e8-8d8e-fa163ec6ef35
Annotations:  <none>
API Version:  midonet.org/v1
Kind:         Translation
Metadata:
  Cluster Name:
  Creation Timestamp:  2018-06-28T00:44:03Z
  Finalizers:
    midonet.org/deleter  Generation:  0
  Owner References:
    API Version:     v1
    Kind:            Endpoints    Name:            kube-dns
    UID:             a55ca29d-2330-11e8-8d8e-fa163ec6ef35
  Resource Version:  12040595
  Self Link:         /apis/midonet.org/v1/namespaces/kube-system/translations/endpoints-port.3.kube-dns-dns-10.96.0.10-10.1.0.24
-53-udp-9562d46816ee2c652d4d0f8a880489a13ec1e073
  UID:               5adf0dce-7a6c-11e8-ba60-fa163ec6ef35
Resources:
  Body:    {"id":"5cfc9435-b3f4-5a79-8093-3c6a8c0f5fff","tenantId":"midonetkube","name":"KUBE-SEP-kube-dns/dns/10.96.0.10/10.1.0.24/53/UDP"}
  Kind:    Chain
  Parent:
  Body:    {"id":"e75a66ad-a7e6-569f-83bc-543d4a01107a","type":"jump","jumpChainId":"5cfc9435-b3f4-5a79-8093-3c6a8c0f5fff"}
  Kind:    Rule
  Parent:  2ba08097-dcb1-548e-817e-d259736a286c
  Body:    {"id":"0f0dbbff-40d9-537a-ae8f-3fdc533cb9fd","type":"dnat","flowAction":"accept","natTargets":[{"addressFrom":"10.1.1.177","addressTo":"10.1.1.177","portFrom":53,"portTo":53}]}
  Kind:    Rule
  Parent:  5cfc9435-b3f4-5a79-8093-3c6a8c0f5fff
  Body:    {"id":"2399be50-d5bd-5691-a6ff-8ffbc14df8a8","type":"snat","dlType":2048,"nwSrcAddress":"10.1.1.177","nwSrcLength":32,"flowAction":"continue","natTargets":[{"addressFrom":"10.96.0.10","addressTo":"10.96.0.10","portFrom":30000,"portTo":60000}]}
  Kind:    Rule
  Parent:  5cfc9435-b3f4-5a79-8093-3c6a8c0f5fff
Events:
  Type    Reason                   Age               From                      Message
  ----    ------                   ----              ----                      -------
  Normal  TranslationUpdatePushed  52s (x2 over 1m)  midonet-kube-controllers  Translation Update pushed to the backend
</pre>

### Prometheus

midonet-kube-controllers provides a few metrics for Prometheus.

<pre>
k% curl -s http://localhost:9453/metrics|grep -E "^# (HELP|TYPE) midonet_"
# HELP midonet_kube_controllers_midonet_client_request_duration_seconds Latency of MidoNet API call
# TYPE midonet_kube_controllers_midonet_client_request_duration_seconds histogram
# HELP midonet_kube_controllers_midonet_client_requests_total Number of MidoNet API calls
# TYPE midonet_kube_controllers_midonet_client_requests_total counter
</pre>

### Go net/http/pprof

midonet-kube-controllers provides Go net/http/pprof on
the same port as the Prometheus metrics.

<pre>
k% go tool pprof http://localhost:9453/debug/pprof/heap
</pre>

#### Examples queries

See [Prometheus Querying documentation][prometheus-query] for details.

- Number of successful MidoNet API calls per seconds.
<pre>
sum(rate(midonet_kube_controllers_midonet_client_requests_total{code=~"2.*"}[5m])) by (resource,method)
</pre>

- Latency of successful MidoNet API calls.
<pre>
histogram_quantile(0.9, sum(rate(midonet_kube_controllers_midonet_client_request_duration_seconds_bucket{code=~"2.*"}[5m])) by (resource,method,le))
</pre>

[prometheus-query]: https://prometheus.io/docs/prometheus/latest/querying/basics/
