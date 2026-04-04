.PHONY: build docker-build docker-push clean run

BOT_NAME := find-keeper
REGISTRY := registery.toadonload.ru
IMAGE_NAME := $(REGISTRY)/$(BOT_NAME)
TAG := latest

build:
	go build -o bin/app ./cmd/find-keeper

docker-build:
	docker build --network=host -t $(IMAGE_NAME):$(TAG) .
	docker build -t $(IMAGE_NAME)-wg:$(TAG) -f Dockerfile.wireguard .

docker-push: docker-build
	docker push $(IMAGE_NAME):$(TAG)
	docker push $(IMAGE_NAME)-wg:$(TAG)

clean:
	rm -f $(BOT_NAME)
	docker rmi $(IMAGE_NAME):$(TAG) 2>/dev/null || true

run: build
	./bin/$(BOT_NAME)

deploy: docker-push
	@echo "Image pushed to $(IMAGE_NAME):$(TAG)"

install-deps:
	go mod download
	go mod tidy
