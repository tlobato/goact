package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	server.ApiRoutes()

	fs := http.FileServer(http.Dir("./client/dist/"))
	http.Handle("/", fs)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
