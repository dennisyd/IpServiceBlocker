# IPBlockerService

This service hosts an HTTP-based API to allow determining if a given IP address is within
a list of allowable countries.

# Dev setup
// TODO

## Prereqs
- Install protoc
// TODO

# Configuration
Configuration of the service is handled via the [base/config.go](base/config.go) file.

# Build, test, and deploy
To successfully run all unit tests and deploy, you'll need to ensure you have a copy of the GeoLite2 database. Instructions can be found [here](/data/README.md)

# General Design

API layer (HTTP or gRPC)
    ->  Handler layer
        -> Data layer


# Considerations for Improvements
At the moment, this service expects to house the GeoLite2 DB locally, per microservice. This is cumbersone on multiple fronts:
1. It means that DB updates need to be part of the actual deployment of the service
2. Each service has a vastly increased memory footprint

A more ideal solution would be to use a cloud based storage system (e.g. AWS S3 or Google Cloud Storage) to house the actual
database file and then access said file over the internet. This would at least solve the first issue.

More considerations may be necessary to reduce the footprint as reading the binary database file will inevitably read it all
directly into memory. One possible option would be to make a whole extra service to simply manage the database. This is
assuming that the increased memory footprint (about 5.8 MB) is actually a concern.