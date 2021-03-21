# Puppet Master

Simple Authentication and Authorization app in Go.

## Status

In development.

## Architecture

Clean Architecture (Robin C. Martin).

## Prerequisites

- Go 1.16.x

- Docker 20.10.x

## Run

Run database (Postgre):

```sh
 docker-compose up -d
```

Then, run the server:

```sh
 cd cmd/puppet_master && go run main.go
```

Access:

- GraphiQL: http://localhost:8000/graphql

- Adminer: http://localhost:8080
