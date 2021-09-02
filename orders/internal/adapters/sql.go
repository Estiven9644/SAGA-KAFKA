package adapters

import (
	"database/sql"
	"fmt"
	"log"

	"leal.co/orders/internal/domain/repository"
)

type sqlRepository struct {
  db *sql.DB
}

func NewSqlRepository(db *sql.DB) repository.Repository {
  return &sqlRepository{
    db,
  }
}

func (s *sqlRepository) SaveOrder(orders repository.NewOrderDTO) (int, error) {
  var id int
	if err := s.db.QueryRow("INSERT INTO orders (client_id, city_code, delivery_address, delivery_cost, status) VALUES ($1, $2, $3, $4, $5) returning id",
		orders.ClientId, orders.CityCode, orders.DeliveryAddress, orders.DeliveryCost, orders.Status).
		Scan(&id); err != nil {
		return 0, fmt.Errorf("error creating user, %v", err)
	}

  for _, product := range orders.Products {
    _, err := s.db.Exec("INSERT INTO orders_products (order_id, product_id, quantity) VALUES ($1, $2, $3)",
    id, product.Id, product.Quantity)
    if err != nil {
      return 0, fmt.Errorf("error saving product: %d", product)
    }
  }

  return id, nil
}

func (s *sqlRepository) GetCity(code int) repository.CityDTO {
  var ans repository.CityDTO
  err := s.db.QueryRow("SELECT * FROM cities WHERE code = $1", code).
  Scan(&ans.Code, &ans.DeliveryCost)
  if err != nil {
    fmt.Printf("could not get cities: %v", err)
    return ans
  }

  return ans
}

func (s *sqlRepository) GetClient(id int) repository.ClientDTO {
  var ans repository.ClientDTO

  err := s.db.QueryRow("SELECT * FROM clients WHERE id = $1", id).
  Scan(&ans.Id)
  if err != nil {
    fmt.Printf("could not get client: %v", err)
    return ans
  }

  return ans
}

func (s *sqlRepository) UpdateOrderStatus(orderId int, status string) {
  res, err := s.db.Exec("UPDATE orders SET status = $2 WHERE id = $1", orderId, status)
  if err != nil {
    log.Printf("error updating order: %d status: %v", orderId, err)
  }

  log.Printf("result from updating status: %+v", res)
}


