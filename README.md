# My tech-challenge
This rest api application can authenticate and authorize users to manage tasks.

## Demo of the manager notification
![](assets/demo.gif)


## Building and Running
### locally 

```bash
# Up the api and database services in docker containers
docker-compose up -d # detached
# or
make run-with-docker 

# Only build the binary
make build

# Build and Run the binary
make run-locally

# Test
make test

# Health check endpoint: http://127.0.0.1:8060/healthz
```

## Environment variables
```bash
# APPLICATION
LOG_LEVEL="Debug"

# HTTP SERVER
HTTP_SERVER_PORT="8060"

# MYSQL
MYSQL_USER=dev
MYSQL_PASSWORD=dev
MYSQL_DB_NAME=dev
MYSQL_HOST="0.0.0.0"
MYSQL_DB_PORT=3306

MYSQL_DATABASE=dev
MYSQL_ROOT_PASSWORD=dev
```

## Endpoints
#### /healthz
* `GET` : Check the api health status  
#### /login/authorize
* `POST` : Validate the user credentials and return an jwt access token  
    * input
        * user (string)
        * password (string)
    * output
        * message (string)
        * token (string) // should be send in the next requests headers

#### /tasks?user_id={user_id}
* `GET` : Get a list of tasks that can be filtered by user_id
     _(needs the bearer token inside the request header)_
    * input
        * role (string)
        * user_id (string)
    * output
        * array
            * id (string)
            * user_id (string)
            * summary (string)
            * perform (string)
* `POST` : Add a new performed task
    * input
        * user_id (string)
        * summary (string)
        * performed_at (datetime)
    _(needs the bearer token inside the request header)_


TODO:
- more unit tests (added only the login_handler_test.go, missing others)
- enconde / decode sensitive data in task summary and user password fields
- replace all fmt.* by log
- add the DELETE task route
