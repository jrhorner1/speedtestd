.PHONY:= all build package 

version = 0.1.1
repo = jrhorner/ookla-speedtest

all: build

build: cmd/speedtest2influx/main.go
	docker build . -f ./build/package/Dockerfile \
		--tag $(repo):$(version) \
		--tag $(repo):latest

xbuild: cmd/speedtest2influx/main.go
	docker buildx build . -f ./build/package/Dockerfile \
		--platform linux/amd64,linux/arm64,linux/arm \
		--tag $(repo):$(version) \
		--tag $(repo):latest 

package: bin/speedtest2influx configs/speedtest2influx.yaml docs/README.txt LICENSE 
	tar cf build/package/speedtest2influx-v$(version).tgz docs/ bin/ configs/ LICENSE \

clean:
	rm -rf bin/ build/package/*.tgz
