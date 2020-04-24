<Master README Template>

# Master README Project title
<!-- One Paragraph of project description goes here, it describes the functionality without details about the context in the system -->

- Background/Contexts
`This should be a relatively long section putting every app category in its context and discribe thee overall interaction between each componant in the system`

- Local development
`general information about running the whole system localy and check basicc functionality`
  - prerequisites `here goes the instactions for installing required tools for setting up the dev env, 
  like : docker, task, DB clients, protoc, ...`
  - dependencies `here goes an overview of all dependences between system componants.`
  - code checkout `List git repos to checkout and build`
  - running automated tests `describe generally the standard and basic way to run tests, like "go test ./...`
  - Services configuration `description of the different configuration information needed for the application to start locally`
  - start the system
  - running e2e tests
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
- CICD
  - Versioning
  - System specifics for every environment
  - How to release a new version
  - How to trigger a build
  - How to trigger a deployment
- Known Issues    
- Relevant links
  - Related Confluence documentation
  - Build pipeline
  - Deploy pipeline
  - What is Clean Architecture
  - What is Micro-services Architecture
  - What is gRPC