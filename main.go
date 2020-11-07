package main

import (
	"flag"
	"log"

	"git.sr.ht/~emersion/go-scfg"
)

var configPath = "/etc/kimchi/config"

func main() {
	flag.StringVar(&configPath, "config", configPath, "configuration file")
	flag.Parse()

	cfg, err := scfg.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}

	srv := NewServer()
	for _, dir := range cfg {
		if err := parseSite(srv, dir); err != nil {
			log.Fatal(err)
		}
	}

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	select {}
}
