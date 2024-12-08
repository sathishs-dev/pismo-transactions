## PISMO TRANSACTIONS
***Service to manage and maintain transactions as well as accounts.***

This repository/service consists of the `pismo-transaction` service along with the `migrator` (DB migrations service), and a `PostgreSQL` database (version 17.2).

The entire service lifecycle, including building, running the services, and managing associated dependencies, is powered by Docker.

This project is built using **Go**, but basic functionalities such as building the service, executing it, and testing it do not require **Go**. However, extending its features requires **Go** to be installed on your machine.

---

## Table of Contents
1. [Dependencies & Installation](#dependencies--installation)
2. [Project Structure](#project-structure)
3. [Usage](#usage)
4. [API References](#api-references)
    1. [Create Accounts](#1-create-accounts)
    2. [Fetch Account](#2-fetch-account)
    3. [Create Transactions](#3-create-transaction)

---

## Dependencies & Installation
*This services is entirely built using **docker** as a infra. And it sticks multiple components together using **Makefile** and **docker compose**.*

So to get started here are the following prerequisites needed,

1. #### Clone this repository:

    over https
    ```bash
    git clone https://github.com/sathishs-dev/pismo-transactions.git
    ```

    over ssh
    ```bash
    git clone git@github.com:sathishs-dev/pismo-transactions.git
    ```

2. #### Install docker:
    > **Docker** is essential for the entire setup to work seamlessly.

    This Comprehensive guide from docker official site, helps us to install and get started with docker on your desired operating system.

    - https://docs.docker.com/engine/install/

    Once you’ve followed the guide, Docker should be installed on your machine. :)

    > **Bonus:** Docker Compose is bundled with the latest Docker installation, so you don't need to install it separately. You can use the new `docker compose` command instead of the older `docker-compose`.
 

3. #### Make:
    > This setup uses ***Makefile*** to reduce most of the complexities to build and run and organize the dependencies

    To install **GNU Make** processor, please follow along this

   **Install it on Linux ( Debian )**

    ```bash
    sudo apt update
    sudo apt install build-essential
    sudo apt install make
    ```

    **Install it on Mac**
    ```bash
    brew install make
    ```

    **Ensure the installation by running**
    ```bash
    make --version
    ```
4. #### Mockery(Optional):
    > This repo does requires [Mockery](https://vektra.github.io/mockery/latest/), if we are extending its functionalities, to generate mocks.

    To install mockery

    ```bash
    go install github.com/vektra/mockery/v2@v2.50.0
    ```

    And please make sure your go binaries path is added to **$PATH**

---
## Project Structure

```
    ├── cmd                     - Contains individual service entrypoints
    │   ├── migrator            - Service for manage db migrations
    │   └── pismo-transactions  - Service for handling the Transactions Routine
    ├── docker-compose.yml
    ├── Dockerfile
    ├── go.mod
    ├── go.sum
    ├── internal                - Contains basic utilities required only to the service
    │   └── meta
    ├── Makefile
    ├── migrator.Dockerfile
    ├── pkg                     - Contains batteries for db and handlers
    │   ├── enums
    │   ├── handler
    │   └── repository
    ├── Readme.md
    ├── schema                  - Contains migration files required for this service
    │   └── migrations
    ├── unittest.Dockerfile
    └── vendor
```

> **Note:** This project ( go project ) follows vendor approach to maintain its dependencies.

---

## Usage

1. **Build the services ( migrator & pismo-transactions ):**

    To build both services along with unit-test image

    ```bash
    make build
    ```

    To build migrator

    ```bash
    make build-migrator
    ```

    To build pismo-transactions service

    ```bash
    make build-service
    ```

2. **Run the unit tests:**
    > Note: Before running unit tests, ensure you built the unit test image using above step

    To run the unit tests

    ```bash
    make unit-test
    ```

3. **Run the services:**
    > Note: To run the services we do require the migrations applied to the db, our docker compose file desgined in a way to support this. 

    To run the service

    ```bash
    make run-service
    ```

    Now our service is listening in port 8080
    And db is listening in port 5432

4. **Optional:**

    To view the logs and follow through it after running the services,

    ```bash
    make logs
    ```

    To generate mocks while extending the features, ( it does requires mockery implementaion )

    ```bash
    make mock
    ```

    To update or install the dependencies, 

    ```bash
    make dep
    ```

    To destroy or brought down the service and its associated components

    ```bash
    make down
    ```
---
## API References

### 1. **Create Accounts**
- **Method**: `POST`
- **Endpoint**: `/accounts`
- **Description**: This endpoint creates a new account.

#### Request
- **Headers**:
    ```bash
        Content-Type: application-json
    ```
- **Body (JSON)**:
    ```json
    {
        "document_number": "1234567"
    }
    ```

#### Responses

- **Status Code**: `201`
    - **Description**: account created successfully
    - **Body** (Success): No Body

- **Status Code**: `400`
    - **Description**: invalid request / invalid body
- **Status Code**: `409`
    - **Description**: account already exists with document_number
- **Status Code**: `500`
    - **Description**: internal server error

- **Body** ( Failure ):
    ```json
    {
        "message": "<failure reason>"
    }
    ```

### 2. **Fetch Account**
- **Method**: `GET`
- **Endpoint**: `/accounts/:accountId`
- **Description**: This endpoint fetches account for :accountId passed.

#### Request
- **URL Param**:
   `accountId: (int)`

#### Responses

- **Status Code**: `200`
    - **Description**: account fetched successfully
    - **Body** (Success):
        ```json
        {
            "account_id": 1,
            "document_number": "document_number"
        }
        ```

- **Status Code**: `400`
    - **Description**: invalid request / invalid body / account doesn't exists

- **Status Code**: `500`
    - **Description**: internal server error

- **Body** ( Failure ):
    ```json
    {
        "message": "<failure reason>"
    }
    ```

### 3. **Create Transaction**
- **Method**: `POST`
- **Endpoint**: `/transactions`
- **Description**: This endpoint creates a new transaction record.

#### Request
- **Headers**:
    ```bash
        Content-Type: application-json
    ```
- **Body (JSON)**:
    ```json
    {
        "account_id": 1,
        "operation_type_id": 4,
        "amount": 123.45
    }
    ```

#### Responses

- **Status Code**: `201`
    - **Description**: transactions created successfully
    - **Body** (Success): No Body

- **Status Code**: `400`
    - **Description**: invalid request / invalid body / account not found / operation not not found

- **Status Code**: `500`
    - **Description**: internal server error

- **Body** ( Failure ):
    ```json
    {
        "message": "<failure reason>"
    }
    ```
---
