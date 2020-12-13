package main

import (
	"log"
	"net/http"
	"search/config"
	"search/server/connectors"
	"search/server/routes"
	"time"
)

var cfg *config.Config

func init() {
	cfg, _ = config.LoadConfig()
}

func main() {
	for dbConnResult := false; !dbConnResult; {
		log.Println("\033[32m ⇨ Initializing database hang on..\033[0m")
		dbConnResult = connectors.Initialize(cfg)
		if !dbConnResult {
			log.Printf("Unable to initialize database. Sleeping for %d seconds...", 10)
			time.Sleep(time.Duration(10) * time.Second)
		}
	}
	log.Println("\033[32m ⇨ Initializing router almost done..\033[0m")
	routes.InitializeRouter()
	log.Println("\033[32m ⇨ http server started at " + cfg.Server.Host + ":" + cfg.Server.Port + "\033[0m")
	http.ListenAndServe(":"+cfg.Server.Port, routes.RouterInstance())
}
