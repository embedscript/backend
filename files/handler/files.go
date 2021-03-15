package handler

import (
	"context"
	"errors"
	"strings"

	files "github.com/embedscript/backend/files/proto"
	"github.com/micro/micro/v3/service/auth"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"
)

type Files struct {
	db               model.Model
	dbPartialIndexed model.Model
}

func NewFiles() *Files {
	i := model.ByEquality("project")
	i.Order.Type = model.OrderTypeUnordered

	byOwnerID := model.ByEquality("username")
	//byOwnerID.Order.FieldName = "created"
	byOwnerID.Order.Type = model.OrderTypeUnordered

	return &Files{
		db: model.New(
			files.File{},
			&model.Options{
				Key:     "Id",
				Indexes: []model.Index{i, byOwnerID},
				Debug:   false,
			},
		),
		dbPartialIndexed: model.New(
			files.File{},
			&model.Options{
				Key:     "Id",
				Indexes: []model.Index{i},
				Debug:   false,
			},
		),
	}
}

func (e *Files) Save(ctx context.Context, req *files.SaveRequest, rsp *files.SaveResponse) error {
	// @todo return proper micro errors
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.New("Files.Save requires authentication")
	}

	log.Info("Received Files.Save request")
	for _, file := range req.Files {
		f := files.File{}
		err := e.db.Read(model.QueryEquals("Id", file.Id), &f)
		if err != nil && err != model.ErrorNotFound {
			return err
		}
		// if file exists check ownership
		if f.Id != "" && f.Owner != acc.ID {
			return errors.New("Not authorized")
		}

		if acc.Metadata != nil && acc.Metadata["username"] != "" {
			file.Username = acc.Metadata["username"]
		}
		if !strings.Contains(file.Project, "preview") {
			err = e.db.Create(file)
		} else {
			err = e.dbPartialIndexed.Create(file)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Files) List(ctx context.Context, req *files.ListRequest, rsp *files.ListResponse) error {
	log.Info("Received Files.List request")
	rsp.Files = []*files.File{}

	if len(req.Ids) > 0 {
		for _, id := range req.Ids {
			f := files.File{}
			err := e.db.Read(model.QueryEquals("Id", id), &f)
			if err != nil {
				return err
			}
			rsp.Files = append(rsp.Files, &f)
		}
		return nil
	}
	if req.Project != "" {
		err := e.db.Read(model.QueryEquals("project", req.GetProject()), &rsp.Files)
		if err != nil {
			return err
		}
		// @todo funnily while this is the archetypical
		// query for the KV store interface, it's not supported by the model
		// so we do client side filtering here
		if req.Path != "" {
			filtered := []*files.File{}
			for _, file := range rsp.Files {
				if strings.HasPrefix(file.Path, req.Path) {
					filtered = append(filtered, file)
				}
			}
			rsp.Files = filtered
		}
		return nil
	}
	if req.Username != "" {
		q := model.QueryEquals("username", req.GetUsername())
		err := e.db.Read(q, &rsp.Files)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("owner or project required")
}
