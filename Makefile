.PHONY:= all build package 

version = 0.1.1
repo = jrhorner/speedtestd

all: build package

build: cmd/speedtestd/main.go
	mkdir -p bin/ \
	go build -o bin/speedtestd cmd/speedtestd/main.go 

buildx: cmd/speedtestd/main.go
	docker buildx build . -f ./build/package/Dockerfile \
		--platform linux/arm64 \
		--tag $(repo):$(version) \
		--tag $(repo):latest 

package: bin/speedtestd configs/speedtestd.yaml docs/README.txt LICENSE 
	tar cf build/package/speedtestd-v$(version).tgz docs/ bin/ configs/ LICENSE \

clean:
	rm -rf bin/ build/package/*.tgz
