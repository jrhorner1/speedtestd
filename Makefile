# .PHONY:= all build package docker
.PHONY:= all docker

# all: build package docker
all: docker

# build: cmd/speedtest2influx/main.go
# 	./scripts/build.sh

# package: bin/speedtest2influx configs/speedtest2influx.yaml docs/README.txt LICENSE 
# 	./scripts/package.sh

docker: 
	./scripts/docker.sh

# clean:
# 	rm -rf bin/ build/package/*.tgz
