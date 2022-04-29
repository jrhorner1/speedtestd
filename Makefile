.PHONY:= all build package 

version = 0.1.1
repo = jrhorner/speedtestd

all: build package

build: cmd/speedtestd/main.go
	mkdir -p bin/
	go build -o bin/speedtestd cmd/speedtestd/main.go 

buildx: cmd/speedtestd/main.go
	docker buildx build . -f ./build/package/Dockerfile.amd64 \
		--platform "linux/amd64"
		--tag $(repo):$(version)-amd64
		--push
	docker buildx build . -f ./build/package/Dockerfile.arm64 \
		--platform "linux/arm64/v8"
		--tag $(repo):$(version)-arm64
		--push
	docker buildx build . -f ./build/package/Dockerfile.armhf \
		--plaform "linux/arm/v7"
		--tag $(repo):$(version)-armhf
		--push
	docker manifest create $(repo):$(version) \
		--amend $(repo):$(version)-amd64 \
		--amend $(repo):$(version)-arm64 \
		--amend $(repo):$(version)-armhf
	docker manifest push $(repo):$(version)
	docker pull $(repo):$(version)
	docker tag $(repo):$(version) $(repo):latest
	docker push $(repo):latest

package: bin/speedtestd configs/speedtestd.yaml docs/README.txt LICENSE 
	tar cf build/package/speedtestd-v$(version).tgz docs/ bin/ configs/ LICENSE 

clean:
	rm -rf bin/ build/package/*.tgz
