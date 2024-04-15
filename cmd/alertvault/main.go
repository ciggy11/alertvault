package main

import (
	"encoding/json"
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/ciggy11/alertvault/pkg/config"
	"github.com/ciggy11/alertvault/pkg/server"
)

const (
	configFileOption = "config"
)

func main() {
	cfgFile := flag.String(configFileOption, "", "Path to the configuration file")
	flag.Parse()
	cfg, err := config.LoadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Unable to parse log level: %s", err)
	}
	log.SetLevel(lvl)
	cfgJSON, _ := json.Marshal(cfg)
	log.Infof("Configuration: %s", cfgJSON)
	s, err := server.New(*cfg)
	if err != nil {
		log.Fatalf("Unable to create server: %s", err)
		return
	}
	s.Start(cfg.HTTPListenAddress)
}
