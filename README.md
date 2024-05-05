# Margo Prototype Project

This repository is to be used only for research purposes and not intended to be a reference implementation for the [Margo specification](https://github.com/margo/specification).

The current implementation only works with the example `app-description.yaml` and `deviceDescription.json` located in the `./install/examples` folder.

The current implementation only works with installing the hello-world application one time on the edge device. Updating and deleting the application, or installing other applications, is not implemented yet.

## Prerequisites

The following prerequisites are need on the host machine if you want to work with the source code.

- [VS Code](https://code.visualstudio.com/Download)
- [Docker with Docker Compose](https://www.docker.com/)

### VSCode Dev Containers

This project uses the [VS Code](https://code.visualstudio.com/Download) [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) to run the development environment.

When running the project using dev containers a `docker-app-builder` container is also started. This container must be used to build the docker images and helm charts because it has the helm and docker software installed.

When running the project using dev containers a `gogs` container is also started. This container runs are [GOGs server](https://gogs.io/) and is used for the git repos for this prototype.

## Source (./src folder)

The source code is under the `src` folder

### /src/apps

This folder contains the simple hello-world application, dockerfile and helm chart.

The `build-images.sh` script can be used to build the docker image for the application

### /src/WOS

This folder contains the services that run on the workload orchestration solution. This folder contains the source code, dockerfiles and helm chart for these services.

The `build-images.sh` script can be used to build the docker images for all these services.

#### /src/WOS/chart

This is the helm chart that packages up all of the workload orchestration related services so they can be deployed to kubernetes.

#### /src/WOS/gitops_pullservice

This is the source code for the `gitops_pullservice.` This service pulls the application description file from the GOGs repository.

#### /src/WOS/gitops_pushservice

This is the source code for the `gitops_pushservice.` This service pushes the files needed to be able to deploy the application to the edge device to the git repository for the device's desired state.

#### /src/WOS/orchestration_service

This is the source code for the `orchestration_service.` This is the main service that contains the APIs used by the other services.

#### /src/WOS/orchestration_portal

This is the source code for the `orchestration_portal.` This is a web service that provides a basic UI for interacting with the orchestration service to add application registries and install the application on the device.

### /src/WOSA

This folder contains the services that run on the workload orchestration agent on the edge device. This folder contains the source code, dockerfiles and helm charts for these services.

the `build-images.sh` script can be used to build the docker image for the `gitops_client`

#### /src/WOSA/chart

This is the helm chart that packages up all of the workload orchestration agent related services so they can be deployed to kubernetes.

#### /src/WOSA/gitops_client

This is the source code for the `gitops_client.` This service pulls down the desired state and adds the solution and instance CRDs to the kubernetes cluster. Right now this service only works to install a new instance of the hello-world application.

#### /src/WOS/orchestration-operator

This is the source code for the `orchestration-operator.` This is using the [Kubernetes Operator SDK](https://sdk.operatorframework.io/). This is a kubernetes operator that will use the information in the solution.yaml and instance.yaml CRDs to install the helm chart indicated by the files.

Use the `MAKE` file to run the operator and create the docker image.

The current implementation is based on [Symphony's](https://github.com/eclipse-symphony/symphony) instance and solution object model.

## Installation (./install folder)

The [install readme.md](./install/readme.md) file has information about how to setup the environment and deploy the helm charts.
