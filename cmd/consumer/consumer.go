// Receives events from Kafka
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/VinukaThejana/notifications/internal/models"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	cl, err := kgo.NewClient(
		kgo.ConsumeTopics("notifications"),
	)
	if err != nil {
		logger.ErrorfWithMsg(err, "Failed to initialize the Kafka client")
	}
	defer cl.Close()

	ctx := context.Background()

	recordChan := make(chan models.Notification)
	var record models.Notification

	go func() {
		for {
			fetches := cl.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				logger.ErrorfWithMsg(err, "Return Poll errors")
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				err := json.Unmarshal(iter.Next().Value, &record)
				if err != nil {
					logger.ErrorWithMsg(err, "Failed to unmarshal the record")
				} else {
					recordChan <- record
				}
			}
		}
	}()

	for {
		select {
		case record := <-recordChan:
			fmt.Printf("\nNotification from %s to %s\n", record.From.Name, record.To.Name)
		default:
			continue
		}
	}
}
