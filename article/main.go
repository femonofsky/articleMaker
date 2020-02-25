package main

import (
	"flag"
	"fmt"
	"github.com/femonofsky/articleMaker/article/config"
	"github.com/femonofsky/articleMaker/article/controller"
	"github.com/femonofsky/articleMaker/article/model"
	"log"
	"net/http"
	"os"
)

func main() {

	configPath := flag.String("config", "./config/config.json", "path of the config file")

	flag.Parse()

	// Initialize Logger
	logger := log.New(os.Stdout, "article-api ", log.LstdFlags)

	logger.Println("Starting the application...")

	// load Config from file
	cfg, err := config.FromFile(*configPath)
	if err != nil {
		logger.Fatal("file not found" + err.Error())
	}
	// Initialize Database
	DB, err := model.New(cfg)
	if err != nil {
		logger.Fatalf("could not initialize DB connection : %v ", err.Error())
	}

	// Migrate all Table into DB if it doesn't exist
	DB = model.Migrate(DB)

	defer DB.Close()

	// Register all Controllers and its routes
	sm := controller.New(logger)

	// listens on the TCP network address addr
	ADDR := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	// start http server
	logger.Fatal(http.ListenAndServe(ADDR, sm))
}
