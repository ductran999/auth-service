package main

import (
	"auth-service/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalln("failed to run app:", err.Error())
	}
}
