package repository

type Repository interface {
	SaveOrder(NewOrderDTO) (int, error)
	GetCity(int) CityDTO
	GetClient(int) ClientDTO
  UpdateOrderStatus(int, string)
}
