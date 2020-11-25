package handler

import (
	"context"
	"time"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/micro/dev/model"

	proto "github.com/embedscript/backend/posts/proto"
	tags "github.com/embedscript/backend/tags/proto"
	"github.com/gosimple/slug"
)

const (
	tagType = "post-tag"
)

func getPostModel(website string) model.Model {
	createdIndex := model.ByEquality("created")
	createdIndex.Order.Type = model.OrderTypeDesc
	idIndex := model.ByEquality("Id")
	idIndex.Order.Type = model.OrderTypeUnordered

	return model.New(
		store.DefaultStore,
		proto.Post{},
		model.Indexes(model.ByEquality("slug"), createdIndex),
		&model.ModelOptions{
			Debug:     false,
			Namespace: website,
			IdIndex:   idIndex,
		},
	)
}

type Website struct {
	// website url eg. example.com
	Id      string
	OwnerID string
}

type Posts struct {
	Tags           tags.TagsService
	websites       model.Model
	websiteIDIndex model.Index
}

func NewPosts(tagsService tags.TagsService) *Posts {
	createdIndex := model.ByEquality("created")
	createdIndex.Order.Type = model.OrderTypeDesc

	websiteIDIndex := model.ByEquality("Id")
	websiteIDIndex.Order.Type = model.OrderTypeUnordered

	return &Posts{
		Tags: tagsService,
		websites: model.New(store.DefaultStore, Website{}, nil, &model.ModelOptions{
			IdIndex: websiteIDIndex,
		}),
		websiteIDIndex: websiteIDIndex,
	}
}

func (p *Posts) Save(ctx context.Context, req *proto.SaveRequest, rsp *proto.SaveResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("proto.save.input-check", "Not logged in")
	}
	if len(req.Website) == 0 {
		return errors.Unauthorized("proto.save.input-check", "Website missing")
	}
	websites := []Website{}
	err := p.websites.List(p.websiteIDIndex.ToQuery(req.Website), &websites)
	if err != nil {
		return err
	}
	if len(websites) == 0 {
		// allow save, tie website to user account
		err = p.websites.Save(Website{
			Id:      req.Website,
			OwnerID: acc.ID,
		})
		if err != nil {
			return err
		}
	} else {
		if websites[0].OwnerID != acc.ID {
			return errors.Unauthorized("proto.save.input-check", "Not authorized")
		}
	}

	// read by post
	posts := []*proto.Post{}
	q := model.Equals("Id", req.Id)
	q.Order.Type = model.OrderTypeUnordered
	err = getPostModel(req.Website).List(q, &posts)
	if err != nil {
		return errors.InternalServerError("proto.save.store-id-read", "Failed to read post by id: %v", err.Error())
	}
	postSlug := slug.Make(req.Title)
	// If no existing record is found, create a new one
	if len(posts) == 0 {
		post := &proto.Post{
			Id:       req.Id,
			Title:    req.Title,
			Content:  req.Content,
			Tags:     req.Tags,
			Slug:     postSlug,
			Created:  time.Now().Unix(),
			Metadata: req.Metadata,
			Image:    req.Image,
		}
		err := p.savePost(ctx, nil, post)
		if err != nil {
			return errors.InternalServerError("proto.save.post-save", "Failed to save new post: %v", err.Error())
		}
		return nil
	}
	oldPost := posts[0]

	post := &proto.Post{
		Id:       req.Id,
		Title:    oldPost.Title,
		Content:  oldPost.Content,
		Slug:     oldPost.Slug,
		Tags:     oldPost.Tags,
		Created:  oldPost.Created,
		Updated:  time.Now().Unix(),
		Metadata: req.Metadata,
		Image:    req.Image,
	}
	if len(req.Title) > 0 {
		post.Title = req.Title
		post.Slug = slug.Make(post.Title)
	}
	if len(req.Slug) > 0 {
		post.Slug = req.Slug
	}
	if len(req.Content) > 0 {
		post.Content = req.Content
	}
	if len(req.Tags) > 0 {
		// Handle the special case of deletion
		if len(req.Tags) == 0 && req.Tags[0] == "" {
			post.Tags = []string{}
		} else {
			post.Tags = req.Tags
		}
	}

	postsWithThisSlug := []*proto.Post{}
	err = getPostModel(req.Website).List(model.Equals("slug", postSlug), &postsWithThisSlug)
	if err != nil {
		return errors.InternalServerError("proto.save.store-read", "Failed to read post by slug: %v", err.Error())
	}

	if len(postsWithThisSlug) > 0 {
		if oldPost.Id != postsWithThisSlug[0].Id {
			return errors.BadRequest("proto.save.slug-check", "An other post with this slug already exists")
		}
	}

	return p.savePost(ctx, oldPost, post)
}

