package main

import (
	"github.com/embedscript/backend/emails/handler"
	pb "github.com/embedscript/backend/emails/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("emails"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterEmailsHandler(srv.Server(), handler.NewEmailsHandler())

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
