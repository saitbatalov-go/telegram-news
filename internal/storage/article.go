package storage

import (
	"context"
	"database/sql"
	"telegram-news/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (a *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {

}

func (a *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) (model.Article, error) {
}

func (a *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
}

type dbArticle struct {
	ID          int64        `json:"id"`
	SourceID    int64        `json:"source_id"`
	Title       string       `json:"title"`
	Link        string       `json:"link"`
	Summary     string       `json:"summary"`
	PublishedAt time.Time    `json:"published_at"`
	PostedAt    sql.NullTime `json:"posted_at"`
}
