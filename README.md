Poll
====

This experimental repo gives an example of how the web server usage model can be used to shape
its infrastructure definition and configuration.

This repo is a Go module. Packages `votes` with `cmd/pollsvc` provide the web service implementation that 
exposes REST API for collecting votes that represent if audience agrees with a speaker during a
presentation.

This repo is also a [CUE](https://cuelang.org) module. Package `infra` contains the definition of the 
service usage model expressed in CUE (see `infra/model`), as well as the deployment code (see `infra/deployment`)
configured with the parameters derived from the usage model.
