package main

import (
	"fmt"
	"tgbot/configs"
	"tgbot/db"
	"tgbot/handlers"
	"tgbot/server"
	"tgbot/storage"

	"go.uber.org/fx"
)

func main() {
	fmt.Println("man func")

	app := fx.New(
		fx.Provide(
			configs.Config,
			db.Init,
			server.InitTelegram,
			storage.New,
			handlers.New,
		),

		fx.Invoke(
			server.Start,
		),
	)

	app.Run()
}
