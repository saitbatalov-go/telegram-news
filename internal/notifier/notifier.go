package notifier

import (
	"context"
	"io"
	"net/http"
	"strings"

	"telegram_news/internal/model"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkPosted(ctx context.Context, id int64) error
}

type Summarizer interface {
	Summarize(ctx context.Context, article model.Article) (string, error)
}

type Notifier struct {
	articles         ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelID        int64
}

func New(articles ArticleProvider, summarizer Summarizer, bot *tgbotapi.BotAPI, sendInterval time.Duration, lookupTimeWindow time.Duration, channelID int64) *Notifier {
	return &Notifier{
		articles:         articles,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelID:        channelID,
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneAcrtiles, err := n.articles.AllNotPosted(ctx, time.Now().Add(-n.lookupTimeWindow), 1)
	if err != nil {
		return err
	}
	if len(topOneAcrtiles) == 0 {
		return nil
	}

	article := topOneAcrtiles[0]

	summary, err := n.extractSummary(ctx, article)
	if err != nil {
		return err
	}
	if err := n.sendAricle(article, summary); err != nil {
		return err
	}
	if err := n.articles.MarkPosted(ctx, article.ID); err != nil {
		return err
	}
	return nil
}

func (n *Notifier) extractSummary(ctx context.Context, article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		r = resp.Body
	}


	doc, err: = readability.FromReader(r)
	if err != nil {
		return "", err
	}
	
	summary,err:=n.simmarizer.Summarize(ctx, clearText(oc.TextContent))
	if err != nil {
		return "", err
	}
	return "\n\n"+summary, nil

}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

func clearText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}
