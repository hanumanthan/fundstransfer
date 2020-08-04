## FundsTransfer
API to transfer funds between two wallets

## Running the service
This service can be started by any of the below 3 options

### docker-compose
```bash
docker-compose up
```
This will start the service and port forward to 8080  
Note - To rebuild the image run `docker-compose build` followed by `docker-compose up`

### docker run
```bash
docker build -t fundstransfer:1.0 .
docker run -p 8080:8080 --name fundstransfer fundstransfer:1.0
```
The image is built with the tag **fundstransfer:1.0** and container is started as **fundstransfer**

### Run locally
```bash
go run main.go
```

## Tests
```bash
go test -v ./...
```
To execute the tests


## Choices/Assumptions
1. DB Transactions are not handled.
2. Database is recreated by docker when the image is re-built. Volumes are not retained
3. Due to simplicity of the app, requests execute in single go routine. DB queries are not handled in seperate goroutines
4. Each wallet is awarded $100 when creating the app
5. Authentication is done via api_key instead of username/password or basic auth
6. The code structure follows layered architecture due to the simplicity of the app
7. If started via docker, logging happens as inside the container. Log files are not mounted to host system
8. As the app heavily relies on SQL operations, integration tests have been written to have better coverage. Unit tests wont be effective as most of the code will get mocked.
9. Gorm and Gin errors are handled but they are not translated to user errors. 
10. api_key for each user can be found in [file](pkg/models/setup.go)

## Initial Data
Data setup is done via this [file](pkg/models/setup.go)

## API documentation
1. `http://localhost:8080/liveness/healthcheck` - returns a 200 if the app is working and hasnt suffered by deadlocks or panics
2. `http://localhost:8080/metrics` - Provides metrics like the error counts and response times for the user apis
3. Admin apis - 
Admin apis needs `api_key` to be passed via header.
Header

    1. `http://localhost:8080/admin/users`  - Gets all users in the system
    2. `http://localhost:8080/admin/wallets`  - Gets all wallets in the system
    3. `http://localhost:8080/admin/transactions`  - Gets all transactions in the system
4. User apis - 
User apis needs `api_key` for the corresponding user to be passed via header.
Header

    1. `http://localhost:8080/api/v1/user/:user_id/transact`  - To transfer money from one user to another user
    Post body
    2. `http://localhost:8080/api/v1/user/:user_id`  - To get the user details like name, balance along with credit and debit transactions
    
    

## Killer feature
1. For every $100 transferred by a user, user gets a scratch card valued between $1-$20 added as cashback to users wallet
2. For every $50 transferred by user, user gets 10% discount for the next cab ride if paid using the wallet

### TODO
1. Add random delay
2. Killer feature