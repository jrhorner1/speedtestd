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
    Ping struct {
		Jitter float64
		Latency float64
	} `json:"ping"`
    Download struct {
		Bandwidth int
		Bytes int
		Elapsed int
	} `json:"download"`
    Upload struct {
		Bandwidth int
		Bytes int
		Elapsed int
	} `json:"upload"`
    PacketLoss float64
    Isp string
    Interface struct {
		InternalIp string
		Name string
		MacAddr string
		IsVpn bool
		ExternalIp string
	} `json:"interface"`
    Server struct {
		Id int
		Name string
		Location string
		Country string
		Host string
		Port int
		Ip string
	} `json:"server"`
    Result struct {
		Id string
		Url string
	} `json:"result"`
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