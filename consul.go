package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// ReadConsulPath will look for values at specific paths
// and return them via a channel
func ReadConsulPath(c chan []byte, kv *api.KV, urlPath string) {

	actionExists, _, err := kv.Keys(appConfig.Consul.Action, "/", nil)
	if err != nil {
		log.Print("Error finding action key from Consul", err)
	}

	urlExists, _, err := kv.Keys(urlPath, "/", nil)
	if err != nil {
		log.Print("Error finding urlPath key from Consul", err)
	}

	if len(actionExists) > 0 && len(urlExists) > 0 {

		action, _, err := kv.Get(appConfig.Consul.Action, nil)
		if err != nil {
			log.Print("Error getting action value from Consul", err)
		}

		url, _, err := kv.Get(urlPath, nil)
		if err != nil {
			log.Print("Error getting urlPath value from Consul", err)
		}

		returnString := fmt.Sprintf("%s,%s", action.Value, url.Value)
		returnByte := []byte(returnString)

		c <- returnByte
	}
}

// WriteConsulPath writes 'value' to the Consul 'path'
func WriteConsulPath(kv *api.KV, path string, value string) {

	p := &api.KVPair{Key: path, Value: []byte(value)}
	_, err := kv.Put(p, nil)
	if err != nil {
		log.Println(err)
	}
}
