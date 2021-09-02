package adapters

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"leal.co/orders/internal/domain/repository"
)

func TestSaveOrder(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Errorf("an error %v was not expected when opening a stub database connection", err)
  }

  defer db.Close()

  newOrderRepoDTO := repository.NewOrderDTO{
    ClientId: 123,
    CityCode: 123456,
    DeliveryAddress: "address",
    Products: []int{123},
    DeliveryCost: 789,
  }


  mock.ExpectQuery("INSERT INTO orders \\(client_id, city_code, delivery_address, delivery_cost\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) returning id").
  WithArgs(newOrderRepoDTO.ClientId, newOrderRepoDTO.CityCode, newOrderRepoDTO.DeliveryAddress, newOrderRepoDTO.DeliveryCost).
  WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

  mock.ExpectExec("INSERT INTO orders_products \\(order_id, product_id\\) VALUES \\(\\$1, \\$2\\)").
  WithArgs(1, 123).
  WillReturnResult(sqlmock.NewResult(1, 1))

  sqlRepo := NewSqlRepository(db)

  err = sqlRepo.SaveOrder(newOrderRepoDTO)
  assert.NoError(t, err)
  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expectations: %s", err)
  }
}

func TestGetCity(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Errorf("an error %v was not expected when opening a stub database connection", err)
  }

  defer db.Close()

  mock.ExpectQuery("SELECT \\* FROM cities WHERE code = \\$1").
  WithArgs(123).
  WillReturnRows(sqlmock.NewRows([]string{"code", "delivery_price"}).AddRow(123, 789))

  sqlRepo := NewSqlRepository(db)
  got := sqlRepo.GetCity(123)

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expectations: %s", err)
  }

  want := repository.CityDTO{123, 789}
  assert.Equal(t, got, want)
}

func TestGetClient(t *testing.T) {
  db, mock, err := sqlmock.New()
  if err != nil {
    t.Errorf("an error %v was not expected when opening a stub database connection", err)
  }

  defer db.Close()

  mock.ExpectQuery("SELECT \\* FROM clients WHERE id = \\$1").
  WithArgs(123).
  WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

  sqlRepo := NewSqlRepository(db)
  got := sqlRepo.GetClient(123)

  if err := mock.ExpectationsWereMet(); err != nil {
    t.Errorf("there were unfulfilled expectations: %s", err)
  }
  want := repository.ClientDTO{123}
  assert.Equal(t, got, want)
}
