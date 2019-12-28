package main

import (
	"gwahangmi-backend/server"
	"log"
)

func main() {
	s, err := server.New()
	if err != nil {
		log.Fatal(err)
	}
	s.Run(":3000")
}
