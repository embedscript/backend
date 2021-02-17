package handler

import (
	"context"
	"strings"

	pb "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/logger"
	filesproto "github.com/micro/services/files/proto"
)

type V1 struct{}

func (e *V1) Serve(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	files := filesproto.NewFilesService("files", client.DefaultClient)
	logger.Infof("Serving %v", req.Path)

	resp, err := files.List(ctx, &filesproto.ListRequest{
		Project: req.Path,
	})
	if err != nil {
		return err
	}

	// ? huh
	rsp.Header["Content-Type"] = &pb.Pair{
		Key:    "Content-Type",
		Values: []string{"text/html", "charset=UTF-8"},
	}
	rsp.Body = resp.Files[0].FileContents
	return nil
}

func (e *V1) ServeSeparate(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	files := filesproto.NewFilesService("files", client.DefaultClient)
	logger.Infof("Serving %v", req.Path)

	resp, err := files.List(ctx, &filesproto.ListRequest{
		Project: req.Path,
	})
	if err != nil {
		return err
	}
	htmlFile := ""
	jsFile := ""
	cssFile := ""
	for _, file := range resp.Files {
		switch {
		case strings.Contains(file.Path, "main"):
			jsFile = file.FileContents
		case strings.Contains(file.Path, "index"):
			htmlFile = file.FileContents
		case strings.Contains(file.Path, "style"):
			cssFile = file.FileContents
		}
	}

	rendered := `<html><head><script src="https://embedscript.com/assets/micro.js"></script></head><body><div><style>` +
		cssFile +
		`</style>` +
		htmlFile +
		`</div><script>` +
		jsFile +
		`</script></body></html>`

	// ? huh
	rsp.Header["Content-Type"] = &pb.Pair{
		Key:    "Content-Type",
		Values: []string{"text/html", "charset=UTF-8"},
	}
	rsp.Body = rendered
	return nil
}
