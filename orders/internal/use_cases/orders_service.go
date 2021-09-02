package usecases

import (
	"google.golang.org/protobuf/proto"
	"leal.co/orders/internal/domain"
	"leal.co/orders/internal/dto"
	"leal.co/orders/internal/pkg/kafka"
	"leal.co/orders/internal/pkg/pb"
)


type Service interface {
  CreateOrder(dto.CreateOrderDTO) error
  UpdateOrderStatus(int, string)
}

func NewService(domain domain.Domain, kafkaClient kafka.Client) Service {
  return &service{
   domain,
   kafkaClient,
  }
}

type service struct {
 orders domain.Domain
 kafkaClient kafka.Client
}

func (s *service) CreateOrder(input dto.CreateOrderDTO) error {
  id, err := s.orders.NewOrder(input)
  if err != nil {
    return err
  }

  productsPayload := make([]*pb.Product, len(input.Products))

  for i, product := range input.Products {
    productsPayload[i] = &pb.Product{
      ProductId: int64(product.Code),
      Quantity: int64(product.Quantity),
    }
  }

  orderCreatedEventPayload := &pb.OrderCreated{
    OrderId: int64(id),
    Products: productsPayload,
  }
  out, err := proto.Marshal(orderCreatedEventPayload)
  if err != nil {
    panic(err)
  }

  s.kafkaClient.SendRecord("order_created", out)

  return nil
}

func (s *service) UpdateOrderStatus(orderId int, status string) {
  s.orders.UpdateOrderStatus(orderId, status)
}
