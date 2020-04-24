<!-- Project Title -->
# Gutsy Gopher backend Services - General README

Gutsy Gopher is the Betaspace platform backend. it is microservice architecture based system build with Go and 
implements the required business logic for a faster time to market application development for Audi connected cars.

## Table of Contents

- [Background & Contexts](#background)
- [Local development](#local_development )
- [Usage](#usage)
- [CICD](#cicd)
- [Known Issues](#known_issues)
- [Relevant links](#links)

## Background & Contexts
Gutsy gopher bachend is a set of core services that serve requests through an API Gateway. This backend is connected to :
* the User portal
* iOS App
* Alexa
* ...

The api gw forwards requests to backend core services after checking a couple of requirements.
The registry helps the api gw and core services to dermine if the target service available and provide the corrsponding endpoint url.


```
                                                      core svcs
                                                  +----------------+       
                                                  | [ users ]      |           abstraction layers
    browser ---> portal bff --+                   | [ services ]   |           +----------------+
                              |                   | [ disclaimer ] |           | [ seriesAl ]   |  ---> [ external services ]
                              +-->[ api-gw ] ---> | [ token ]      | --------> | [ storageAl ]  |
                              |        |          | [ vehicles ]   |           | [ cacheAl ]    |  ---> [ DB ]
           app ---> app bff --+        |          | [ feedback ]   |           +----------------+
                                       |          +----------------+                   ^
                                       |                  ^                            |
                                       |                  |                            |
                                       |                  |                            |
                                       |                  |                            |
                                       +----------->[ registry ]<----------------------+
```     

## Local development / Getting Started 
### Prerequisites
* #### Git: 
  This project uses GIT for SCV
* #### Go
* #### Task
* #### Docker
Docker should be installed on the dev machine.
To run the application locally, it requires a running database and other application to be already started in containers.


    +---------+-------------------------------------------------+
    |  Tool   |    Purpose                                      |
    |---------|-------------------------------------------------|
    | Git     | This project uses GIT for SCV                   |
    |---------|-------------------------------------------------|
    | Go      |                                                 |
    |---------|-------------------------------------------------|
    | Task    | Task is used to run command locally for         |
    |         | common task like testing, code generation, ...  |
    |---------|-------------------------------------------------|
    | Doccker | Docker should be installed on the dev machine.  |
    |         | To run the application locally, it requires a   |
    |         | running database and other application to be    |
    |         | already started in containers.                  |
    +---------+-------------------------------------------------+

### Dependencies
For a proper functioning, this application depends on the following services:


    +------------------+-------------------------------------------+
    |     Service      |    Purpose                                |
    |------------------|-------------------------------------------|
    | Storage service  | to store skeleton dara in mongo DB        |
    |------------------|-------------------------------------------|
    | Vehicles service | to store skeleton dara in mongo DB        |
    |------------------|-------------------------------------------|
    | Storage service  | to store skeleton dara in mongo DB        |
    |------------------|-------------------------------------------|
    | Users service    | to request basic user data                |
    +------------------+-------------------------------------------+

### Code checkout
#### branching model
* master: this branch is the latest stable version, to checkout a specific release, search for the corresponding git-tag
* develop: this branch is always under active development, to get the latest chages this is the right place.
* other branches: bugfix or feature branch
#### Git clone
```bash
git clone https://collaboration.msi.audi.com/stash/scm/betacore/gopher_skeleton.git
```
#### running automated tests
To run tests
```bash
go test ./...
```
To run tests with coverage
```bash
gocov test ./...
```

#### Service configuration
#### running the application
#### running simple manual tests like smoke tests or health check requests
### Usage
#### Use cases
* ##### Manage user information
* ##### manage user services
* ...
#### Provided Interfaces
* ##### manage user information
     - GetUserinfo: get userinfo from myAudi idp and store data in users DB
     - GetBetaflag: get the value of the betaflag and check if it is accepted and latest 
     - ...
#### REST API - Request examples
more details about the Rest API is found on the swagger page [here]()

### Developement
#### Code style
To check code style:
```bash
go vet ./...
```
#### Available and relevant libraries
main used libraries are  :
* bootstrap
* cache
* common.go
* discovery
* errors
* goflow
* grpc
* http
* openapi

For more details check [lib-go-common]() documentation:
* [README]()
* [Confluence]()
 
#### Writing tests
#### Error handling
#### Logging/Monitoring

## CICD
### Versioning
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags).
### Configuration for a specific environment
### How to release a new version
### How to trigger a build
### How to trigger a deployment
## Known Issues
List of issues, limitations and things that are still to be done.

#### Todo-list
[] enhance test coverage

## Relevant links
- [Related Confluence documentation]()
- [Build pipeline]()
- [Deploy pipeline]()
- [What is Clean Architecture]()
- [What is Micro-services Architecture]()
- [What is gRPC]()
  
## Built With
* [Clean architecture]()
* [Go]()
* [gRPC]()

## Runs on
* Docker
* Cloud foundry
* Kubernetes
* All OSs that Go supports.

## Documentation

## Authors
* [Thomas Liebeskind](https://collaboration.msi.audi.com/stash/users/thomas.liebeskind_audi.de)
* [Kamal Koubaa](https://collaboration.msi.audi.com/stash/users/kamal.koubaa_valtech-mobility.com)
* [Daniel Spelmezan](https://collaboration.msi.audi.com/stash/plugins/servlet/network/profile/daniel.spelmezan@valtech-mobility.com)

See also the list of [contributors](https://collaboration.msi.audi.com/stash/plugins/servlet/network/contributorsgraph/BETACORE/gopher_skeleton) who participated in this project.
