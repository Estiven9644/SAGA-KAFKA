package repository

type NewOrderDTO struct {
  ClientId int
  CityCode int
  DeliveryAddress string
  Products []int
  DeliveryCost int
}

type CityDTO struct {
  Code int
  DeliveryCost int
}

type ClientDTO struct {
  Id int
}
