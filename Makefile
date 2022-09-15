BINARY_NAME=app

build:
	echo "Building..."
	go mod download
	go mod verify
	go build -o ${BINARY_NAME} cmd/main.go
	echo "Finished..."

clean:
	go clean
	rm ${BINARY_NAME}
	docker-compose down --remove-orphans

lint:
	go vet ./...
	staticcheck ./...
	golint ./...

run-with-docker:
	docker-compose up --build

run-locally:
	docker-compose up --build -d
	make build
	./app

test:
	go test -v ./...
	go test ./... -coverprofile=coverage.out
