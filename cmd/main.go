package main

import (
	"github.com/mastrasec/firewow/internal/adapter/client"
	"github.com/mastrasec/firewow/internal/entrypoint"
	"github.com/mastrasec/firewow/internal/service"
)

func main() {
	svc := service.New(client.New())

	entry := entrypoint.New("localhost:8765", svc)

	if err := entry.Run(); err != nil {
		return
	}
}
