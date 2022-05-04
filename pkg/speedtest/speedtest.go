package speedtest

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

type Results struct {
	Type      string
	Timestamp time.Time
	Ping      struct {
		Jitter  float64
		Latency float64
	} `json:"ping"`
	Download struct {
		Bandwidth int
		Bytes     int
		Elapsed   int
	} `json:"download"`
	Upload struct {
		Bandwidth int
		Bytes     int
		Elapsed   int
	} `json:"upload"`
	PacketLoss float64
	Isp        string
	Interface  struct {
		InternalIp string
		Name       string
		MacAddr    string
		IsVpn      bool
		ExternalIp string
	} `json:"interface"`
	Server struct {
		Id       int
		Name     string
		Location string
		Country  string
		Host     string
		Port     int
		Ip       string
	} `json:"server"`
	Result struct {
		Id  string
		Url string
	} `json:"result"`
}

func Run(retries int) *Results {
	args := [2]string{"", ""}
	return speedtest(args, retries)
}

func RunWithServerId(serverid, retries int) *Results {
	args := [2]string{"-s", strconv.Itoa(serverid)}
	return speedtest(args, retries)
}

func RunWithHost(host string, retries int) *Results {
	args := [2]string{"-o", host}
	return speedtest(args, retries)
}

func speedtest(testargs [2]string, retries int) *Results {
	var results Results
	for i := 0; i < retries; i++ {
		var test *exec.Cmd
		if testargs[0] != "" && testargs[1] != "" {
			test = exec.Command("speedtest", "-f", "json", testargs[0], testargs[1])
		} else {
			test = exec.Command("speedtest", "-f", "json", "-vv")
		}
		var out bytes.Buffer
		test.Stdout = &out

		log.Info("Running Speedtest...")
		err := test.Run()
		if err != nil {
			log.Error("speedtest-cli error ", err.Error())
			continue
		}
		log.Info("Speedtest completed.")

		jsonOutput := bytes.Split(out.Bytes(), []byte("\n"))
		for _, json := range jsonOutput {
			if string(json) == "" {
				continue
			}
			log.Info(string(json))
		}
		err = json.Unmarshal(jsonOutput[len(jsonOutput)-2], &results)
		if err != nil {
			log.Error("JSON unmarshal error", err.Error())
			continue
		}
		if results.Download.Bandwidth == 0 && results.Upload.Bandwidth == 0 {
			continue
		}
		break
	}
	return (&results)
}
