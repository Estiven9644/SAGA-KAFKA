package ports

import (
	"encoding/json"
	"net/http"

	"leal.co/orders/internal/dto"
	usecases "leal.co/orders/internal/use_cases"
)


type ordersServer struct {
  app usecases.Service
}

func NewOrdersServer(app usecases.Service) ordersServer {
  return ordersServer{
    app: app,
  }
}

func (o *ordersServer) CreateOrder(w http.ResponseWriter, r *http.Request) {
  dto := &dto.CreateOrderDTO{}

  json.NewDecoder(r.Body).Decode(dto)

  responseMessage := "order created correctly"
  responseCode := http.StatusCreated

  if err := o.app.CreateOrder(*dto); err != nil {
    responseMessage = err.Error()
    responseCode = http.StatusBadRequest
  }

  response := struct {
    Message string `json:"message"`
  }{
    Message: responseMessage,
  }

  responseBytes, err := json.Marshal(response)
  if err != nil {
    http.Error(w, http.StatusText(500), http.StatusInternalServerError)
    return
  }

  w.WriteHeader(responseCode)
  w.Write(responseBytes)
}

