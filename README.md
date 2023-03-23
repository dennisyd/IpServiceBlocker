# IPBlockerService

This service hosts an HTTP-based API, as well as a gRPC API, to allow determining if a given IP address is within
a list of allowable countries.

This service uses the following technologies:
- Golang 1.20.2 found [here](https://go.dev/doc/install)
- Install gRPC prereqs found [here](https://grpc.io/docs/languages/go/quickstart/#prerequisites)
- Gin Web Framework found [here](https://github.com/gin-gonic/gin)
- GeoIP2 Reader for Go found [here](https://github.com/oschwald/geoip2-golang)
- Docker found [here](https://docs.docker.com/engine/install/)
- Kubernetes found [here](https://kubernetes.io/docs/setup/)

# Configuration
Configuration of the service is handled via environment variables read in the [base/config.go](base/config.go) file.

# Build, test, and deploy
To successfully run all unit tests and deploy, you'll need to ensure you have a copy of the GeoLite2 database. Instructions can be found [here](/data/README.md)

# General Design
API layer (HTTP and gRPC)
    ->  Handler layer
        -> Data layer

# Updating the database
To perform a database update, you will need to download the newest version of the GeoLite2 database (instructions found [here](/data/README.md)), then perform a `make docker.push`. This will package the database into the Docker file prior to pushing to Docker Hub. See the [# Considerations for Improvements](#considerations-for-improvements) section for a more ideal solution.

# Considerations for Improvements
At the moment, this service expects to house the GeoLite2 DB locally, per microservice. This is cumbersone on multiple fronts:
1. It means that DB updates need to be part of the actual deployment of the service
2. Each service has an increased memory footprint

A more ideal solution would be to use a cloud based storage system (e.g. AWS S3 or Google Cloud Storage) to house the actual
database file and then access said file over the internet. This would at least solve the first issue.

More considerations may be necessary to reduce the footprint as reading the binary database file will inevitably read it all
directly into memory. One possible option would be to make a whole extra service to simply manage the database. This is
assuming that the increased memory footprint (about 5.8 MB) is actually a concern.

Additional improvements:
- OpenAPI 3.0 server generation.
  - So far I haven't found a great Go module for this. There are some for OpenAPI 2.0 but the ones I've found for 3.0 all seem to have various annoyances/learning curves that I'm not prepared to delve into.
- Interface for the underlying database handling.
  - This isn't really necessary unless there's need of using other databases and/or using a GeoLite2 webservice but it could provide some flexibility.
- Security.
  - Both the HTTP and gRPC APIs are currently insecure. There's plenty of work to be done in this area, from TLS support on the HTTP API to credential management on the gRPC API.