# service-catalog

service-catalog is a backend service written in go to support UI fronted. here is the [figma](https://www.figma.com/file/zeaWiePnc3OCe34I4oZbzN/Service-Card-List?node-id=0%3A1) for UI. 

It supports addition of a service into the service catalog and retriving a service informatoin.

## Features 

1. service catalog exposes create API to create a new service instance. 
2. service catalog exposes GET API to get a particular service and also users can list the services available.


## Architecture 


### System overview

<img width="666" alt="Screenshot 2024-06-25 at 6 03 35 PM" src="https://github.com/rajendragosavi/service-catalog/assets/36725053/6629b6fd-5112-4cf7-a354-c9d4f81b4bdf">


from above diagram - we have our backend service exposes APIs (add link to APIs doc) to outside world. 

We are using postgres database to store the service details. 

Why SQL and not No-SQL ?

There few reasons why we have chosen the SQL for this case - 

1. The service data is a well defined strctured data and also if in future even if we add more features/tables etc we can easily build relationships with mulitple tables. this will still be a good choice to build robust and consistent queries.

2. This CRUD APIs are classic use cases for - relational databases where data consistency and ACID transactions are critical.


### Code Overview

<img width="1005" alt="Screenshot 2024-06-25 at 6 03 24 PM" src="https://github.com/rajendragosavi/service-catalog/assets/36725053/4cf4d00c-c288-43a6-9931-eb1b74098cdd">




If you see above diagram we have isolation, between Repository, Model, Service and Delivery layers.

Here we are not exposing the repository directly to the Delivery layer, instead only our business service has access to it. And once we have received an API call we just expose the business service to the delivery layer. 

This modular approach will enable us to segregate and implement model, API handlers and buisness logic easily.

We have repository interface{} - which provides methods to interact with database from buisness logic

```go
type Repository interface {
	Create(ctx context.Context, obj *model.ServiceCatalog) (string, error)
	Get(ctx context.Context, serviceName string) (*model.ServiceCatalog, error)
	List(ctx context.Context) ([]model.ServiceCatalog, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	CreatewithCommit(ctx context.Context, obj *model.ServiceCatalog) error
}
```

and Service interface{} - which provides methods to cater buiseness logic

```go
type ServiceCatalogService interface {
	Create(ctx context.Context, params CreateParams) (string, error)
	List(ctx context.Context) ([]*model.ServiceCatalog, error)
	Get(ctx context.Context, name string) (*model.ServiceCatalog, error)
}
```

### List of third party libraries used 

1. ``` https://github.com/gorilla/mux ``` for http request routing and handling
2. ```https://github.com/sirupsen/logrus```  Logrus is a structured logger for Go (golang), completely API compatible with the standard library logger.
3. ```https://github.com/jmoiron/sqlx``` sqlx is a library which provides a set of extensions on go's standard database/sql library
4. ```https://github.com/lib/pq``` A pure Go postgres driver for Go's database/sql package
5. ```https://github.com/asaskevich/govalidator``` A package of validators and sanitisers for strings, structs and collections
6. ```https://github.com/joho/godotenv``` for managing configuration data from environment variables (we can use viper as well)


## How to Build 


### Pre-requisites

* Go 1.22 
* Git
* make (To run commands with ease)
* Docker (If you want to build and run docker container)
* Postgres (Dont worry - If you dont have a postgres running on your machine.
You can use docker to host postgres for you.) 

Run - ``` docker run -d --name <name of the container> -p 5432:5432 -e POSTGRES_PASSWORD=<password> postgres``` 

This will run a postgres in docker container and make it accessible over localhost on port 5432
 

### Installation

1. Clone the github repo 

	```
	git clone https://github.com/rajendragosavi/service-catalog
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
