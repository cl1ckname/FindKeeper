.PHONY: build docker-build docker-push clean run

BOT_NAME := find-keeper
REGISTRY := <your registry>
IMAGE_NAME := $(REGISTRY)/$(BOT_NAME)
TAG := latest

build:
	go build -o bin/$(BOT_NAME) ./cmd/find-keeper

docker-build:
	docker build -t $(IMAGE_NAME):$(TAG) .

docker-push: docker-build
	docker push $(IMAGE_NAME):$(TAG)

clean:
	rm -f $(BOT_NAME)
	docker rmi $(IMAGE_NAME):$(TAG) 2>/dev/null || true

run: build
	./$(BOT_NAME)

deploy: docker-push
	@echo "Image pushed to $(IMAGE_NAME):$(TAG)"

install-deps:
	go mod download
	go mod tidy
