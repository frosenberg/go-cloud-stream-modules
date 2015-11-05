DOCKER_REPO?=frosenberg
IMAGE_NAME?=gcs-module-log-sink
IMAGE_TAG?=1.0.0

default: builddocker

buildgo:
	go get github.com/tools/godep
	CGO_ENABLED=0 GOOS=linux godep go build -ldflags "-s" -a -installsuffix cgo -o main *.go

builddocker:
	docker build -t ${DOCKER_REPO}/${IMAGE_NAME}:${IMAGE_TAG} -f ./Dockerfile.build .
	docker run -t ${DOCKER_REPO}/${IMAGE_NAME}:${IMAGE_TAG} /bin/true
	docker cp `docker ps -q -n=1`:/main .
	chmod 755 ./main
	docker build --rm=true --tag=${DOCKER_REPO}/${IMAGE_NAME}:${IMAGE_TAG} -f Dockerfile.static .

run: builddocker
	docker run ${DOCKER_REPO}/${IMAGE_NAME}:${IMAGE_TAG}

push:
	docker push ${DOCKER_REPO}/${IMAGE_NAME}:${IMAGE_TAG}

clean:
	rm main

.PHONY: clean
