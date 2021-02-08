# Ookla Speedtest to InfluxDB

This program aims to run the Ookla Speedtest CLI and export the results to an InfluxDB server for storage/visualization

## Build

Prequisites:
- go
- docker

To build the application and docker container, simply run:
```sh
make
```

Clean up your build environment with: 
```sh
make clean
```

## Deployment  

### Docker Compose

```sh
cd deployments/compose
docker-compose up -d
```

### Helm

```sh
cd deployments/helm
helm install speedtest -f my-values.yaml -n speedtest ./ookla-speedtest --create-namespace
```

### Terraform

```sh
cd deployments/terraform
terraform init
terraform plan -out=plan.out
terraform apply "plan.out"
```
