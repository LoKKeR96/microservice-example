# Example code for DDD microservice architecture in Go
This repository gives an overview of how you could structure a containerised RESTful and/or GraphQL-Application written in Go following a DDD and CQRS approach. It uses gqlgen to generate necessary GraphQL server files, Echo as the HTTP server and framework, gorm for ORM data handling and mapping to DB. Two container orchestration configurations are implemented (Docker Compose and Kubernetes) which provide more flexibility in case scalability is an important requirement.
A CI approach is also implemented making use of GitHub Actions to run tests when a merge into branch `main` is done and when a pull request is created.

A complete list of the Tools and Frameworks used can be found later in this README.


## Contents
1. [Intentions](#intentions)
2. [Getting Started](#getting-started)
    * [Noteworthy things that should be considered while developing](#noteworthy-things-that-should-be-considered-while-developing)
3. [Golang Frameworks and Tools](#golang-frameworks-and-tools)
4. [Other tools](#other-tools)
5. [Core Concepts](#core-concepts)
    * [Domain Driven Design Architecture](#domain-driven-design-architecture)
    * [Echo](#echo)
    * [Application](#application)
    * [Domain](#domain)
    * [Persistence](#persistence)
6. [Docker Compose](#docker-compose)
7. [Kubernetes](#kubernetes)
8. [CI/CD Pipeline](#ci/cd-pipeline)
9. [TODOs](#todos)


## Intentions
* Providing an example project that can be used as boilerplate for starting up a project
* Demonstrating the used core-concepts for well-structuring a Go application
* Demonstrating proficiency with GraphQL and RESTful API protocols
* Demonstrating proficiency with Docker and Kubernetes orchestration tools
* Learning from mutual feedback and improve this approach


## Getting started
* Clone or download this project

* Generate the GraphQL-Server and Resolver-Interfaces by running `go run github.com/99designs/gqlgen generate` in your project directory

* Build the server image using `docker compose build --no-cache server'

* Starting up the stack using docker compose or kubernetes

* Visit `localhost:8080/playground` for getting started with GraphQL API

* Send HTTP requests to `localhost:8080` based on the path of the endpoint you want to use


### Noteworthy things that should be considered while developing
**Schema-Changes lead you to re-generate**
* Every change in the `schema.graphqls` should be finished with re-generating the GraphQL-Files by running `go run github.com/99designs/gqlgen generate`.

**New Actions are not automatically added to the resolvers-implementation**
* New Actions inside the schema do not lead gqlgen to add them to your resolvers. You are responsible for implement them on your own (The interface containing all necessary functions can be found in `src/infrastructure/graph/server_generated.go`).


## Golang Frameworks and Tools
|Framework/Tool|Description|
|---|---|
|[gorm.io/gorm](https://pkg.go.dev/gorm.io/gorm)|ORM library used to map and manage entities in the database|
|[99design/gqlgen](https://github.com/99designs/gqlgen)|Generating GraphQL-Server and Resolver-Interface based on `schema.graphql`. Can also generate missing models|
|[palantir/stacktrace](https://github.com/palantir/stacktrace)|Used to manage and trace errors and log them on the console|
|[gofrs/uuid](https://github.com/gofrs/uuid)|Used for generation UUIDs for the in-memory database used in this example|
|[dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go)|Used for generation of a JWT toket for user authentication
|[labstack/echo/v4](https://github.com/labstack/echo/v4)|HTTP server and framework to serve APIs, manage middleware and respond to HTTP requests|
|[go-playground/validator/v10](https://github.com/go-playground/validator/v10)|Used to validate data from HTTP requests before processing them|

## Other tools
These are tools that are required to download in order to build Docker images and run Kubernetes network locally.

|Tool|Description|
|---|---|
|[Docker Desktop](https://www.docker.com/products/docker-desktop/)|Containerisation software for building and running container images|
|[minikube](https://github.com/kubernetes/minikube)|Lightweight Kubernetes implementation for fast local development|

## Core Concepts

### Domain Driven Design Architecture
The microservice architecture is structured as shown below where the data flow from an HTTP request follows the graph from top down in a DDD fashion.

<p align="center">
    <img src="./docs/DDD_structure.png" alt="ddd overview" />
</p>


### Echo
The framework used to serve HTTP requests. Here the request will be served to either a resolver from graphql or a controller for RESTful APIs.

### Route
Here is where the routes (API endpoints) are defined and an handler and middleware is associated with them.

#### Resolver
Resolvers are GraphQL construct for handling requests. They are described in the `schema.graphql` file.
An interface for each type of action (Query, Mutation) is created through gqlgen.
This interface holds all possible user-actions and needs to be implemented.
In this example, the implementation is done in the `src/infrasctructure/graph/resolver/resolver.go` file.
For larger projects it might be a good idea to use multiple files instead.

* **Logging**: GraphQL logs everything interesting happening in data-transmission. In this example gqlgen should perform this aspect of logging.

* **Errors**: Errors from the domain should be passed through from the application layer and reported to the user.

#### Controller
Here live the controllers of the RESTful APIs implemented.

* **Logging**: Ideally we should implement a logger but for this project but I have avoided it to keep things simple. Might implement in the future.

* **Errors**: Errors from the domain should be passed through from the application layer and reported to the user.

### Application
Creates and fetches data from the domain layer. Here CQRS is implemented to separate into queries and commands such that the former implements only read requests with no side effects while the latter implements read requests that have side effects and write requests.

### Domain
Responsible for representing concepts for the business, information about the business situation, and business logic. Within the domain there are repositories which define how data is accessed from the domain layer.

* **Logging:** Should not provide logic-based, but technical logging (e.g. connection-failures, timeouts).

* **Errors:** Should provide sensible domain based errors that user can interpret to avoid issues with the next requests. Not too much technical details, no sensitive informations.

#### Repository
The Repository layer is holding logic for interactions with any kind of data store. This could be a database or streaming-service.

A repository is always defined through an interface before implementing it. This is important to keep it testable. They live inside the `repository`-directory.

* **Logging:** Should not provide logic-based, but technical logging (e.g. connection-failures, timeouts). A not found database-row is mostly not a technical error, even if your db-driver returns an error in this case.

* **Errors:** Errors returned by this layer are not concerned for being seen by the end-user. Errors could contain sensible informations, we should filter this out.


#### Service
Your business-logic lives in the Service layer.
Every action a user can do should be implemented in a service layer, but is not to this.

A service could also provide functions used by other services, e.g. retrieving configuration.

Similar to repositories, services should always be described through an interface. So you could mock it during tests. They live inside the `service`-directory.

* **Logging**: Responsible for logic-based logging. Should log faulty actions (retrieving not existing entities, try accessing data without permission).

* **Errors**: Errors should be non-technical and suited for the end-user. Not too much technical details, no sensitive informations.


### Persistence
Implements repositories interfaces to provide access to permanent data stores hosted within the microservice itself (cache) or databases.

## Docker Compose
Docker Compose is used here for container orchestration as a quick option to test and develop the codebase. It can also be used for simple applications where there's no need to scale the application on multiple nodes. The architecture consists of two containers/services, one for PostgreSQL and another for the Golang server. A secret to share the database password across the containers has also been implemneted with the password stored in the `db/password.txt` file. 

You can quickly test it using `docker compose up`.

## Kubernetes
A Kubernetes container architecture is implemented under the `kubernetes`-directory which can be used alternatively of Docker Compose.
This architecture can be expanded and used for projects where scaling and flexibility are mission critical.

The architecture implemeted here is similar to the one used in Docker and it requires loading the server image built from Docker.

The kubernetes implementation can be started using the script at `kubernetes/start_kubernetes.sh` which setups the volumes, secrets, deployments and services. This script uses 'minikube' as it provides a simple Kubernetes cluster with a single node for development purposes. More information can be found in the script.

Other scripts have been provided for stopping all services and closing down the minikube VM and cluster using `kubernetes/stop_kubernetes.sh`. There is also an additional script called `clean_up.sh` inside the kubernetes folder that can be used to fully remove all deployments, volumes and services that might be left.

## CI/CD Pipeline
A continous integration (CI) approach is implemented using GitHub Actions to run a CI pipeline. In the `.github/workflows` folder a configuration yaml file has been written to use the docker compose configuration files under the folder `.ci`.
The dockerfile here is the same as `Dockerfile.dev` but it is using a docker image with golang binaries. This allows the CI to use the `go` binary to run the unit tests.

Other approaces could have been used. For example, I could have installed the `go` binary in the CI runner machine directly as part of the CI process and run the tests on the machine instead of inside the container. This approach is faster because it doesn't require to spin up the containers. This would be a better approach for unit tests and when functional tests will be implemented, I will separate the unit tests run from the functional tests run and only run the latter in the container since functional tests require the microservice to be running.

No approach for continous deployment is implemented as this is a demo project and it is not delivered anywhere.
I have provided a sample file at `.github/workflows/deploy.yml.sample` which demonstrates how a docker image is built and then pushed to DockerHub and then the image is pulled in the remote machine through SSH and the app is redeployed.

## TODOs

* Finish implementing graphql types, queries and modifications
* Implement authorisation check in graphql
* Implement Observability through Open Telemetry
