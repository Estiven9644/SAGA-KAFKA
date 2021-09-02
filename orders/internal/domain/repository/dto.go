package repository

type NewOrderDTO struct {
  ClientId int
  CityCode int
  DeliveryAddress string
  Products []ProductDTO
  DeliveryCost int
  Status string
}

type ProductDTO struct {
  Id int
  Quantity int
}

type CityDTO struct {
  Code int
  DeliveryCost int
}

type ClientDTO struct {
  Id int
}
