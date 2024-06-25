# service-catalog

service-catalog is a backend service written in go to support UI fronted. here is the [figma](https://www.figma.com/file/zeaWiePnc3OCe34I4oZbzN/Service-Card-List?node-id=0%3A1) for UI. 

It supports addition of a service into the service catalog and retriving a service informatoin.

## Features 

1. service catalog exposes create API to create a new service instance. 
2. service catalog exposes GET API to get a particular service and also users can list the services available.


## Architecture 


### System overview

<img width="666" alt="Screenshot 2024-06-25 at 6 03 35 PM" src="https://github.com/rajendragosavi/service-catalog/assets/36725053/6629b6fd-5112-4cf7-a354-c9d4f81b4bdf">


We have our backend service exposes APIs to clients which can be any frontend or any REST client. Upon receiving a request, corresponding web handler interacts with database to return the results



**SQL v/s NoSQL**

We are using postgres database which is a relational database to store the service details. 

**Why  ?**

There few reasons why we have chosen the SQL for this case :  

1. The service object which we are dealing with is a well defined strctured data also if in future even if we add more features/tables we can easily build relationships with mulitple tables. This is classic use case for relational database.

2. This CRUD APIs are ideal use cases for - relational databases where data consistency and ACID transactions are critical.



### Code Overview

<img width="1005" alt="Screenshot 2024-06-25 at 6 03 24 PM" src="https://github.com/rajendragosavi/service-catalog/assets/36725053/4cf4d00c-c288-43a6-9931-eb1b74098cdd">




If you see above diagram we have isolation, between Repository, Model, Service and Delivery layers.

```Delivery (REST, gRPC) -> Business Service -> Repository -> Model```

We are not exposing database layer to the buiseness logic. Instead we create a repository interface which uses models (_"model" refers to a struct type that represents a table in your database. This struct will have fields that correspond to the columns of the table. Models are used to map the database structure to Go code, making it easier to interact with the database using Go’s type system_)  

We have http handlers seperated and they are using service interface which is again another layer to connect with bbuiseness logic. The idea was to have more modular handlers.


Whenever a request comes to any API - the http handler corresponding to the API calls corresponding service method provided by service interface, and then this service method uses methods exposes by repository interface to interact with database.


This modular approach will enable us to segregate and implement model, API handlers and buisness logic easily.

Go Interfaces in brief : 

Repository interface - which provides methods to interact with database from buisness logic

```go
type Repository interface {
	Create(ctx context.Context, obj *model.ServiceCatalog) (string, error)
	Get(ctx context.Context, serviceName string) (*model.ServiceCatalog, error)
	List(ctx context.Context) ([]model.ServiceCatalog, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	CreatewithCommit(ctx context.Context, obj *model.ServiceCatalog) error
}
```

Service interface - which provides methods to cater buiseness logic

```go
type ServiceCatalogService interface {
	Create(ctx context.Context, params CreateParams) (string, error)
	List(ctx context.Context) ([]*model.ServiceCatalog, error)
	Get(ctx context.Context, name string) (*model.ServiceCatalog, error)
}
```

## How to Build 


### Pre-requisites

* Go 1.22 
* Git
* make (To run commands with ease)
* Docker (If you want to build and run docker container)
* Postgres (Dont worry - If you dont have a postgres running on your machine.
You can use docker to host postgres for you.) 

Run - 

1. ``` docker run -d --name <name of the container> -p 5432:5432 -e POSTGRES_PASSWORD=<password> postgres``` 

2. ```docker exec -it <name of the container>  psql -U postgres```

3. Once you exec into the container. create sql database - service -  and connect to it - ```\c service```

4. Run following SQL commands to setup table and schemas.

```sql

-- create new serice database
CREATE DATABSE service;

-- use service database
USE serice;

-- Enable the pgcrypto extension for UUID generation
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create the services table
CREATE TABLE service (
    service_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(200),
    status SMALLINT NOT NULL,
    creation_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_updated_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deletion_time TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    versions TEXT[] DEFAULT '{}',
    is_deleted BOOLEAN DEFAULT FALSE
);
```


This will run a postgres in docker container and make it accessible over localhost on port 5432
 

### Installation

1. Clone the github repo 

	```
	git clone https://github.com/rajendragosavi/service-catalog
	```

2. Run setup (for database configs. Edit the database password)

	```
	make setup
	```


2. Build the Project:
    ```
    make build
    ```

3. Rub the binary in local
    ```
    make run
    ```

4. Run Unit test
    ```
    make test
    ```

5. Build Docker Image
    ```
    make docker-build
    ```

6. Run the docker image
    ```
    make deploy
    ```

7. Delete the binary generated
    ```
    make clean
    ```

## How to test

#### create a service 

Request

```sh
curl --location 'http://localhost:80/api/v1/services' \
--header 'Content-Type: application/json' \
--data '{
	"name":"service-2",
	"description" : "service-2 description",
	"versions":[
		"1.0"
	]
}'
```

Response 

```json
{
    "id": "5fdeff89-2261-4fa9-b8fd-d6117473c90d"
}
```

#### Get a specific service 

Request 
```sh
curl --location 'http://localhost:80/api/v1/services/service-2'
```

Response

```json
{
    "service_id": "5fdeff89-2261-4fa9-b8fd-d6117473c90d",
    "service_name": "service-2",
    "description": "service-2 description",
    "created_time": "2024-06-25T07:35:03.434233Z"
}
```

#### Get all services

Request

```sh
curl --location 'http://localhost:80/api/v1/services'
```

Response

```json
[
    {
        "ID": "4d52f7af-4a54-4c84-963a-41633181b64f",
        "Name": "service-1",
        "Description": "service-1 description",
        "CreatedOn": "2024-06-24T04:37:51.3456Z",
        "UpdatedOn": null,
        "DeletedOn": null,
        "Versions": [
            "1.0"
        ],
        "IsDeleted": false
    },
    {
        "ID": "cb30cb70-50b9-4f0f-8fc8-c05ece379f89",
        "Name": "service-2",
        "Description": "service-2 description",
        "CreatedOn": "2024-06-24T04:54:49.000238Z",
        "UpdatedOn": null,
        "DeletedOn": null,
        "Versions": [
            "1.0"
        ],
        "IsDeleted": false
    }
]
```

## Feature Status

✅  We have three apis available as of now for create service , get a particular service , list all services.

✅ We have clear repository interface for database access.

✅ We have an interface for service catalog which exposes a way from http handlers to interact with database

✅ We have unit tests for repositor access which enables testing our features using mocks. I have tried to handle dependency injection at service , handler and repository level which avoids [monkey patching](https://en.wikipedia.org/wiki/Monkey_patch#:~:text=In%20computer%20programming%2C%20monkey%20patching,altering%20the%20original%20source%20code.). 

✅ Any idea production grade service should have versioned apis.
 We have versioned APIs. 

✅ We have structured logging in place using logrus.

✅ Our webserver handles Graceful termination. It checks if all open transactions are completed or not before terminating the server.

### Following tasks are WIP

🔶  Adding Integration tests.

🔶  Currently filtering and pagination is not implemented for list API.

🔶  swagger doc is WIP in swagger-1 branch.

🔶  POST to create a new service. service_name is a unique key in the database  [PENDING]

🔶  PUT to update an existing service [Pending]