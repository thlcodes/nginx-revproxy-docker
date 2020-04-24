# Project title
![](https://img.shields.io/badge/App_Version-v0.1.1_latest-red.svg)
![](https://img.shields.io/badge/App_Version-v0.1.1-blue.svg)
![](https://img.shields.io/badge/Api_Version-v1-green.svg)
![](https://img.shields.io/badge/Code_quality-A+-green.svg)
![](https://img.shields.io/badge/Build-passing-green.svg)

## README Template

- Background/Contexts
`This part is about the actual app in the context of the whole system. It would help 
the developer understand what functionality to expect from the application`
- Local development
`Information about how to prepare the local environment for building, testing and developping the application code`
  - prerequisites `here goes the instactions for installing required tools for setting up the dev env, 
  like : docker, task, DB clients, protoc, ...`
  - dependencies `here goes the other components on which the application depends, ex. 
  "users" services depends on:
   disclaimer services in-order to get the latest disclaimer version and compare it with the betaflag accepted version
   storage al service in-order to stoge user-related data in the longterm storage`
  - code checkout `basic git clone to run directly in the terminal, or copy the url an use it in the IDE`
  - running automated tests `describe the standard and basic way to run tests, like "go test ./...`
  - Service configuration `description of the different configuration information needed for the application to start locally`
  - running the application
  - running simple manual tests like smoke tests or health check requests
- Usage `now that we have the code built and tested, we can go ahead with advanced usage, as needed`
  - use cases `list the main use cases the application provides, ex:
    -manage user information
    -manage user services
    -...`
  - Provided Interfaces `deeper details about how to fulfill specific use case:
    -manage user information
     -GetUserinfo: get userinfo from myAudi idp and store data in users DB
     -GetBetaflag: get the value of the betaflag and check if it is accepted and latest 
     -...`
  - request examples
- Known Issues