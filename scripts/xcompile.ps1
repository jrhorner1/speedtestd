$version = "0.1.0"

docker buildx build . `
	--platform linux/amd64,linux/arm64,linux/arm `
	-f ./build/package/Dockerfile `
	--tag jrhorner/ookla-speedtest:$version `
	--tag jrhorner/ookla-speedtest:latest `
	--push