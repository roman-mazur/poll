Poll
====

This experimental repo gives an example of how the web server usage model can be transformed into
an infrastructure definition making it possible to deploy the service with a cloud provider utilizing the
resources derived from the usage model.

This repo is a Go module. Package `votes` with `cmd/svc` provide the web service implementation that 
exposes REST API for collecting votes that represent if audience agrees with a speaker during a
presentation.

This repo is also a [CUE](https://cuelang.org) module. Package `defs` contains the definition of the 
service usage model expressed in CUE, as well as its transformation into a set of specific AWS cloud 
resources that can be instantiated with Terraform in order to deploy the web service. 
**Still work in progress**

Running `go generate ./...` is necessary before working with the data in the repo.

Some examples
-------------
(execute from the `defs` directory)

Check samples of the stored data.
```shell
cue eval -e summary.samples
```

Get memory requirement for the web service.
```shell
cue eval -e summary.req 
```
