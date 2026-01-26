package main

import (
	"log"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
