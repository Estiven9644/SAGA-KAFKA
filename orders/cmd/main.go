package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"leal.co/orders/internal/adapters"
	"leal.co/orders/internal/domain"
	apihelpers "leal.co/orders/internal/pkg/api"
	"leal.co/orders/internal/pkg/database"
	"leal.co/orders/internal/pkg/kafka"
	"leal.co/orders/internal/ports"
	usecases "leal.co/orders/internal/use_cases"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

  var wg sync.WaitGroup
  wg.Add(1)
  go runServer()
  go runBrokerListener()
  wg.Wait()
}

func runServer() {
  r := chi.NewRouter()

  r.Use(middleware.RequestID)
  r.Use(middleware.RealIP)
  r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)
  r.Use(apihelpers.JSONContentTypeMiddleware)

  kafkaClient := kafka.NewClient("order_result")

  repo := adapters.NewSqlRepository(database.GetDB())
  domain := domain.NewDomain(repo)
  app := usecases.NewService(domain, kafkaClient)
  api := ports.NewOrdersServer(app)


  r.Post("/orders", api.CreateOrder)

  log.Println("starting server")
  http.ListenAndServe(":8080", r)
}

func runBrokerListener() {
  kafkaClient := kafka.NewClient("order_result")
  repo := adapters.NewSqlRepository(database.GetDB())
  domain := domain.NewDomain(repo)
  app := usecases.NewService(domain, kafkaClient)
  brokerListener := ports.NewBrokerListener(kafkaClient, app)

  brokerListener.StartListening()

}
