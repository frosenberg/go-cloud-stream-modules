DOCKER_REPO?=frosenberg
IMAGE_NAME?=gcs-module-httpclient-processor
IMAGE_TAG?=1.0.3

default: builddocker

buildgo:
	go get github.com/tools/godep
	CGO_ENABLED=0 GOOS=linux godep go build -ldflags "-s" -a -installsuffix cgo -o main httpclient-processor.go

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

refresh-deps:
	rm -rf Godeps
	godep save

run-kafka-sample:
	godep go run httpclient-processor.go -v --spring.cloud.stream.binder.kafka.brokers=localhost:9092 --spring.cloud.stream.bindings.output.destination=topic:httpclient_out --spring.cloud.stream.bindings.input.destination=topic:bridge_out --url=http://google.com

.PHONY: clean run-sample


