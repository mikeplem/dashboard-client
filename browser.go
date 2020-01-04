package main

import (
	"log"

	"github.com/raff/godet"
)

// ReloadBrowser will reload the current page open in the browser
func ReloadBrowser() {

	remote, err := godet.Connect(ChromeConnString, false)
	if err != nil {
		log.Println("cannot connect to Chrome instance:", err)
		return
	}

	defer remote.Close()

	log.Print("Reloading browser page")
	err = remote.Reload()
	if err != nil {
		log.Print(err)
		return
	}
}

// OpenURLInBrowser will pass 'url' to the browser
func OpenURLInBrowser(url string) {

	remote, err := godet.Connect(ChromeConnString, false)
	if err != nil {
		log.Println("cannot connect to Chrome instance:", err)
		return
	}

	defer remote.Close()

	log.Printf("Requested to open %s\n", url)
	_, _ = remote.Navigate(url)

}
