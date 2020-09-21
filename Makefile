.PHONY:= all build package

all: build package

build: cmd/speedtest2influx/main.go
	./scripts/build.sh

package: bin/speedtest2influx configs/speedtest2influx.yaml docs/README.txt LICENSE 
	./scripts/package.sh

clean:
	rm -rf bin/ speedtest2influx-v*.tgz
