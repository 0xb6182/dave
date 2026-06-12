package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/0xb6182/dave/config"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt any) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Printf("%s from %v\n", v.Message.GetConversation(), v.Info.PushName)
	}
}

func Run() error {
	dbLog := waLog.Noop
	clientLog := waLog.Noop

	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", config.DB, dbLog)
	if err != nil {
		return err
	}

	// TODO: multi device support
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return err
	}
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	err = client.Connect()
	if err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Disconnecting...")
	client.Disconnect()
	return nil
}
