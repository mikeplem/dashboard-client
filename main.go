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
	NewURL     string
	RunningURL string
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
	configFile = flag.String("conf", "", "Config for Chromium, Consul, delay interval. If not provided, config.toml in this directory will be read by default.")

	flag.Parse()

	if *configFile == "" {
		*configFile = "config.toml"
	}

	if _, err := toml.DecodeFile(*configFile, &appConfig); err != nil {
		log.Fatal(err)
	}

	ChromeConnString = fmt.Sprintf("%s:%d", appConfig.Chrome.Host, appConfig.Chrome.Port)

	log.Println("Chrome Connection: ", ChromeConnString)
	log.Println("Consul Address: ", appConfig.Consul.Address)
	log.Println("Consul Scheme: ", appConfig.Consul.Scheme)
	log.Println("Consul Datacenter: ", appConfig.Consul.Datacenter)
	log.Println("Consul Action Path: ", appConfig.Consul.Action)
	log.Println("Consul New URL Path: ", appConfig.Consul.NewURL)
	log.Println("Consul Running URL Path: ", appConfig.Consul.RunningURL)
	log.Println("Loop Delay: ", appConfig.Delay.Interval*time.Millisecond)
}

func main() {

	consulConfig := api.DefaultConfig()
	consulConfig.Address = appConfig.Consul.Address
	consulConfig.Scheme = appConfig.Consul.Scheme
	consulConfig.Datacenter = appConfig.Consul.Datacenter

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal("Error creating consul client", err)
	}

	kv := client.KV()

	// using []byte because that is the format
	// the data is stored in Consul
	c := make(chan []byte)

	for {
		// ReadConsulPath also looks at the action path
		// since that is always going to be used to determine
		// if a path should be opened or reloaded.
		go ReadConsulPath(c, kv, appConfig.Consul.NewURL)
		chanValue := <-c

		s := strings.Split(string(chanValue), ",")
		action, url := s[0], s[1]

		switch string(action) {
		case "open":
			if runningURL != url {
				log.Println("Open URL: ", url)
				OpenURLInBrowser(url)

				// writing to Consul in case we need that data from
				// the admin interface but using the local variable
				// since that will be faster to compare against.
				WriteConsulPath(kv, appConfig.Consul.RunningURL, url)
				runningURL = url
			}
		case "reload":
			log.Println("Reload browser")
			ReloadBrowser()
			// setting the action path to open so that reload
			// does not cause the application to stay in the
			// reload loop.
			WriteConsulPath(kv, appConfig.Consul.Action, "open")
		default:
			//log.Println("switch default")
		}

		time.Sleep(appConfig.Delay.Interval * time.Millisecond)
	}
}
