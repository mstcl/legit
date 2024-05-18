package main

import (
	"flag"
	"fmt"
	"log"
	"embed"
	"net/http"

	"github.com/mstcl/legit/config"
	"github.com/mstcl/legit/routes"
)

//go:embed templates/*
var fs embed.FS

func main() {
	var cfg string
	flag.StringVar(&cfg, "config", "./config.yaml", "path to config file")
	flag.Parse()

	c, err := config.Read(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := UnveilPaths([]string{
		c.Repo.ScanPath,
	},
		"r"); err != nil {
		log.Fatalf("unveil: %s", err)
	}

	routes.Templates = fs
	mux := routes.Handlers(c)
	addr := fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
	log.Println("starting server on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
