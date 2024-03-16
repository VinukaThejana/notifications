// Produces events to be sent to Kafka
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/VinukaThejana/notifications/internal/lib"
	"github.com/VinukaThejana/notifications/internal/models"
	"github.com/fatih/color"
	"github.com/twmb/franz-go/pkg/kgo"
)

func send(ctx context.Context, cl *kgo.Client, payload struct {
	From models.User
	To   models.User
},
) {
	notification := models.Notification{
		From: payload.From,
		To:   payload.To,
	}

	notificationJSON, err := json.Marshal(notification)
	if err != nil {
		logger.ErrorfWithMsg(err, "Failed to marshal the notification JSON")
	}

	record := &kgo.Record{
		Topic: "notifications",
		Value: notificationJSON,
	}
	cl.Produce(ctx, record, func(_ *kgo.Record, err error) {
		if err != nil {
			logger.ErrorfWithMsg(err, "Failed to produce the Kafka event")
		}
	})
}

func main() {
	users := []models.User{
		{ID: 1, Name: "JohnDoe"},
		{ID: 2, Name: "Bruno"},
		{ID: 3, Name: "Jack"},
		{ID: 4, Name: "Sparrow"},
	}

	cl, err := kgo.NewClient()
	if err != nil {
		logger.ErrorfWithMsg(err, "Failed to initialize Kafka")
	}
	defer cl.Close()
	ctx := context.Background()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	userInput := make(chan struct {
		From models.User
		To   models.User
	})

	var fromStr, toStr string

	go func() {
		for {
			fmt.Print("\nFrom user ID : ")
			fmt.Scanln(&fromStr)

			fmt.Print("To user ID   : ")
			fmt.Scanln(&toStr)

			to, err := strconv.Atoi(toStr)
			if err != nil {
				fmt.Println(color.RedString("user ID must be an integer"))
				continue
			}
			from, err := strconv.Atoi(fromStr)
			if err != nil {
				fmt.Println(color.RedString("user ID must be an integer"))
				continue
			}

			toUser := lib.FindUserByID(to, users)
			fromUser := lib.FindUserByID(from, users)
			if toUser == nil {
				fmt.Println(color.RedString("To user with the given ID does not exsist"))
				continue
			}
			if fromUser == nil {
				fmt.Println(color.RedString("From user with the given ID does not exsist"))
				continue
			}

			userInput <- struct {
				From models.User
				To   models.User
			}{
				From: *fromUser,
				To:   *toUser,
			}

			fmt.Print(color.GreenString("\nFrom : %s\n", fromUser.Name))
			fmt.Print(color.GreenString("To   : %s\n", toUser.Name))
		}
	}()

	for {
		select {
		case sig := <-shutdown:
			logger.Log(color.RedString(fmt.Sprintf("[%s] : Program aborted !", sig.String())))
			return
		case input := <-userInput:
			send(ctx, cl, struct {
				From models.User
				To   models.User
			}{
				From: input.From,
				To:   input.To,
			})
		}
	}
}
