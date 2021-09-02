package ports

import (
	"log"

	"google.golang.org/protobuf/proto"
	"leal.co/orders/internal/pkg/kafka"
	"leal.co/orders/internal/pkg/pb"
	usecases "leal.co/orders/internal/use_cases"
)

type BrokerListener interface {
  StartListening()
}

type brokerListener struct {
  client kafka.Client
  app usecases.Service
}

func NewBrokerListener(cl kafka.Client, app usecases.Service) BrokerListener {
  return &brokerListener{
    client: cl,
    app: app,
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
    orderResult := &pb.OrderResult{}
    if err := proto.Unmarshal(record, orderResult); err != nil {
      log.Printf("error unmarshaling proto orderResult: %v", err)
    }
    log.Printf("order result received: %+v", orderResult)
    b.app.UpdateOrderStatus(int(orderResult.OrderId), orderResult.Status)
  }
}
