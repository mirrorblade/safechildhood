ifeq (${NAME},)
	NAME := "safechildhood"
endif

ifeq (${TAG},)
	TAG := "latest"
endif

IMAGE_TAG := ${NAME}:${TAG}

build: deps
	go build -o ./bin/safechildhood ./cmd/safechildhood/main.go

build-docker:
	docker build -t ${IMAGE_TAG} .
	
run: build
	./bin/safechildhood

run-docker: build-docker
	docker stop ${NAME} && docker rm ${NAME} || true

	docker run --name ${NAME} --volume=./log/:/app/log/ --net=host ${IMAGE_TAG}

deps:
	go mod download

tidy:
	go mod tidy

clean:
	rm -r ./bin

clean-docker:
	docker stop ${NAME} && docker rm ${NAME}
	docker rmi ${IMAGE_TAG}