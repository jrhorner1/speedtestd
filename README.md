# Speedtestd

This project is a Go wrapper for the Ookla Speedtest CLI that parses the output and sends the data to an InfluxDB server. 

![Grafana Dashboard](dashboard.png)

## Quick Start

You can quickly deploy this to kubernetes with the included helm chart.

```
helm repo add speedtestd https://jrhorner1.github.io/speedtestd/
helm install my-release speedtestd 
```
