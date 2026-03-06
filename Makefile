.PHONY: build run tidy clean docker-build docker-push

IMAGE ?= unsw-comp3900-app

build:
	go build -o bin/server .

run:
	go run .

tidy:
	go mod tidy

clean:
	rm -rf bin/

docker-build:
	docker build -t $(IMAGE) .

docker-push: docker-build
	docker push $(IMAGE)
