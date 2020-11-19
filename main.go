package main

import (
	"flag"
	"log"
)

var configPath = "/etc/kimchi/config"

func main() {
	flag.StringVar(&configPath, "config", configPath, "configuration file")
	flag.Parse()

	srv := NewServer()
	if err := loadConfig(srv, configPath); err != nil {
		log.Fatal(err)
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
