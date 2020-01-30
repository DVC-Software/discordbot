package main

import (
	"fmt"
	"github.com/DVC-Software/discordbot/bot"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"time"
)

// Golbal
var dev_port string = "0.0.0.0:7080"
var test_port string = "0.0.0.0:7070"
var url_prefix string = "/"

func getPortFromEnv() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" {
		fmt.Print("development")
		return dev_port
	} else if env == "test" {
		fmt.Print("test")
		return test_port
	} else {
		fmt.Print("development")
		return dev_port
	}
}

func main() {
	r := mux.NewRouter()
	http.Handle(url_prefix, r)
	port := getPortFromEnv()
	srv := &http.Server{
		Handler: r,
		Addr:    port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	fmt.Println("starting server...")
	bot.Start()

	if err := srv.ListenAndServe(); err != nil {
		panic(err.Error())
	}
	defer srv.Close()
}
