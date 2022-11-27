package server

import (
	"fmt"
	"net/http"
	"tgbot/configs"
	"tgbot/handlers"
)

// Start ...
func Start(config *configs.Configuration, h *handlers.HandlerService) {

	err := http.ListenAndServe(config.ServerPort, http.HandlerFunc(h.GlobalHandler))
	if err != nil {
		fmt.Println("error on staring server!")
	}
	fmt.Println("server running on port : ", config.ServerPort)
}
