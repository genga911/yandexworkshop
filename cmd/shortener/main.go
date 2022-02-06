package main

import (
	"fmt"
	"github.com/genga911/yandexworkshop/internal/app/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env")
	}

	server.SetUpServer()
}
