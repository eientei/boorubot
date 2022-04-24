// Package boorubot provides fediverse/pleroma bot implementation for new booru posts
package boorubot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eientei/boorubot/integration/booru/danbooru"
	"github.com/eientei/boorubot/integration/pleroma"
)

// Config for the bot instance
type Config struct {
	PleromaClient  *pleroma.Client
	DanbooruClient *danbooru.Client
	StateProvider  StateProvider
	HTTPClient     *http.Client
	Interval       time.Duration
	PostInterval   time.Duration
}

// NewBot returns new bot instance
func NewBot(config *Config) (*Bot, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	if config.PleromaClient == nil {
		return nil, errors.New("nil pleroma client")
	}

	if config.DanbooruClient == nil {
		return nil, errors.New("nil danbooru client")
	}

	if config.StateProvider == nil {
		return nil, errors.New("nil state provider")
	}

	if config.Interval <= 0 {
		return nil, errors.New("non-positive polling interval")
	}

	if config.PostInterval <= 0 {
		return nil, errors.New("non-positive post interval")
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{}
	}

	return &Bot{
		Config: *config,
	}, nil
}

// Bot instance
type Bot struct {
	Config
}

// Start blockingly polls for new posts and uploads them to pleroma
func (bot *Bot) Start() {
	for {
		err := bot.run()
		if err != nil {
			fmt.Println("Error: ", err)
		}

		time.Sleep(bot.Interval)
	}
}

func (bot *Bot) fetchPosts(ctx context.Context, statelast uint64) (posts []*danbooru.Post, err error) {
	page := 1

	var locposts []*danbooru.Post

outter:
	for {
		locposts, err = bot.DanbooruClient.PostList(ctx, &danbooru.PostListQuery{
			Tags:  nil,
			Page:  page,
			Limit: 100,
		})
		if err != nil {
			return
		}

		if len(locposts) == 0 {
			break
		}

		for _, p := range locposts {
			if p.ID == 0 || p.IsPending || time.Since(time.Time(p.CreatedAt)) < bot.Interval {
				continue
			}

			if uint64(p.ID) <= statelast {
				break outter
			}

			posts = append(posts, p)
		}

		page++
	}

	return
}

func (bot *Bot) run() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), bot.Interval)

	defer cancel()

	state, err := bot.StateProvider.StateLoad(ctx)
	if err != nil {
		return err
	}

	posts, err := bot.fetchPosts(ctx, state.LastPost)
	if err != nil {
		return err
	}

	for i := len(posts) - 1; i >= 0; i-- {
		err = bot.processPost(ctx, posts[i], state)
		if err != nil {
			return err
		}

		time.Sleep(bot.PostInterval)
	}

	return nil
}

func (bot *Bot) processPost(ctx context.Context, post *danbooru.Post, state *State) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, post.FileURL, nil)
	if err != nil {
		return err
	}

	resp, err := bot.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	mediaID, err := bot.PleromaClient.MediaUpload(ctx, post.FileURL, resp.Body)
	if err != nil {
		return err
	}

	postID := post.ID.String()

	var tagstring string

	if len(post.TagStringCharacter) > 0 {
		var tags []string

		for _, t := range strings.Split(post.TagStringCharacter, " ") {
			tags = append(tags, pleroma.MakeTag(t))
		}

		tagstring = strings.Join(tags, " ")

		if len(tagstring) > 0 {
			tagstring = "<br/>" + tagstring
		}
	}

	_, err = bot.PleromaClient.StatusCreate(ctx, &pleroma.StatusCreateRequest{
		Status:      `<a href="//booru.eientei.org/posts/` + postID + `">Post #` + postID + `</a>` + tagstring,
		ContentType: "text/html",
		MediaIDs:    []string{mediaID},
		Sensitive:   post.Rating != danbooru.RatingSafe,
	})
	if err != nil {
		return err
	}

	state.LastPost = uint64(post.ID)

	return bot.StateProvider.StateSave(ctx, state)
}
