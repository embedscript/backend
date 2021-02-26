package handler

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/micro/micro/v3/service/auth"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/model"

	datastore "github.com/embedscript/backend/datastore/proto"
)

var indexIndex = model.Index{
	FieldName: "TypeOf",
}

type IndexRecord struct {
	ID     string
	TypeOf string
	Index  model.Index
}

type Owner struct {
	ID string
}

type Rule struct {
	datastore.Rule
	ID string
}

type Datastore struct {
}

func (e *Datastore) Create(ctx context.Context, req *datastore.CreateRequest, rsp *datastore.CreateResponse) error {
	log.Info("Received Datastore.Create request")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(req.Value), &m)
	if err != nil {
		return err
	}
	err = e.authAction(ctx, req.Project, req.Table, "write")
	if err != nil {
		return err
	}

	indexes, err := e.getIndexes(ctx, req.Project, req.Table)
	if err != nil {
		return err
	}
	db := model.New(map[string]interface{}{}, &model.Options{
		Indexes:   indexes,
		Namespace: req.Project + req.Table,
	})
	return db.Create(m)
}

func (e *Datastore) authAction(ctx context.Context, project, table, action string) error {
	owners, err := e.getOwner(ctx, project, table)
	rules, err := e.getRules(ctx, project, table)
	if err != nil {
		return err
	}
	unregisteredEnabled := false
	userEnabled := false
	for _, rule := range rules {
		switch {
		case rule.Action == action && rule.Role == "unregistered":
			unregisteredEnabled = true
		case rule.Action == action && rule.Role == "user":
			userEnabled = true
		}
	}
	acc, ok := auth.AccountFromContext(ctx)
	if !unregisteredEnabled && !ok {
		return errors.New("unauthorized: need to be registered")
	}
	if !userEnabled {
		owner, err := e.getOwner(ctx, project, table)
		if err != nil {
			return err
		}
		if len(owners) == 0 {
			return errors.New("unauthorized: no owners")
		}
		if owner[0].ID != acc.ID {
			return errors.New("unauthorized: not an owner")
		}
	}
	return nil
}

func (e *Datastore) Update(ctx context.Context, req *datastore.UpdateRequest, rsp *datastore.UpdateResponse) error {
	log.Info("Received Datastore.Update request")
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(req.Value), &m)
	if err != nil {
		return err
	}
	indexes, err := e.getIndexes(ctx, req.Project, req.Table)
	if err != nil {
		return err
	}
	db := model.New(map[string]interface{}{}, &model.Options{
		Indexes: indexes,
	})
	return db.Update(m)
}

func (e *Datastore) getIndexes(ctx context.Context, project, table string) ([]model.Index, error) {
	indexDb := model.New(map[string]interface{}{}, &model.Options{
		Namespace: project + table + "indexes",
	})
	result := []IndexRecord{}
	err := indexDb.Read(model.QueryAll(), &result)
	if err != nil {
		return nil, err
	}
	indexes := []model.Index{}
	for _, v := range result {
		indexes = append(indexes, v.Index)
	}
	return indexes, nil
}

func (e *Datastore) getRules(ctx context.Context, project, table string) ([]Rule, error) {
	indexDb := model.New(map[string]interface{}{}, &model.Options{
		Namespace: project + table + "rules",
	})
	result := []Rule{}
	err := indexDb.Read(model.QueryAll(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Datastore) getOwner(ctx context.Context, project, table string) ([]Owner, error) {
	indexDb := model.New(map[string]interface{}{}, &model.Options{
		Namespace: project + table + "owners",
	})
	result := []Owner{}
	err := indexDb.Read(model.QueryAll(), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (e *Datastore) saveOwner(ctx context.Context, project, table, ownerID string) error {
	ownerDb := model.New(map[string]interface{}{}, &model.Options{
		Namespace: project + table + "owners",
	})
	ownerDb.Create(Owner{
		ID: ownerID,
	})
	return nil
}

func (e *Datastore) saveRule(ctx context.Context, rule *datastore.Rule) error {
	ownerDb := model.New(map[string]interface{}{}, &model.Options{
		Namespace: rule.Project + rule.Table + "rules",
	})
	return ownerDb.Create(Rule{
		ID: rule.Role + rule.Action,
	})
}

func (e *Datastore) Read(ctx context.Context, req *datastore.ReadRequest, rsp *datastore.ReadResponse) error {
	log.Info("Received Datastore.Read request")
	q := toQuery(req.Query)
	result := []map[string]interface{}{}
	indexes, err := e.getIndexes(ctx, req.Project, req.Table)
	if err != nil {
		return err
	}
	db := model.New(map[string]interface{}{}, &model.Options{
		Indexes: indexes,
	})
	err = db.Read(q, &result)
	if err != nil {
		return err
	}
	js, err := json.Marshal(result)
	rsp.Value = string(js)
	return err
}

func (e *Datastore) CreateIndex(ctx context.Context, req *datastore.CreateIndexRequest, rsp *datastore.CreateIndexResponse) error {
	log.Info("Received Datastore.Index request")

	index := toIndex(req.Index)
	indexRecord := IndexRecord{
		ID: index.FieldName + index.Type + index.Order.FieldName + string(index.Order.Type),
	}
	db := model.New(IndexRecord{}, &model.Options{
		Namespace: req.Project + req.Table + "indexes",
	})
	return db.Create(indexRecord)
}

func (e *Datastore) CreateRule(ctx context.Context, req *datastore.CreateRuleRequest, rsp *datastore.CreateRuleResponse) error {
	log.Info("Received Datastore.CreateRule request")
	owners, err := e.getOwner(ctx, req.Rule.Project, req.Rule.Table)
	if err != nil {
		return err
	}
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.New("unauthorized")
	}
	if len(owners) > 0 && owners[0].ID != acc.ID {
		return errors.New("unauthorized - not owner")
	}
	err = e.saveOwner(ctx, req.Rule.Project, req.Rule.Table, acc.ID)
	if err != nil {
		return err
	}
	return e.saveRule(ctx, req.Rule)
}

func (e *Datastore) Delete(ctx context.Context, req *datastore.DeleteRequest, rsp *datastore.DeleteResponse) error {
	log.Info("Received Datastore.Delete request")
	q := toQuery(req.Query)
	return model.New(map[string]interface{}{}, nil).Delete(q)
}

func toQuery(pquery *datastore.Query) model.Query {
	q := model.QueryEquals(pquery.Index.FieldName, pquery.Value)
	if pquery.Order != nil {
		q.Order.FieldName = pquery.Order.FieldName
		q.Order.Type = model.OrderType(pquery.Order.OrderType.String())
	}
	return q
}

func toIndex(pindex *datastore.Index) model.Index {
	i := model.Index{
		FieldName: pindex.FieldName,
		Type:      pindex.Type,
		Unique:    pindex.Unique,
	}
	if pindex.Order != nil {
		i.Order = model.Order{
			FieldName: pindex.FieldName,
			Type:      model.OrderType(strings.ToLower(pindex.Order.OrderType.String())),
		}
	}
	return i
}
