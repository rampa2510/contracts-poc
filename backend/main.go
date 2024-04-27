package main

import (
	"log/slog"
	"os"
)

type Address struct {
	name    string `json:"name"`
	street  string `json:"street"`
	city    string `json:"city"`
	state   string `json:"state"`
	Pincode int    `json:"pincode"`
}

func main() {
	slog.Info("hello, world", "user", os.Getenv("USER"))
}
