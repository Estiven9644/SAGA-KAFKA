package ports

import (
	"log"

	"google.golang.org/protobuf/proto"
	"leal.co/inventory/internal/pkg/database"
	"leal.co/inventory/internal/pkg/kafka"
	"leal.co/inventory/internal/pkg/pb"
)

type BrokerListener interface {
  StartListening()
}

type brokerListener struct {
  client kafka.Client
}

func NewBrokerListener(cl kafka.Client) BrokerListener {
  return &brokerListener{
    client: cl,
  }
}

func (b *brokerListener) StartListening() {
  buffer := make(chan []byte)
  go b.client.ConsumeRecords(buffer)

  for i := 0; i<100; i++ {
    go b.work(buffer)
  }
}

func (b *brokerListener) work(buffer chan []byte) {
  for record := range buffer {
    orderCreated := &pb.OrderCreated{}
    if err := proto.Unmarshal(record, orderCreated); err != nil {
      log.Printf("error unmarshaling proto orderCreated: %v", err)
    }

    validInventory := true
    for _, product := range orderCreated.Products {
      if !ValidateInventory(int(product.ProductId), int(product.Quantity)) {
        validInventory = false
      }
    }

    if validInventory {
      for _, product := range orderCreated.Products {
        UpdateInventory(int(product.ProductId), int(orderCreated.OrderId), int(product.Quantity))
      }
    }

    var resultStatus string
    if validInventory {
      resultStatus = "success"
    } else {
      resultStatus = "failure"
    }

    orderConfirmed := &pb.OrderResult{
      OrderId: orderCreated.OrderId,
      Status: resultStatus,
    }
    out, err := proto.Marshal(orderConfirmed)
    if err != nil {
      panic(err)
    }

    b.client.SendRecord("order_result", out)
  }
}

func ValidateInventory(productId, quantity int) bool {
  db := database.GetDB()
  var qty int
  db.QueryRow("SELECT quantity from inventory WHERE product_id = $1", productId).
  Scan(&qty)

  return qty - quantity >= 0
}

func UpdateInventory(productId, orderId, quantity int){
  db := database.GetDB()
  res, err := db.Exec("INSERT INTO transactions (order_id, product_id) VALUES($1, $2)", orderId, productId)
  if err != nil {
    panic(err)
  }

  log.Printf("result from updating transactions: %+v", res)

  res, err = db.Exec("UPDATE inventory SET quantity = (quantity - $1) WHERE product_id = $2", quantity, productId)
  if err != nil {
    panic(err)
  }

  log.Printf("result from updating transactions: %+v", res)
}


