TAG ?= stable
PUSH ?= false
IMAGE_REGISTRY = quay.io/adarsh0728/test/sink-log:${TAG}
ARCHITECTURES = arm64

.PHONY: build
build:
	for arch in $(ARCHITECTURES); do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$${arch} go build -v -o ./dist/log-example-$${arch} main.go; \
	done

.PHONY: image-push
image-push: build
	docker buildx build -t ${IMAGE_REGISTRY} --platform linux/arm64 --target log . --push

.PHONY: image
image: build
	docker build -t ${IMAGE_REGISTRY} --target log .
	@if [ "$(PUSH)" = "true" ]; then docker push ${IMAGE_REGISTRY}; fi

clean:
	-rm -rf ./dist