package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"leal.co/orders/internal/domain/repository"
	"leal.co/orders/internal/dto"
	"leal.co/orders/mocks"
)

func TestNewOrder(t *testing.T) {
	input := dto.CreateOrderDTO{
		ClientId:        123,
		CityCode:        123456,
		DeliveryAddress: "address",
		Products: []dto.ProductDTO{
			{123, "prod 1"},
			{456, "prod 2"},
		},
	}

	newOrderRepoDTO := repository.NewOrderDTO{
		ClientId:        123,
		CityCode:        123456,
		DeliveryAddress: "address",
		Products:        []int{123, 456},
		DeliveryCost:    789,
	}

	t.Run("creates the order correctly", func(t *testing.T) {
		mockRepository := new(mocks.Repository)
		mockRepository.On("GetCity", 123456).Return(repository.CityDTO{123456, 789})
		mockRepository.On("GetClient", 123).Return(repository.ClientDTO{123})
		mockRepository.On("SaveOrder", newOrderRepoDTO).Return(nil)
		orders := NewDomain(mockRepository)
		got := orders.NewOrder(input)

		mockRepository.AssertExpectations(t)
		assert.NoError(t, got)
	})

	t.Run("should return error if the client doesn't exist", func(t *testing.T) {
		mockRepository := new(mocks.Repository)
		mockRepository.On("GetCity", 123456).Return(repository.CityDTO{123456, 789})
		mockRepository.On("GetClient", 123).Return(repository.ClientDTO{})

		orders := NewDomain(mockRepository)
		got := orders.NewOrder(input)

		mockRepository.AssertNotCalled(t, "SaveOrder")
		mockRepository.AssertExpectations(t)
		assert.Error(t, got)
	})

	t.Run("should return error if the city doesn't exist", func(t *testing.T) {
		mockRepository := new(mocks.Repository)
		mockRepository.On("GetCity", 123456).Return(repository.CityDTO{})

		orders := NewDomain(mockRepository)
		got := orders.NewOrder(input)

		mockRepository.AssertNotCalled(t, "SaveOrder")
		mockRepository.AssertNotCalled(t, "GetClient")
		mockRepository.AssertExpectations(t)
		assert.Error(t, got)
	})
}
