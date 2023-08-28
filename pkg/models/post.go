package models

import (
	"context"
	"html/template"
	"time"

	"github.com/muhwyndhamhp/marknotes/pkg/post/values"
	"github.com/muhwyndhamhp/marknotes/utils/scopes"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string
	Content        string
	EncodedContent template.HTML
	Status         values.PostStatus
	PublishedAt    time.Time
	FormMeta       map[string]interface{} `gorm:"-"`
	UserID         uint
	User           User
	Tags           []*Tag `gorm:"many2many:post_tags;"`
}

type PostRepository interface {
	Upsert(ctx context.Context, value *Post) error
	GetByID(ctx context.Context, id uint) (*Post, error)
	Get(ctx context.Context, funcs ...scopes.QueryScope) ([]Post, error)
	Delete(ctx context.Context, id uint) error
}

func (m *Post) AppendFormMeta(page int, onlyPublished bool, sortQuery string, keyword string) {
	m.FormMeta = map[string]interface{}{
		"IsLastItem":    true,
		"Page":          page,
		"PublishedOnly": onlyPublished,
		"SortQuery":     sortQuery,
	}

	if keyword != "" {
		m.FormMeta["Keyword"] = keyword
	}
}
