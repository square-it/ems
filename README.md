# TIBCO EMS Go client

This repository contains the source code for the TIBCO EMS Go client library.

## Prerequisites
This client is designed to work with the EMS 8.3 client libraries as shipped with TIBCO EMS.

## Prepare

1. Pull a Docker image of TIBCO EMS from a private Docker registry:
```
export TIBCO_EMS_DOCKER_IMAGE=<docker-registry>:5000/tibco/tibco-ems:8.3.0
docker pull $TIBCO_EMS_DOCKER_IMAGE
```

2. Run a TIBCO EMS container
```
docker run --name ems -d -p 7222:7222 $TIBCO_EMS_DOCKER_IMAGE
```

3. Extract TIBCO EMS directory
```
export TIBCO_EMS_DIRECTORY=/tmp/ems
mkdir $TIBCO_EMS_DIRECTORY
docker export ems | tar -xf - -C $TIBCO_EMS_DIRECTORY
docker rm -f ems
```

## Build

1. Export the cgo CFLAGS and LDFLAGS directives to the correct location of your local EMS client libraries:
```
export CGO_CFLAGS="-I. -I$TIBCO_EMS_DIRECTORY/opt/tibco/ems/8.3/include/tibems"
export CGO_LDFLAGS="-L$TIBCO_EMS_DIRECTORY/opt/tibco/ems/8.3/lib -ltibems64"             
```

2. Build the library
```
git clone https://github.com/square-it/ems.git
cd ems

go build .
```

## Test

1. Run a TIBCO EMS container
```
docker run --name ems -d -p 7222:7222 $TIBCO_EMS_DOCKER_IMAGE
```

2. Launch Go tests
```
export LD_LIBRARY_PATH=$TIBCO_EMS_DIRECTORY/opt/tibco/ems/8.3/lib

go test .
```

## Reporting bugs

Please report bugs by raising issues for this project in github https://github.com/square-it/ems/issues
