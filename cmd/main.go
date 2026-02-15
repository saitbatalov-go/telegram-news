package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"telegram_news/internal/config"
	"telegram_news/internal/fetcher"
	"telegram_news/internal/notifier"
	"telegram_news/internal/storage"
	"telegram_news/internal/summary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
)


func main() {

	BotAPI,err:= tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot %v",err)
	}

	db,err:= sqlx.Connect("postgres",config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to database %v",err)
	}
	defer db.Close()

	var (
		sourcesStorage = storage.NewSourceStorage(db)
		articlesStorage = storage.NewArticleStorage(db)
		fetcher = fetcher.New(
			articlesStorage,
			sourcesStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		notifier = notifier.New(
			articlesStorage,
			summary.NewOpenAISummarizer(config.Get().OpenAIKey,config.Get().OpenAIPrompt),
			BotAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	ctx, cancel:=signal.NotifyContext(context.Background(), os.Interrupt,syscall.SIGTERM)
	defer cancel()

	go func (ctx context.Context)  {
		if err:=fetcher.Start(ctx); err!=nil {

if !errors.Is(err,context.Canceled) {
	log.Printf("failed to start fetcher %v",err)
			return
}
log.Printf("fetcher stopperd")
		}


	}(ctx)
	}


	// go func (ctx context.Context)  {
		if err:=notifier.Start(ctx); err!=nil {

if !errors.Is(err,context.Canceled) {
	log.Printf("failed to start notifier %v",err)
			return
}
log.Printf("notifier stopperd")
		}


	// }(ctx)
}