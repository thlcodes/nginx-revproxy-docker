<!-- Project Title -->

# ecomy service 
![](https://img.shields.io/badge/App_Version-v0.1.1_latest-red.svg)
![](https://img.shields.io/badge/Api_Version-v1-green.svg)

<!-- One Paragraph of project description goes here, it describes the functionality without details about the context in the system -->
Purpose of this ecomy is to make the creation of a new platform service as easy and fast es possible.

## Table of Contents

- [Background & Contexts](#background)
- [Local development](#local_development )
- [Usage](#usage)
- [CICD](#cicd)
- [Known Issues](#known_issues)
- [Relevant links](#links)

## Background & Contexts
This project is created to help solve the issue of started a new service from scratch or adding a new cross-functionality like using a logging library or interating service discovery.
This application is deployed in the dev environment and will be interacting with the service registry, api gw, and storage service.
It will also serve requests via a REST API. These requests are mainly for checking the general interaction between backend components.
If this application fails to function properly, there will be impact on other applications.

```
( internet ) ---> [ api-gw ] ---> [ gopher-user-ecomy ] 
                      ^                   |
                      |                   |
                      +----[ registry ]<--+
```     

## Local development / Getting Started 
### Prerequisites
specific to the actual service

### Dependencies
For a proper functioning, this application depends on the following services:


    +------------------+-------------------------------------------+
    |     Service      |    Purpose                                |
    |------------------|-------------------------------------------|
    | Storage service  | to store ecomy dara in mongo DB        |
    |------------------|-------------------------------------------|
    | Vehicles service | to store ecomy dara in mongo DB        |
    |------------------|-------------------------------------------|
    | Storage service  | to store ecomy dara in mongo DB        |
    |------------------|-------------------------------------------|
    | Users service    | to request basic user data                |
    +------------------+-------------------------------------------+

### Local build
#### Git clone
```bash
git clone https://collaboration.msi.audi.com/stash/scm/betacore/gopher-user-ecomy.git
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

###CICD
Build service : 
####Todo-list
[] enhance test coverage

## Runs on
* Docker
* Cloud foundry
* Kubernetes

## Authors
* [Thomas Liebeskind](https://collaboration.msi.audi.com/stash/users/thomas.liebeskind_audi.de)
* [Kamal Koubaa](https://collaboration.msi.audi.com/stash/users/kamal.koubaa_valtech-mobility.com)
* [Daniel Spelmezan](https://collaboration.msi.audi.com/stash/plugins/servlet/network/profile/daniel.spelmezan@valtech-mobility.com)

See also the list of [contributors](https://collaboration.msi.audi.com/stash/plugins/servlet/network/contributorsgraph/BETACORE/gopher-user-ecomy) who participated in this project.
