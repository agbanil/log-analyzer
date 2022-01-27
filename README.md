# Log Analyzer

what do we have here?

- [Objective](#Objective)
- [Installation and Running Application](#Installation-and-Running-Application)
    - [Prerequisites](#Prerequisites)
    - [Clone and Build Project](#Clone-and-Build-Project)
    - [Running the application with `go run`](#Running-the-application-with-go-run)
- [Docker](#Docker)
- [Running Tests](#Running-Tests)
- [Deployment](#Deployment)
  - [Helm](#Helm)
- [Improvements](#Improvements)

## Objective

Create a solution (a script or a program) that takes the contents of a log file, and outputs the devices and their
classification.

This program runs as a sort of daemon listening on files, analyzing the inputs and producing outputs via stdout.
It works by grabbing comma delimited file paths from an env variable (FILE_PATHS) and listening on and analyzing those files.

At the moment, the outputs have the file names as the key in order to tell what files the outputs are for.


## Installation and Running Application
### Prerequisites

Go 1.17 or higher. To check, run `go version`:

```bash
$ go version
go version go1.17.6 darwin/amd64
```

### Clone and Build Project

```bash
# Clone Repository
git clone https://github.com/raznerdeveloper/log-analyzer.git
```

### Running the application with `go run`

```bash
# Runs the app with the file inside the project
FILE_PATHS=file.txt go run main.go

or

make run-local
```

Output:

```bash
$ go run main.go
{"file.txt":{"temp-1":"precise"}}
{"file.txt":{"temp-1":"precise","temp-2":"ultra precise"}}
{"file.txt":{"hum-1":"keep","temp-1":"precise","temp-2":"ultra precise"}}
```

## Docker

The application was dockerized using this [Dockerfile](./Dockerfile) but there's a makefile provided to make things easy.
To run the app using docker (with the provided makefile helper), you can do:
```bash
# Build the Docker Image
make build

# Run the Docker Image - by default it uses the test file in the project in the repository
make run

# You can also run the docker image (with our darling makefile) using a file location as an argument
# This file location has to be an absolute path on your machine. E.g., /Users/Alice/Golang/awesomeProject/file.txt
make run-file file=file.txt
```

## Running Tests

To run the unit tests, execute the following:

```bash
go test -v

# You can also use our wonderful makefile for this by doing:
make test
```

## Deployment
### Helm
The application was packaged using helm (with files in the repo). In order to deploy this to your kubernetes cluster, do:
```bash
helm install gatsbytakehome ./k8s
```
Manual steps would be to:
- Build the docker image
- Push to your favourite container registry
- Reference the CR inside of `k8s/values.yaml`
- do helm install

However, I have already built and pushed an image to docker hub and referenced it so, the helm install step should
work with what's available.

There is an infrastructure config for this in [this repo](https://github.com/raznerdeveloper/log-infra.git). Kindly go there and
follow the instructions in order to setup a kubernetes infrastructure on AWS before running the `helm install` if you don't 
already have a kubernetes cluster to run this on.

## Improvements

- Currently, new reading results are printed when there's a new line with a new device. The improvement here would be to print
all the values as they're calculated.
- Improve on how the results are being displayed. Currently, since it's tailing the file, it displays for every device as it gets
the results