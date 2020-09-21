# Ookla Speedtest to InfluxDB

This program aims to run the Ookla Speedtest CLI and export the results to an InfluxDB server for storage/visualization

## Build for ARM  

### Linux  

```bash
GOOS=linux GOARCH=arm go build -o bin/speedtest2influx cmd/speedtest2influx/main.go
```

### Windows  

```powershell
set GOOS=linux
set GOARCH=arm
go build -o bin/speedtest2influx cmd/speedtest2influx/main.go
```

## Docker  

```
docker buildx build --platform linux/arm64 . -t speedtest2influx:latest
```

## Helm  
