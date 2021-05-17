# Development environment build instructions

## Environment variables

There's a .env file included in the repo with some nice default values for enviroment variables and their descriptions. This will be loaded by the app and the docker integration test system unless overriden in the cli.
## Pre-requisites

    * [go](https://golang.org/)
    * [docker-compose](https://docs.docker.com/compose/)
## Directory setup

After download from github run:

    go mod tidy

To download all the go modules.

## Integration testing

Presently to run one integration test:

  docker-compose down
  go run .
  docker-compose up -d

TODO:

 * Make the go program clear the database at the start
 * Make a proper go integration test instead of it being the main program

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
