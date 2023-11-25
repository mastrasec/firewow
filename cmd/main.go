package main

import (
	"net/http"

	"github.com/mastrasec/firewow/internal/adapter/client"
	"github.com/mastrasec/firewow/internal/entrypoint"
	"github.com/mastrasec/firewow/internal/service"
	"github.com/mastrasec/firewow/internal/service/canary"
)

func main() {

	canaryService := canary.New()
	svc := service.New(
		client.New(&http.Client{}),
		canaryService,
	)

	entry := entrypoint.New("localhost:8765", svc)

	if err := entry.Run(); err != nil {
		return
	}
}
