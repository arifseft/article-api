BINARY=engine
test: 
	go test -v -cover -covermode=atomic ./...

coverage:
	go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

engine:
	go build -o ${BINARY} app/*.go

unittest:
	go test -short  ./...

lint-prepare:
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint:
	./bin/golangci-lint run ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t article-api .

build: stop
	docker-compose up --build

run:
	docker-compose up

stop:
	docker-compose down

watch:
	@air -c air.conf


.PHONY: clean install unittest lint-prepare lint build docker run stop watch
