ifeq ($(NAME),)
	NAME := "safechildhood"
endif

ifeq ($(TAG),)
	TAG := "latest"
endif

IMAGE_TAG := ${NAME}:${TAG}

build: deps
	go build -o ./bin/safechildhood ./cmd/safechildhood/main.go
	
run: build
	mkdir log || true
	./bin/safechildhood

build-docker:
	./scripts/docker_object_exist.sh ${NAME} && docker rm ${NAME} || true
	./scripts/docker_object_exist.sh ${IMAGE_TAG} && docker rmi ${IMAGE_TAG} || true

	docker build -t ${IMAGE_TAG} .

run-docker: build-docker
	docker run --name ${NAME} --volume=./log/:/app/log/ --net=host ${IMAGE_TAG}

deps:
	go mod download

tidy:
	go mod tidy