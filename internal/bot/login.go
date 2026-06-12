package bot

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/0xb6182/dave/config"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func Login() error {
	// dbLog := waLog.Stdout("Database", "DEBUG", true)
	// clientLog := waLog.Stdout("Client", "DEBUG", true)
	dbLog := waLog.Noop
	clientLog := waLog.Noop
	// TODO: better logging

	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", config.DB, dbLog)
	if err != nil {
		return err
	}

	deviceStore := container.NewDevice()
	client := whatsmeow.NewClient(deviceStore, clientLog)

	qrChan, _ := client.GetQRChannel(context.Background())
	err = client.Connect()
	if err != nil {
		return err
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateWithConfig(evt.Code, qrterminal.Config{
				Writer:     os.Stdout,
				Level:      qrterminal.L,
				HalfBlocks: true,
			})
		} else {
			fmt.Println("Login event:", evt.Event)
		}
	}

	time.Sleep(10 * time.Second) // TODO: do better (events.AppStateSyncComplete)
	client.Disconnect()
	return nil
}
