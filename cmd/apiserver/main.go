package main

import (
	apiserver "github.com/Andronovdima/job-avito-assignment/internal/app"
	"log"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
