package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
)

// hold onto the currently running URL
var runningURL string

var actionPath = "foo/action"
var urlPath = "foo/url"
var delay time.Duration = 1000

func readConsulPath(c chan []byte, kv *api.KV) {

	actionExists, _, err := kv.Keys(actionPath, "/", nil)
	if err != nil {
		log.Print("Error getting Keys from Consul", err)
	}

	urlExists, _, err := kv.Keys(urlPath, "/", nil)
	if err != nil {
		log.Print("Error getting Keys from Consul", err)
	}

	if len(actionExists) > 0 && len(urlExists) > 0 {

		action, _, err := kv.Get(actionPath, nil)
		if err != nil {
			log.Print("Error getting action from Consul", err)
		}

		url, _, err := kv.Get(urlPath, nil)
		if err != nil {
			log.Print("Error getting url from Consul", err)
		}

		returnString := fmt.Sprintf("%s,%s", action.Value, url.Value)
		returnByte := []byte(returnString)

		c <- returnByte
	}
}

func main() {

	config := api.DefaultConfig()
	config.Address = "127.0.0.1:8500"
	config.Scheme = "http"
	config.Datacenter = "datacenter1"

	client, err := api.NewClient(config)
	if err != nil {
		log.Print("Error creating client", err)
	}

	kv := client.KV()

	c := make(chan []byte)

	for {
		time.Sleep(delay * time.Millisecond)
		go readConsulPath(c, kv)
		chanValue := <-c
		s := strings.Split(string(chanValue), ",")

		action, url := s[0], s[1]

		switch string(action) {
		case "open":
			if runningURL != url {
				log.Println("Open URL: ", url)
				// chromium open call happens here

				// save the new URL to the global var
				runningURL = url
			}
		case "reload":
			log.Println("Reload browser")
			// chromium reload call happens here

			// once the reload happens swich the action to open
			// so that we don't stay in this loop
			p := &api.KVPair{Key: "foo/action", Value: []byte("open")}
			_, err = kv.Put(p, nil)
			if err != nil {
				log.Println(err)
			}
		default:
			//log.Println("switch default")
		}
	}
}
