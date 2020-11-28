# article-api

#### Run the Testing

```bash
$ make test
```

#### Run the Application

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
$ docker ps

# Execute the call
$ curl localhost:9090/articles

# Stop
$ make stop
```
