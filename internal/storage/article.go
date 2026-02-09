package storage

import (
	"context"
	"database/sql"
	"telegram_news/internal/model"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (a *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published_at) 
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}
	return nil

}

func (a *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles
		 WHERE posted_at IS NULL
		 AND published_at >= $1::timestamp
		 ORDER BY published DESC
		LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PublishedAt: article.PublishedAt,
			PostedAt:    article.PostedAt.Time,
			CreatedAt:   article.CreatedAt,
		}
	}), nil

}

func (a *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := a.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2`,
		time.Now().UTC().Format(time.RFC3339),
		id,
	); err != nil {
		return err
	}
	return nil
}

type dbArticle struct {
	ID          int64        `json:"id"`
	SourceID    int64        `json:"source_id"`
	Title       string       `json:"title"`
	Link        string       `json:"link"`
	Summary     string       `json:"summary"`
	PublishedAt time.Time    `json:"published_at"`
	PostedAt    sql.NullTime `json:"posted_at"`
	CreatedAt   time.Time    `json:"created_at"`
}
