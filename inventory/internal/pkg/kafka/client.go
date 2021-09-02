package kafka

import (
	"context"
	"log"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Client interface {
  SendRecord(string, []byte)
  ConsumeRecords(chan []byte)
}

type kafkaClient struct {
  cl *kgo.Client
}

func NewClient(topic string) Client {
  seeds := os.Getenv("SEEDS")
  cl, err := kgo.NewClient(
    kgo.SeedBrokers(seeds),
    kgo.ConsumerGroup("my-group-identifier"),
    kgo.ConsumeTopics(topic),
    kgo.DefaultProduceTopic(topic),
    kgo.AllowAutoTopicCreation(),
    kgo.WithLogger(kgo.BasicLogger(os.Stdout, kgo.LogLevelInfo, func() string { return "Kafka" })),
  )

  if err != nil {
    log.Fatalf("could not start kafka client: %v", err)
  }

  return &kafkaClient{
    cl,
  }
}

func (k *kafkaClient) SendRecord (topic string, message []byte) {
	ctx := context.Background()

  record := &kgo.Record{Topic: topic, Value: message}
	k.cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
		if err != nil {
			log.Printf("record had a produce error: %v\n", err)
		}
	})

}

func (k *kafkaClient) ConsumeRecords (buffer chan []byte) {
	ctx := context.Background()
	for {
		fetches := k.cl.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			panic(errs)
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
      buffer<- record.Value
	}
		k.cl.CommitOffsetsSync(ctx, k.cl.UncommittedOffsets(), nil)
	}

}
