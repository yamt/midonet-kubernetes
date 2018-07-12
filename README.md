[![Go Report Card](https://goreportcard.com/badge/github.com/midonet/midonet-kubernetes)](https://goreportcard.com/report/github.com/midonet/midonet-kubernetes)
[![GoDoc](https://godoc.org/github.com/midonet/midonet-kubernetes?status.svg)](https://godoc.org/github.com/midonet/midonet-kubernetes)

# MidoNet Kubernetes Integration

## Overview

This software provides a way to use [MidoNet][MidoNet] as a backend
for Kubernetes networking.  Namely, it provides the following
Kubernetes networking functionalitites.

* Basic cluster network, that is, connectivity among Pods, Nodes, and the apiserver
* Services with ClusterIP type (Note: externalIPs are ignored)

[MidoNet]: https://github.com/midonet/midonet

### Supported versions

This software is tested with:

* Kubernetes v1.10
* MidoNet 5.6

### Limitations

* Even if a Service has multiple Endpoints, only one endpoint which happens
  to be first is always used.  I.e. no load-balancing. [MNA-1264][MNA-1264]

[MNA-1264]: https://midonet.atlassian.net/browse/MNA-1264

### References

* The [doc][doc] directry contains internal documentations

* The [design doc][design] might have more details

[doc]: ./doc
[design]: https://docs.google.com/document/d/1dYwz26I6NXO0MnbUf_pnC2Ihoz1Kdp0Pdm0DmEmGn4I/edit

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## License

Apache License 2.0.  See [LICENSE](./LICENSE)
