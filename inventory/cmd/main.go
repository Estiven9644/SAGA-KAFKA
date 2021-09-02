package main

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
	"leal.co/inventory/internal/pkg/database"
	"leal.co/inventory/internal/pkg/kafka"
	"leal.co/inventory/internal/ports"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env, %v", err)
	}
}

func main() {
	defer database.GetDB().Close()
  kafka := kafka.NewClient("order_created")
  brokerListener := ports.NewBrokerListener(kafka)
  var wg sync.WaitGroup
  wg.Add(1)
  brokerListener.StartListening()
  wg.Wait()
}


