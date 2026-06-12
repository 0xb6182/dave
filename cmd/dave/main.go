package main

import (
	"context"
	"os"

	"github.com/0xb6182/dave/internal/bot"
	"github.com/urfave/cli/v3"
)

func main() {
	cli.HelpFlag = nil
	cmd := cli.Command{
		Name:  "dave",
		Usage: "A discord-like whatsapp bot",
		Commands: []*cli.Command{
			{
				Name:    "login",
				Aliases: []string{"l"},
				Usage:   "Login the bot to whatsapp with a qr code",
				Action: func(ctx context.Context, c *cli.Command) error {
					return bot.Login()
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Run the bot server",
				Action: func(ctx context.Context, c *cli.Command) error {
					return bot.Run()
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
