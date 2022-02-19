package config

import "fmt"

const HOST = "localhost"
const PORT = "8080"
const PROTOCOL = "http"

// Возврат корня сайта
func GetMainLink() string {
	return fmt.Sprintf(
		"%s://%s:%s",
		PROTOCOL,
		HOST,
		PORT,
	)
}

// адрес для запуска сервера, без протокола
func GetServerAddress() string {
	return fmt.Sprintf(
		"%s:%s",
		HOST,
		PORT,
	)
}
