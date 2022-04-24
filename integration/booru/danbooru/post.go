package danbooru

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Rating post rating
type Rating string

// Rating values
const (
	RatingSafe         Rating = "s"
	RatingQuestionable Rating = "q"
	RatingExplicit     Rating = "e"
)

// Post model
type Post struct {
	CreatedAt           Time   `json:"created_at"`
	UpdatedAt           Time   `json:"updated_at"`
	LastNotedAt         Time   `json:"last_noted_at"`
	LastCommentBumpedAt Time   `json:"last_comment_bumped_at"`
	LastCommentedAt     Time   `json:"last_commented_at"`
	Source              string `json:"source,omitempty"`
	MD5                 string `json:"md_5,omitempty"`
	TagString           string `json:"tag_string,omitempty"`
	FileExt             string `json:"file_ext,omitempty"`
	TagStringGeneral    string `json:"tag_string_general,omitempty"`
	TagStringCharacter  string `json:"tag_string_character,omitempty"`
	TagStringCopyright  string `json:"tag_string_copyright,omitempty"`
	TagStringArtist     string `json:"tag_string_artist,omitempty"`
	TagStringMeta       string `json:"tag_string_meta,omitempty"`
	FileURL             string `json:"file_url,omitempty"`
	LargeFileURL        string `json:"large_file_url,omitempty"`
	PreviewFileURL      string `json:"preview_file_url,omitempty"`
	Rating              Rating `json:"rating,omitempty"`
	ID                  ID     `json:"id,omitempty"`
	ParentID            ID     `json:"parent_id,omitempty"`
	UploaderID          ID     `json:"uploader_id,omitempty"`
	ApproverID          ID     `json:"approver_id,omitempty"`
	PixivID             ID     `json:"pixiv_id,omitempty"`
	BitFlags            int64  `json:"bit_flags,omitempty"`
	Score               int    `json:"score,omitempty"`
	UpScore             uint   `json:"up_score,omitempty"`
	DownScore           uint   `json:"down_score,omitempty"`
	FileSize            uint   `json:"file_size,omitempty"`
	ImageWidth          uint   `json:"image_width,omitempty"`
	ImageHeight         uint   `json:"image_height,omitempty"`
	FavCount            uint   `json:"fav_count,omitempty"`
	TagCount            uint   `json:"tag_count,omitempty"`
	TagCountGeneral     uint   `json:"tag_count_general,omitempty"`
	TagCountArtist      uint   `json:"tag_count_artist,omitempty"`
	TagCountCharacter   uint   `json:"tag_count_character,omitempty"`
	TagCountCopyright   uint   `json:"tag_count_copyright,omitempty"`
	TagCountMeta        uint   `json:"tag_count_meta,omitempty"`
	IsPending           bool   `json:"is_pending,omitempty"`
	IsFlagged           bool   `json:"is_flagged,omitempty"`
	IsDeleted           bool   `json:"is_deleted,omitempty"`
	HasChildren         bool   `json:"has_children,omitempty"`
	IsBanned            bool   `json:"is_banned,omitempty"`
	HasActiveChildren   bool   `json:"has_active_children,omitempty"`
	HasLarge            bool   `json:"has_large,omitempty"`
	HasVisibleChildren  bool   `json:"has_visible_children,omitempty"`
}

// PostListQuery request parameters
type PostListQuery struct {
	Tags  []string
	Page  int
	Limit int
}

// PostList returns list of posts fetched by provided parameters
func (client *Client) PostList(ctx context.Context, q *PostListQuery) (posts []*Post, err error) {
	var tags string

	if q != nil {
		tags = strings.Join(q.Tags, " ")
	}

	base := client.base

	base.Path += "/posts.json"

	v := url.Values{}

	if len(tags) > 0 {
		v.Set("tags", tags)
	}

	if q.Limit > 0 {
		v.Set("limit", strconv.Itoa(q.Limit))
	}

	if q.Page > 0 {
		v.Set("page", strconv.Itoa(q.Page))
	}

	base.RawQuery = v.Encode()

	err = client.exchange(ctx, http.MethodGet, base.String(), nil, &posts)

	return
}

// PostCountQuery request parameters
type PostCountQuery struct {
	Tags []string
}

// PostCount returns number of posts matching parameters
func (client *Client) PostCount(ctx context.Context, q *PostCountQuery) (count uint64, err error) {
	var tags string

	if q != nil {
		tags = strings.Join(q.Tags, " ")
	}

	base := client.base

	base.Path += "/counts/posts.json"

	v := url.Values{}

	if len(tags) > 0 {
		v.Set("tags", tags)
	}

	base.RawQuery = v.Encode()

	resp := struct {
		Counts struct {
			Posts uint64 `json:"posts"`
		} `json:"counts"`
	}{}

	err = client.exchange(ctx, http.MethodGet, base.String(), nil, &resp)

	return resp.Counts.Posts, err
}
