package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"leal.co/orders/internal/dto"
	"leal.co/orders/mocks"
)

func TestCreateOrder(t *testing.T) {
  input := dto.CreateOrderDTO{
    ClientId: 123,
    CityCode: 123456,
    DeliveryAddress: "address",
    Products: []dto.ProductDTO{
      {123, "prod 1"},
      {456, "prod 2"},
    },
  }

  mockDomain := new(mocks.Domain)
  mockDomain.On("NewOrder", input).Return(nil)
  mockDomain.NewOrder(input)
  service := NewService(mockDomain)

  mockDomain.AssertExpectations(t)
  err := service.CreateOrder(input)
  assert.NoError(t, err)
}
