package domain

import (
	"fmt"

	"leal.co/orders/internal/domain/repository"
	"leal.co/orders/internal/dto"
)

type Order struct {
	client          Client
	city            City
	shippingCost    int
	deliveryAddress string
	products        []Product
}

type Client struct {
	Id int
}

type City struct {
	Code          int
	ShipphingCost int
}

type Product struct {
	code int
	name string
}

type orders struct {
	repo repository.Repository
}

func NewDomain(repo repository.Repository) Domain {
	return &orders{
		repo: repo,
	}
}

func (o *orders) NewOrder(input dto.CreateOrderDTO) (int, error) {

	city := o.repo.GetCity(input.CityCode)
	if city.Code == 0 {
		return 0, fmt.Errorf("city does not exist")
	}

	client := o.repo.GetClient(input.ClientId)
	if client.Id == 0 {
		return 0, fmt.Errorf("client does not exist")
	}

	productsIds := make([]repository.ProductDTO, len(input.Products))
	for i, product := range input.Products {
		productsIds[i] = repository.ProductDTO{
			Id:       product.Code,
			Quantity: product.Quantity,
		}
	}

	newOrderRepoDTO := repository.NewOrderDTO{
		ClientId:        client.Id,
		CityCode:        input.CityCode,
		DeliveryAddress: input.DeliveryAddress,
		Products:        productsIds,
		DeliveryCost:    city.DeliveryCost,
		Status:          "pending",
	}

	id, err := o.repo.SaveOrder(newOrderRepoDTO)

	return id, err
}

func (o *orders) UpdateOrderStatus(orderId int, status string) {
	var statusToSave string
	if status == "success" {
		statusToSave = "readyForShipping"
	} else if status == "failure" {
		statusToSave = "cancelled"
	} else {
		// handle error
		return
	}

	o.repo.UpdateOrderStatus(orderId, statusToSave)
}
