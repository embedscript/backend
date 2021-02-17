package main

import (
	"github.com/embedscript/backend/api/handler"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/api"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("v1"),
		service.Version("latest"),
	)

	srv.Server().Handle(
		srv.Server().NewHandler(
			new(handler.V1),
			api.WithEndpoint(
				&api.Endpoint{
					Name:    "V1.Serve",
					Path:    []string{"^/v1/.*$"},
					Method:  []string{"GET", "POST", "OPTIONS", "PUT", "HEAD", "DELETE"},
					Handler: "api",
				}),
		))
	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
