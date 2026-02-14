package notifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"telegram_news/internal/botkit/markup"
	"telegram_news/internal/model"
	"time"

	"github.com/go-shiori/go-readability"
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

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summarizer.Summarize(ctx, cleanText(doc.TextContent))

	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil

}

func (n *Notifier) sendAricle(article model.Article, summary string) error {
	const msgFormat = " *%s*%s\n\n%s"

	msg := tgbotapi.NewMessage(n.channelID,
		fmt.Sprintf(
			msgFormat,
			markup.EscapeForMarkdown(article.Title),
			markup.EscapeForMarkdown(summary),
			markup.EscapeForMarkdown(article.Link),
		),
	)
	msg.ParseMode = "markdown"
	_, err := n.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

func cleanText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}
