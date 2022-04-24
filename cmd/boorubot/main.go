package main

import (
	"net/http"
	"os"
	"time"

	"github.com/eientei/boorubot/integration/booru/danbooru"
	"github.com/eientei/boorubot/integration/pleroma"
	"github.com/eientei/boorubot/internal/boorubot"
)

func main() {
	danbooruClient, err := danbooru.NewClient(&danbooru.Config{
		HTTPClient: nil,
		APIKey:     "",
		Login:      "",
		URL:        "https://booru.eientei.org",
	})
	if err != nil {
		panic(err)
	}

	pleromaClient, err := pleroma.NewClient(&pleroma.Config{
		HTTPClient: nil,
		APIKey:     os.Getenv("PLEROMA_API_KEY"),
		URL:        "https://eientei.org",
	})
	if err != nil {
		panic(err)
	}

	bot, err := boorubot.NewBot(&boorubot.Config{
		PleromaClient:  pleromaClient,
		DanbooruClient: danbooruClient,
		StateProvider: &boorubot.FileStateProvider{
			FileName: "state.json",
		},
		HTTPClient:   &http.Client{},
		Interval:     time.Hour,
		PostInterval: time.Second * 10,
	})
	if err != nil {
		panic(err)
	}

	bot.Start()
}
