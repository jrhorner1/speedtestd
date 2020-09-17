package speedtest

import (
	"time"
	"os/exec"
	"bytes"
	log "github.com/sirupsen/logrus"
	"encoding/json"
	"strconv"
)

type Speedtest struct {
    Type string
    Timestamp time.Time
    Ping Ping
    Download DataRate
    Upload DataRate
    PacketLoss float64
    Isp string
    Interface Interface
    Server Server
    Result Result
}

type Ping struct {
	Jitter float64
    Latency float64
}

type DataRate struct {
	Bandwidth int
	Bytes int
	Elapsed int
}

type Interface struct {
	InternalIp string
	Name string
	MacAddr string
	IsVpn bool
	ExternalIp string
}

type Server struct {
	Id int
	Name string
	Location string
	Country string
	Host string
	Port int
	Ip string
}

type Result struct {
	Id string
	Url string
}

func Run() *Speedtest {
	args := [2]string{"", ""}
	return speedtest(args)
}

func RunWithServerId(serverid int) *Speedtest {
	args := [2]string{"-s", strconv.Itoa(serverid)}
	return speedtest(args)
}

func RunWithHost(host string) *Speedtest {
	args := [2]string{"-o", host}
	return speedtest(args)
}

func speedtest(testargs [2]string) *Speedtest {
	var test *exec.Cmd
	if testargs[0] != "" && testargs[1] != "" {
		test = exec.Command("speedtest", "-f", "json", testargs[0], testargs[1])
	} else {
		test = exec.Command("speedtest", "-f", "json")
	}
	var out bytes.Buffer
	test.Stdout = &out

	log.Info("Running Speedtest...")
	err := test.Run()
	if err != nil  {
		log.Error("speedtest-cli error", err.Error())
	}
	log.Info("Speedtest completed..")

	jsonOutput := out.Bytes()
	log.Info(string(jsonOutput))

	var results Speedtest
	err = json.Unmarshal(jsonOutput, &results)
	if err != nil {
		log.Error("JSON unmarshal error", err.Error())
	}
	
	return(&results)
}