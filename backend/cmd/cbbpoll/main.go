package main

import (
	"fmt"
	"github.com/r-cbb/cbbpoll/backend/internal/app"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("hello")

	server, err := app.NewServer()

	if err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
		log.Printf("Defaulting to port %s", port)
	}

	srv := &http.Server{
		Handler: server.Handler(),
		Addr:    fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}