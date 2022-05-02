package main

import (
	"strconv"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/jrhorner1/ookla-speedtest/pkg/speedtest"
	log "github.com/sirupsen/logrus"
)

func influxdbConnect(results *speedtest.Speedtest, config *Config) {
	influxdb_protocol := "http"
	influxdb_server := config.Influxdb.Address
	influxdb_port := config.Influxdb.Port
	influxdb_url := influxdb_protocol + "://" + influxdb_server + ":" + strconv.Itoa(influxdb_port)
	influxdb_user := config.Influxdb.Username
	influxdb_pass := config.Influxdb.Password
	influxdb_token := influxdb_user + ":" + influxdb_pass
	influxdb_org := ""
	influxdb_database := config.Influxdb.Database

	log.Info("Connecting to influxdb server: " + influxdb_url)
	client := influxdb2.NewClient(influxdb_url, influxdb_token)

	writeAPI := client.WriteAPI(influxdb_org, influxdb_database)
	errorsCh := writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Error("write error:", err.Error())
		}
	}()

	tags := map[string]string{
		"serverId":       strconv.Itoa(results.Server.Id),
		"serverName":     results.Server.Name,
		"serverLocation": results.Server.Location,
		"serverCountry":  results.Server.Country,
		"serverHost":     results.Server.Host,
		"serverPort":     strconv.Itoa(results.Server.Port),
		"serverIp":       results.Server.Ip,
		"isp":            results.Isp}

	ping := influxdb2.NewPoint(
		"ping", // measurement
		tags,
		map[string]interface{}{ // fields
			"jitter":  results.Ping.Jitter,
			"latency": results.Ping.Latency},
		results.Timestamp)
	log.Debug("Writing ping measurements")
	writeAPI.WritePoint(ping)

	download := influxdb2.NewPoint(
		"download", // measurement
		tags,
		map[string]interface{}{
			"bandwidth": results.Download.Bandwidth * 8, // Value is in bytes, converting to bits
			"bytes":     results.Download.Bytes,
			"elapsed":   results.Download.Elapsed},
		results.Timestamp)
	log.Debug("Writing download measurements")
	writeAPI.WritePoint(download)

	upload := influxdb2.NewPoint(
		"upload", // measurement
		tags,
		map[string]interface{}{
			"bandwidth": results.Upload.Bandwidth * 8, // Value is in bytes, converting to bits
			"bytes":     results.Upload.Bytes,
			"elapsed":   results.Upload.Elapsed},
		results.Timestamp)
	log.Debug("Writing upload measurements")
	writeAPI.WritePoint(upload)

	packet := influxdb2.NewPoint(
		"packet", // measurement
		tags,
		map[string]interface{}{ // fields
			"loss": results.PacketLoss},
		results.Timestamp)
	log.Debug("Writing packet measurements")
	writeAPI.WritePoint(packet)

	clientData := influxdb2.NewPoint(
		"client",
		tags,
		map[string]interface{}{
			"internalIp":       results.Interface.InternalIp,
			"interfaceName":    results.Interface.Name,
			"interfaceMacAddr": results.Interface.MacAddr,
			"isVpn":            strconv.FormatBool(results.Interface.IsVpn),
			"externalIp":       results.Interface.ExternalIp,
			"isp":              results.Isp},
		results.Timestamp)
	log.Debug("Writing client info")
	writeAPI.WritePoint(clientData)

	serverData := influxdb2.NewPoint(
		"server",
		tags,
		map[string]interface{}{
			"serverId":       strconv.Itoa(results.Server.Id),
			"serverName":     results.Server.Name,
			"serverLocation": results.Server.Location,
			"serverCountry":  results.Server.Country,
			"serverHost":     results.Server.Host,
			"serverPort":     strconv.Itoa(results.Server.Port),
			"serverIp":       results.Server.Ip},
		results.Timestamp)
	log.Debug("Writing server info")
	writeAPI.WritePoint(serverData)

	result := influxdb2.NewPoint(
		"result",
		tags,
		map[string]interface{}{
			"id":  results.Result.Id,
			"url": results.Result.Url},
		results.Timestamp)
	log.Debug("Writing result info")
	writeAPI.WritePoint(result)

	log.Info("Flushing writes from the buffer")
	writeAPI.Flush()
	log.Info("Closing the influxdb client connection")
	client.Close()
}
