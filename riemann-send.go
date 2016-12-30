package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/amir/raidman"
)


const (
	ENV_SERVER_CONFIG_PATH = "RIEMANN_SEND_SERVER_CONFIG"
	DEFAULT_SERVER_CONFIG_PATH = "/etc/riemann-send/server.json"
)

func readJsonFile(filename string, data interface{}) error {
	if bytes, err := ioutil.ReadFile(filename); err != nil {
		return err
	} else {
		if err := json.Unmarshal(bytes, data); err != nil {
			return err
		} else {
			return nil
		}

	}
}

func readJson(r io.Reader, data interface{}) error {
	dec := json.NewDecoder(r)
	return dec.Decode(data)
}

func getEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	} else {
		return defaultValue
	}
}

func loadServerConfig() ServerConfig {
	serverConfig := NewServerConfig()
	serverConfigPath := getEnv(ENV_SERVER_CONFIG_PATH, DEFAULT_SERVER_CONFIG_PATH)
	if err := readJsonFile(serverConfigPath, &serverConfig); err != nil && !os.IsNotExist(err) {
		log.Printf("Unable to read server config %s: %s", serverConfigPath, err)
	}
	return serverConfig
}

func loadEvent() raidman.Event {
	event := raidman.Event{}
	// TODO: read default event from json
	return event
}

func loadEventFromJson(jsonFile string, event *raidman.Event) error {
	if jsonFile == "-" {
		return readJson(os.Stdin, event)
	} else {
		return readJsonFile(jsonFile, event)
	}
}

func connectAndSend(serverConfig ServerConfig, event raidman.Event) {
	timeout := int64(time.Second) * serverConfig.Timeout
	if client, err := raidman.DialWithTimeout(serverConfig.Mode, serverConfig.Addr(), time.Duration(timeout)); err != nil {
		log.Panicf("Error connecting to riemann server %s: %s", serverConfig.Addr(), err)
	} else {
		if err := client.Send(&event); err != nil {
			log.Panicf("Error sending event %v to riemann server %s: %s", event, serverConfig.Addr(), err)
		}
	}
}

func main() {
	serverConfig := loadServerConfig()
	event := loadEvent()

	// riemann server flags
	flag.StringVar(&serverConfig.Host, "riemann-host", serverConfig.Host, "Riemann server")
	flag.UintVar(&serverConfig.Port, "riemann-port", serverConfig.Port, "Riemann port")
	flag.StringVar(&serverConfig.Mode, "riemann-mode", serverConfig.Mode, "Riemann connection type (tcp/tcp4/tcp6/udp/udp4/udp6)")
	flag.Int64Var(&serverConfig.Timeout, "riemann-timeout", serverConfig.Timeout, "Riemann connection timeout in seconds")

	// event flags
	jsonFile := flag.String("json", "", "Read event from file, - for stdin")
	flag.StringVar(&event.Host, "host", event.Host, "Event: host")
	flag.StringVar(&event.State, "state", event.State, "Event: state")
	flag.StringVar(&event.Service, "service", event.Service, "Event: service")
	flag.StringVar(&event.Description, "description", event.Description, "Event: description")
	var metric Metric // TODO: default metric
	flag.Var(&metric, "metric", "Event: metric")
	// TODO: ttl
	// TODO: time
	// TODO: tags
	// TODO: arguments
	flag.Parse()

	event.Metric, _ = metric.Parse()

	if *jsonFile != "" {
		if err := loadEventFromJson(*jsonFile, &event); err != nil {
			log.Panicf("Error reading event from json file %s: %s", *jsonFile, err)
		}
	}

	connectAndSend(serverConfig, event)
}
