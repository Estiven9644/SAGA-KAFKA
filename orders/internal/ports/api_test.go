package ports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"leal.co/orders/internal/dto"
	"leal.co/orders/mocks"
)

func CreateRequest(t *testing.T) (*http.Request, *dto.CreateOrderDTO) {
  newOrder := &dto.CreateOrderDTO{
    ClientId: 123,
    CityCode: 123456,
    DeliveryAddress: "address",
    Products: []dto.ProductDTO{
      {123, "prod 1"},
      {456, "prod 2"},
    },
  }

  payload, err := json.Marshal(newOrder)
  if err != nil {
    t.Fail()
  }
  return httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(payload)), newOrder
}

func TestCreateOrder(t *testing.T) {

  t.Run("returns succesfull message if creation was correct", func(t *testing.T) {
    req, newOrder := CreateRequest(t)
    mockService := new(mocks.Service)
    mockService.On("CreateOrder", *newOrder).Return(nil)

    ordersServer := NewOrdersServer(mockService)

    recorder := httptest.NewRecorder()

    ordersServer.CreateOrder(recorder, req)

    mockService.AssertExpectations(t)
    assert.Equal(t, http.StatusCreated, recorder.Code)
    assert.Equal(t, "{\"message\":\"order created correctly\"}", recorder.Body.String())
  })

  t.Run("returns error message if creation was incorrect", func(t *testing.T) {
    req, newOrder := CreateRequest(t)

    mockService := new(mocks.Service)
    mockService.On("CreateOrder", *newOrder).Return(fmt.Errorf("error message"))

    ordersServer := NewOrdersServer(mockService)

    recorder := httptest.NewRecorder()

    ordersServer.CreateOrder(recorder, req)

    mockService.AssertExpectations(t)
    assert.Equal(t, http.StatusBadRequest, recorder.Code)
    assert.Equal(t, "{\"message\":\"error message\"}", recorder.Body.String())
  })

}
