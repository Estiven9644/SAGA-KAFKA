package dto

type CreateOrderDTO struct {
  ClientId int `json:"clientId"`
  CityCode int `json:"cityCode"`
  DeliveryAddress string `json:"deliveryAddress"`
  Products []ProductDTO `json:"products"`
}

type ProductDTO struct {
  Code int `json:"code"`
  Name string `json:"name"`
  Quantity int `json:"quantity"`
}