func (p *Posts) savePost(ctx context.Context, oldPost, post *proto.Post) error {
	err := getPostModel(post.Website).Save(*post)
	return err
	// @todo do not save tags for now
	if err != nil {
		return err
	}
	if oldPost == nil {
		for _, tagName := range post.Tags {
			_, err := p.Tags.Add(ctx, &tags.AddRequest{
				ResourceID: post.Id,
				Type:       tagType,
				Title:      tagName,
			})
			if err != nil {
				return err
			}
		}
		return nil
	}
	return p.diffTags(ctx, post.Id, oldPost.Tags, post.Tags)
}

func (p *Posts) diffTags(ctx context.Context, parentID string, oldTagNames, newTagNames []string) error {
	oldTags := map[string]struct{}{}
	for _, v := range oldTagNames {
		oldTags[v] = struct{}{}
	}
	newTags := map[string]struct{}{}
	for _, v := range newTagNames {
		newTags[v] = struct{}{}
	}
	for i := range oldTags {
		_, stillThere := newTags[i]
		if !stillThere {
			_, err := p.Tags.Remove(ctx, &tags.RemoveRequest{
				ResourceID: parentID,
				Type:       tagType,
				Title:      i,
			})
			if err != nil {
				logger.Errorf("Error decreasing count for tag '%v' with type '%v' for parent '%v'", i, tagType, parentID)
			}
		}
	}
	for i := range newTags {
		_, newlyAdded := oldTags[i]
		if newlyAdded {
			_, err := p.Tags.Add(ctx, &tags.AddRequest{
				ResourceID: parentID,
				Type:       tagType,
				Title:      i,
			})
			if err != nil {
				logger.Errorf("Error increasing count for tag '%v' with type '%v' for parent '%v': %v", i, tagType, parentID, err)
			}
		}
	}
	return nil
}

func (p *Posts) Query(ctx context.Context, req *proto.QueryRequest, rsp *proto.QueryResponse) error {
	if len(req.Website) == 0 {
		return errors.Unauthorized("proto.save.input-check", "Website missing")
	}

	var q model.Query
	if len(req.Slug) > 0 {
		logger.Infof("Reading post by slug: %v", req.Slug)
		q = model.Equals("slug", req.Slug)
	} else if len(req.Id) > 0 {
		logger.Infof("Reading post by id: %v", req.Id)
		q = model.Equals("Id", req.Id)
		q.Order.Type = model.OrderTypeUnordered
	} else {
		q = model.Equals("created", nil)
		q.Order.Type = model.OrderTypeDesc
		var limit uint
		limit = 20
		if req.Limit > 0 {
			limit = uint(req.Limit)
		}
		q.Limit = int64(limit)
		q.Offset = req.Offset
		logger.Infof("Listing posts, offset: %v, limit: %v", req.Offset, limit)
	}

	return getPostModel(req.Website).List(q, &rsp.Posts)
}

func (p *Posts) Delete(ctx context.Context, req *proto.DeleteRequest, rsp *proto.DeleteResponse) error {
	if len(req.Website) == 0 {
		return errors.Unauthorized("proto.save.input-check", "Website missing")
	}
	logger.Info("Received Post.Delete request")
	q := model.Equals("Id", req.Id)
	q.Order.Type = model.OrderTypeUnordered
	return getPostModel(req.Website).Delete(q)
}
