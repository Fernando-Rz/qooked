package main

import (
	"log"
	"os"
	"qooked/internal/http/server"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	
	if err != nil {
		log.Fatalf("\n%s\nError loading .env file: %v\n%s\n", strings.Repeat("-", 100), err, strings.Repeat("-", 100))
	}

    environment, exists := os.LookupEnv("QOOKED_ENV")
    
    if !exists {
        environment = "local"
    }

	server, err := server.NewServer(environment)

	if err != nil {
		log.Fatalf("\n%s\nFailed to initialize server: %v\n%s\n", strings.Repeat("-", 100), err, strings.Repeat("-", 100))
	}

	if err := server.Run(); err != nil {
		log.Fatalf("\n%s\nServer failed to run: %v\n%s\n", strings.Repeat("-", 100), err, strings.Repeat("-", 100))
	}
}
