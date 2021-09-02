package domain

import "leal.co/orders/internal/dto"

type Domain interface {
  NewOrder(dto.CreateOrderDTO) (int, error)
  UpdateOrderStatus(int, string)
}
