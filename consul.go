package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// ReadConsulPath will look for values at specific paths
// and return them via a channel
func ReadConsulPath(c chan []byte, kv *api.KV) {

	actionExists, _, err := kv.Keys(appConfig.Consul.Action, "/", nil)
	if err != nil {
		log.Print("Error finding action key from Consul", err)
	}

	urlExists, _, err := kv.Keys(appConfig.Consul.URL, "/", nil)
	if err != nil {
		log.Print("Error finding url key from Consul", err)
	}

	if len(actionExists) > 0 && len(urlExists) > 0 {

		action, _, err := kv.Get(appConfig.Consul.Action, nil)
		if err != nil {
			log.Print("Error getting action value from Consul", err)
		}

		url, _, err := kv.Get(appConfig.Consul.URL, nil)
		if err != nil {
			log.Print("Error getting url value from Consul", err)
		}

		returnString := fmt.Sprintf("%s,%s", action.Value, url.Value)
		returnByte := []byte(returnString)

		c <- returnByte
	}
}

// WriteConsulPath writes the value open to the action path
// this is so the for loop in main does not continue trying
// to reload the browser
func WriteConsulPath(kv *api.KV) {

	p := &api.KVPair{Key: appConfig.Consul.Action, Value: []byte("open")}
	_, err := kv.Put(p, nil)
	if err != nil {
		log.Println(err)
	}
}
