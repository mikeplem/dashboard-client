package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
)

var configFile *string
var appConfig tomlConfig

type tomlConfig struct {
	Chrome chromeconfig `toml:"chrome"`
	Consul consulconfig `toml:"consul"`
	Delay  delayconfig  `tomle:"delay"`
}

type chromeconfig struct {
	Host string
	Port int
}

type consulconfig struct {
	Address    string
	Scheme     string
	Datacenter string
	Action     string
	URL        string
}

type delayconfig struct {
	Interval time.Duration
}

// hold onto the currently running URL
var runningURL string

// ChromeConnString holds the chrome address and port
// it is used in browser.go
var ChromeConnString string

// =============================

func init() {
	configFile = flag.String("conf", "", "Config file chromium, Consul, delay interval.")

	flag.Parse()

	if _, err := toml.DecodeFile(*configFile, &appConfig); err != nil {
		log.Fatal(err)
	}

	ChromeConnString = fmt.Sprintf("%s:%d", appConfig.Chrome.Host, appConfig.Chrome.Port)

	log.Println("Chrome Connection: ", ChromeConnString)
	log.Println("Consul Address: ", appConfig.Consul.Address)
	log.Println("Consul Scheme: ", appConfig.Consul.Scheme)
	log.Println("Consul Datacenter: ", appConfig.Consul.Datacenter)
	log.Println("Consul Action Path: ", appConfig.Consul.Action)
	log.Println("Consul URL Path: ", appConfig.Consul.URL)
	log.Println("Loop Delay: ", appConfig.Delay.Interval*time.Millisecond)
}

func main() {

	consulConfig := api.DefaultConfig()
	consulConfig.Address = appConfig.Consul.Address
	consulConfig.Scheme = appConfig.Consul.Scheme
	consulConfig.Datacenter = appConfig.Consul.Datacenter

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Print("Error creating consul client", err)
	}

	kv := client.KV()

	// using []byte because that is the format
	// the data is stored on Consul
	c := make(chan []byte)

	for {
		go ReadConsulPath(c, kv)
		chanValue := <-c

		s := strings.Split(string(chanValue), ",")

		action, url := s[0], s[1]

		switch string(action) {
		case "open":
			if runningURL != url {
				log.Println("Open URL: ", url)
				OpenURLInBrowser(url)
				runningURL = url
			}
		case "reload":
			log.Println("Reload browser")
			ReloadBrowser()
			WriteConsulPath(kv)
		default:
			//log.Println("switch default")
		}

		time.Sleep(appConfig.Delay.Interval * time.Millisecond)
	}
}
