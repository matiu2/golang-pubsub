# Development environment build instructions

## Pre-requisites

    * [go](https://golang.org/)
    * [docker-compose](https://docs.docker.com/compose/)

## Directory setup

After download from github run:

    go mod tidy

To download all the go modules.

## Integration testing

To create the testing *mysql* and *pubsub-emulator* docker images:

    docker-compose up -d

(see [docker-compose.yaml](docker-compose.yaml) for that configuration)

There's a file included [.env](.env), that lets you connect to these services. You'll need to modify this for test/prod deployments.

### Environment access

#### MySQL

 * To run mysql-commands:
   docker-compose exec mysql mysql -uroot -proot pubsub
 * mysql listens no localhost 3307

#### Pubsub emulator

 * Details are loaded from the .env file:

   PUBSUB_EMULATOR_HOST=localhost:8085
   PUBSUB_PROJECT_ID=test
