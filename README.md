# Go Application

Prerequisites: _Go, Node, Docker_

Installation: `npm install -g generator-jhipster`

Link:

`git clone https://github.com/tic-oss/jhipster-blueprints.git `

`cd generator-jhipster-gomicro`

`npm link generator-jhipster-gomicro`

A Sample microservices application can be generated using the provided jdl: `jhipster jdl reminder.jdl`

reminder.jdl

---

<pre>
   application {
    config {
      baseName gomicro,
      applicationType microservice,
      packageName com.cmi.tic,
      authenticationType oauth2,
      databaseType sql,
      prodDatabaseType postgresql,
      devDatabaseType postgresql,
      serviceDiscoveryType eureka,
      blueprints [gomicro],
      serverPort 9002,
    } 
}
 </pre>

---

Below is the directory structure once application code is generated.

```
.
├── be1 (backend microservice 1)
│   ├── README.md
│   ├── auth
│   ├── config
│   ├── controllers
│   ├── db
│   ├── eurekaregistry
│   ├── handler
│   ├── migrate
│   ├── proto
│   ├── resources
│   ├── rabbitmq
│   ├── main.go
|   ├── go.mod
|   ├── Dockerfile
|   ├── docker-compose (docker compose for all microservices and their dependencies)
│   ├── central-server-config
│   ├── docker-compose.yml
└── └── realm-config
```

The generated application also includes keycloak, jhipster registry (based on eureka), postgres mongodb rabbitmq server specific to each service.

## Intro on the services provided

    - auth - The authorization is supported by  [Keycloak](https://www.keycloak.org).Here The token is fetched from the header and retrospection is done for the validity of the token if the token is valid the api is hit on the contrary an error is returned.

    - config - Next for the config there are multiple environments available initially the app.yml configuration is loaded then based on the env the particular env specific file is loaded and finally the env variables are overided thus on the same structure multiple values will be overriden.

    - resources - Here the app.yml and the respective config files are stored.

    - controllers -  Implemented by [mux](https://github.com/gorilla/mux).Here the list of api calls and their respective functionalities are written where whena specific api hit what function must be called.

    - db - Here the link is made to the database and the application and specific client is returned as an output.

    - eureka registry - This file contains the code for Registering the service with eureka Service registry with the configurations.

    - handler - Here the logic for the api calls of the controller is present for the db queries,CRUD operations.

    - migrate - Implemented by [gomigrate](https://github.com/golang-migrate/migrate).This is for the database versioning in golang.Here the database is created with the application basename and the scripts are run to create table along with schema and insert few data.

    - rabbitmq - This file consists of files of consumer and producer for asynchronous communication i.e, for sending data and recieving data from the queue.

    - Dockerfile - To build an image of the microservice.

    - main.go - The source file from where the application starts.

## To run the golang application generated

- Build the go application

---

      cd go
      go mod tidy

---

- First run the postgress or mongodb,keycloak and jhipster if rabbitmq for connection then rabbitmq registry files.

---

      cd docker
      docker-compose -f mongodb.yml up
      docker-compose -f  keycloak.yml up
      docker-compose -f  jhipster-registry.yml up
      docker-compose -f rabbitmq.yml up

---

- Now get back to the root directory of go and start the golang service

---

      go run .

---

## Build an image

---

    docker build -t image_name .

---

### Run by passing values

- docker run -it --network="host" -p <port_number>:<port_number> -e db_url=postgresql://go@localhost:5433/postgres image_name

### Run by passing 3 ENV values

```
docker run -it --network="host" -p <port_number>:<port_number> \

-e DB_URL=postgresql://go@localhost:5433/postgres \

-e KEYCLOAK_URL=http://localhost:9080/realms/jhipster \

-e SERVICE_REGISTRY_URL=http://admin:admin@localhost:8761/eureka/apps/ \

be
```

### Code Formatting

- Gofmt formats Go programs. It uses tabs for indentation and blanks for alignment. Alignment assumes that an editor is using a fixed-width font.

Usage:

```
 gofmt [flags] [path ...]
```

- Flags

  - d - Do not print reformatted sources to standard output.
    If a file's formatting is different than gofmt's, print diffs
    to standard output.

  - e - Print all (including spurious) errors.

  - l - Do not print reformatted sources to standard output.
    If a file's formatting is different from gofmt's, print its name
    to standard output.

  - r rule - Apply the rewrite rule to the source before reformatting.

  - s - Try to simplify code (after applying the rewrite rule, if any).

  - w - Do not print reformatted sources to standard output.

- For formatting whole directory

  ```
    gofmt  .
  ```
