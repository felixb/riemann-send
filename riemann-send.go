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
	"fmt"
	"strings"
)

const (
	ENV_SERVER_CONFIG_PATH = "RIEMANN_SEND_SERVER_CONFIG"
	ENV_EVENT_CONFIG_PATH = "RIEMANN_SEND_EVENT_CONFIG"
	DEFAULT_SERVER_CONFIG_PATH = "/etc/riemann-send/server.json"
	DEFAULT_EVENT_CONFIG_PATH = "/etc/riemann-send/event.json"
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
	eventConfigPath := getEnv(ENV_EVENT_CONFIG_PATH, DEFAULT_EVENT_CONFIG_PATH)
	if err := readJsonFile(eventConfigPath, &event); err != nil && !os.IsNotExist(err) {
		log.Printf("Unable to read event config %s: %s", eventConfigPath, err)
	}
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

func printEvent(event raidman.Event) {
	if bytes, err := json.MarshalIndent(event, "", "  "); err != nil {
		log.Println("Error writing event to stdout", err)
	} else {
		os.Stdout.Write(bytes)
		fmt.Println("")
	}
}

func parseArgs(serverConfig *ServerConfig, event *raidman.Event) {
	// riemann server flags
	riemannHost := flag.String("riemann-host", serverConfig.Host, "Riemann server")
	riemannPort := flag.Uint("riemann-port", serverConfig.Port, "Riemann port")
	riemannMode := flag.String("riemann-mode", serverConfig.Mode, "Riemann connection type (tcp/tcp4/tcp6/udp/udp4/udp6)")
	riemannTimeout := flag.Int64("riemann-timeout", serverConfig.Timeout, "Riemann connection timeout in seconds")

	// event flags
	jsonFile := flag.String("json", "", "Read event from file, - for stdin")
	host := flag.String("host", event.Host, "Event: host")
	state := flag.String("state", event.State, "Event: state")
	service := flag.String("service", event.Service, "Event: service")
	description := flag.String("description", event.Description, "Event: description")
	eventTime := flag.Int64("time", event.Time, "Event: time")
	ttl := flag.Float64("ttl", float64(event.Ttl), "Event: TTL")
	tags := flag.String("tags", "", "Event: tags, comma separated")
	attributes := flag.String("attributes", "", "Event: attributes, format 'key:value[,key:value]...'")
	metric := NewMetric(event.Metric)
	flag.Var(&metric, "metric", "Event: metric")

	flag.Parse()

	serverConfig.Host = *riemannHost
	serverConfig.Port = *riemannPort
	serverConfig.Timeout = *riemannTimeout
	serverConfig.Mode = *riemannMode

	event.Host = *host
	event.State = *state
	event.Service = *service
	event.Description = *description
	event.Time = *eventTime
	event.Ttl = float32(*ttl)
	event.Metric, _ = metric.Parse()
	event.Tags = append(event.Tags, strings.Split(*tags, ",")...)
	for _, e := range strings.Split(*attributes, ",") {
		kv := strings.SplitN(e, ":", 2)
		if len(kv) != 2 {
			log.Panicf("Error parsing attributes, format: key:value[,key:value]...")
		}
		event.Attributes[kv[0]] = kv[1]
	}

	if *jsonFile != "" {
		if err := loadEventFromJson(*jsonFile, event); err != nil {
			log.Panicf("Error reading event from json file %s: %s", *jsonFile, err)
		}
	}
}

func main() {
	serverConfig := loadServerConfig()
	event := loadEvent()

	// generic flags
	quiet := flag.Bool("quiet", false, "suppress output of sent events")
	parseArgs(&serverConfig, &event)

	connectAndSend(serverConfig, event)
	if !*quiet {
		printEvent(event)
	}
}
