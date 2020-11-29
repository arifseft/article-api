# article-api

## Run the Testing

```bash
$ make test
```

## Run the Application

Here is the steps to run it with `docker-compose`

```bash
#move to directory
$ cd workspace

# Clone into YOUR $GOPATH/src
$ git clone https://github.com/arifseft/article-api.git

#move to project
$ cd article-api

# Build the docker image first
$ make docker

# Run the application
$ make run

# check if the containers are running
$ docker-compose ps

# Execute the call
$ curl localhost:9090/articles

# Stop
$ make stop
```

## REST API

### Post Article

#### Request

```bash
$ curl --location --request POST 'http://localhost:9090/articles' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "title",
    "body": "lorem ipsum body",
    "author": "john wick"
}'
```

#### Response

```bash
{
    "id": 1,
    "title": "title",
    "body": "lorem ipsum body",
    "author": "john wick",
    "created_at": "2020-11-29T07:55:24.2709948Z"
}
```

### Search Articles

#### Request

```bash
$ curl --location --request GET 'http://localhost:9090/articles?query=em&author=john%20wick%201'
```

#### Response

```bash
{
    "data": [
        {
            "id": 1,
            "title": "title",
            "body": "lorem ipsum body",
            "author": "john wick",
            "created_at": "2020-11-29T07:55:24.2709948Z"
        }
    ],
    "message": null
}
```
