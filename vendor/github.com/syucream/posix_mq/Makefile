.PHONY: docker
docker:
	docker build -f Dockerfile-alpine -t posix_mq_alpine .
	docker build -f Dockerfile-stretch -t posix_mq_stretch .

.PHONY: build
build:
	go build example/exec/sender.go
	go build example/exec/receiver.go

