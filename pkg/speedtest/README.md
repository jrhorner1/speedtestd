# Speedtest.net Speedtest CLI wrapper package for GoLang

This package requires the appropriate [Speedtest CLI](https://www.speedtest.net/apps/cli) binary for your system to function.  

## Examples

Here's a basic example:

```go
package main

import (
    "fmt"
    "github.com/jrhorner1/ookla-speedtest/pkg/speedtest"
)

func main() {
    results := speedtest.Run()
    /* the bandwidth results are returned as bytes */
    fmt.Println("Download Speed: ", (results.download.bandwidth * 8)/1000000, "mbps")
    fmt.Println("Upload Speed: ", (results.upload.bandwidth * 8)/1000000, "mbps")
    fmt.Println("Latency: ", results.ping.latency, "ms")
}
```

## Sample output data

This is a sample of the JSON output from the speedtest CLI directly. All of this data is collected for later use in your application.  

```json
{
    "type": "result",
    "timestamp": "2020-09-15T22:32:31Z", # RFC3339
    "ping": {
        "jitter": 1.3560000000000001, # milliseconds
        "latency": 9.9130000000000003 # milliseconds
    },
    "download": {
        "bandwidth": 4052019, # bytes per second (Bps)
        "bytes": 34195752,
        "elapsed": 8915 # milliseconds
    },
    "upload": {
        "bandwidth": 1387153, # bytes per second (Bps)
        "bytes": 16350810,
        "elapsed": 10911 # milliseconds
    },
    "packetLoss": 0.69444444444444442,
    "isp": "Example ISP",
    "interface": {
        "internalIp": "192.168.1.100",
        "name": "eth0",
        "macAddr": "FF:FF:FF:FF:FF:FF",
        "isVpn": false,
        "externalIp": "111.1.111.1"
    },
    "server": {
        "id": 00001,
        "name": "Example ISP",
        "location": "Somewhere, DC",
        "country": "United States",
        "host": "speedtest.example.com",
        "port": 8080,
        "ip": "1.111.1.111"
    },
    "result": {
        "id": "f9fcc30a-b2ef-625c-1c0b-957cdef95548",
        "url": "https://www.speedtest.net/result/c/f9fcc30a-b2ef-625c-1c0b-957cdef95548"
    }
}
```
