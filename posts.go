package pinboard

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"
)

// Post represents a bookmark.
type Post struct {
	// URL of bookmark.
	Href *url.URL

	// Title of bookmark. This field is unfortunately named
	// 'description' for backwards compatibility with the
	// delicious API
	Description string

	// Description of the item. Called 'extended' for backwards
	// compatibility with delicious API.
	Extended []byte

	// Tags of bookmark.
	Tags []string

	// If the bookmark is private or public.
	Shared bool

	// If the bookmark is marked to read later.
	Toread bool

	// Create time for this bookmark.
	Time time.Time

	// Change detection signature of the bookmark.
	Meta []byte

	// Hash of the bookmark.
	Hash []byte

	// The number of other users who have bookmarked this same
	// item.
	Others int
}

// post represents intermediate post response data before type
// conversion.
type post struct {
	Href        string
	Description string
	Extended    string
	Tags        string
	Shared      string
	Toread      string
	Time        string
	Meta        string
	Hash        string
	Others      int
}

// toPost converts a post to a type correct Post.
func (p *post) toPost() (*Post, error) {
	href, err := url.Parse(p.Href)
	if err != nil {
		return nil, err
	}

	tags := strings.Split(p.Tags, " ")

	var shared, toread bool
	if p.Shared == "yes" {
		shared = true
	}

	if p.Toread == "yes" {
		toread = true
	}

	dt, err := time.Parse(time.RFC3339, p.Time)
	if err != nil {
		return nil, err
	}

	P := Post{
		Href:        href,
		Description: p.Description,
		Extended:    []byte(p.Extended),
		Tags:        tags,
		Shared:      shared,
		Toread:      toread,
		Time:        dt,
		Meta:        []byte(p.Meta),
		Hash:        []byte(p.Hash),
		Others:      p.Others,
	}

	return &P, nil
}

// postsResponse represents a response from certain /posts/ endpoints.
type postsResponse struct {
	UpdateTime string `json:"update_time,omitempty"`
	ResultCode string `json:"result_code,omitempty"`
	Date       string `json:"date,omitempty"`
	User       string `json:"user,omitempty"`
	Posts      []post `json:"posts,omitempty"`
}

// PostsUpdate returns the most recent time a bookmark was added,
// updated or deleted.
//
// https://pinboard.in/api/#posts_update
func PostsUpdate() (time.Time, error) {
	resp, err := get("postsUpdate", nil)
	if err != nil {
		return time.Time{}, err
	}

	var pr postsResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return time.Time{}, err
	}

	update, err := time.Parse(time.RFC3339, pr.UpdateTime)
	if err != nil {
		return time.Time{}, err
	}

	return update, nil
}

// PostsAddOptions represents the required and optional arguments for
// adding a bookmark.
type PostsAddOptions struct {
	// Required: The URL of the item.
	URL string

	// Required: Title of the item. This field is unfortunately
	// named 'description' for backwards compatibility with the
	// delicious API.
	Description string

	// Description of the item. Called 'extended' for backwards
	// compatibility with delicious API.
	Extended []byte

	// List of up to 100 tags.
	Tags []string

	// Creation time for this bookmark. Defaults to current
	// time. Datestamps more than 10 minutes ahead of server time
	// will be reset to current server time.
	Dt time.Time

	// Replace any existing bookmark with this URL. Default is
	// yes. If set to no, will throw an error if bookmark exists.
	Replace bool

	// Make bookmark public. Default is "yes" unless user has
	// enabled the "save all bookmarks as private" user setting,
	// in which case default is "no".
	Shared bool

	// Marks the bookmark as unread. Default is "no".
	Toread bool
}

// PostsAdd adds a bookmark.
//
// https://pinboard.in/api/#posts_add
func PostsAdd(opt *PostsAddOptions) error {
	if opt.URL == "" {
		return errors.New("error: missing url")
	}

	if opt.Description == "" {
		return errors.New("error: missing description")
	}

	resp, err := get("postsAdd", opt)
	if err != nil {
		return err
	}

	var pr postsResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return err
	}

	if pr.ResultCode != "done" {
		return errors.New(pr.ResultCode)
	}

	return nil
}

// postsDeleteOptions represents the single required argument for
// deleting a bookmark.
type postsDeleteOptions struct {
	URL string
}

// PostsDelete deletes the bookmark by url.
//
// https://pinboard.in/api/#posts_delete
func PostsDelete(url string) error {
	resp, err := get("postsDelete", &postsDeleteOptions{URL: url})
	if err != nil {
		return err
	}

	var pr postsResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return err
	}

	if pr.ResultCode != "done" {
		return errors.New(pr.ResultCode)
	}

	return nil
}

// PostsGetOptions represents the optional arguments for getting
// bookmarks.
type PostsGetOptions struct {
	// Filter by up to three tags.
	Tag []string

	// Return results bookmarked on this day. UTC date in this
	// format: 2010-12-11.
	Dt time.Time

	// Return bookmark for this URL.
	URL string

	// Include a change detection signature in a meta attribute.
	Meta bool
}

// PostsGet returns one or more posts (on a single day) matching the
// arguments. If no date or URL is given, date of most recent bookmark
// will be used.Returns one or more posts on a single day matching the
// arguments. If no date or URL is given, date of most recent bookmark
// will be used.
//
// https://pinboard.in/api/#posts_get
func PostsGet(opt *PostsGetOptions) ([]*Post, error) {
	resp, err := get("postsGet", opt)
	if err != nil {
		return nil, err
	}

	var pr postsResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, p := range pr.Posts {
		post, err := p.toPost()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// PostsRecentOptions represents the optional arguments for returning
// the user's most recent posts.
type PostsRecentOptions struct {
	// Filter by up to three tags.
	Tag []string

	// Number of results to return. Default is 15, max is 100.
	Count int
}

// PostsRecent returns a list of the user's most recent posts,
// filtered by tag.
//
// https://pinboard.in/api/#posts_recent
func PostsRecent(opt *PostsRecentOptions) ([]*Post, error) {
	resp, err := get("postsRecent", opt)
	if err != nil {
		return nil, err
	}

	var pr postsResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, p := range pr.Posts {
		post, err := p.toPost()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// postsDatesResponse represents the response from /posts/dates.
type postsDatesResponse struct {
	User  string            `json:"user"`
	Tag   string            `json:"tag"`
	Dates map[string]string `json:"dates"`
}

// PostsDatesOptions represents the single optional argument for
// returning a list of dates with the number of posts at each date.
type PostsDatesOptions struct {
	// Filter by up to three tags.
	Tag []string
}

// PostsDates returns a list of dates with the number of posts at each
// date.
//
// https://pinboard.in/api/#posts_dates
func PostsDates(opt *PostsDatesOptions) (map[string]string, error) {
	resp, err := get("postsDates", opt)
	if err != nil {
		return nil, err
	}

	var pr postsDatesResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return nil, err
	}

	return pr.Dates, nil
}

// PostsAllOptions represents the optional arguments for returning all
// bookmarks in the user's account.
type PostsAllOptions struct {
	// Filter by up to three tags.
	Tag []string

	// Offset value (default is 0).
	Start int

	// Number of results to return. Default is all.
	Results int

	// Return only bookmarks created after this time.
	Fromdt time.Time

	// Return only bookmarks created before this time.
	Todt time.Time

	// Include a change detection signature for each bookmark.
	//
	// Note: This probably doesn't work. A meta field is always
	// returned. The Pinboard API says the datatype is an int but
	// changing the value has no impact on the results. Using a
	// yes/no string like all the other meta options doesn't work
	// either.
	Meta int
}

// PostsAll returns all bookmarks in the user's account.
//
// https://pinboard.in/api/#posts_all
func PostsAll(opt *PostsAllOptions) ([]*Post, error) {
	resp, err := get("postsAll", opt)
	if err != nil {
		return nil, err
	}

	var pr []post
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return nil, err
	}

	var posts []*Post
	for _, p := range pr {
		post, err := p.toPost()
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

// postSuggestResponse represents the response from /posts/suggest.
type postsSuggestResponse struct {
	Popular     []string `json:"popular"`
	Recommended []string `json:"recommended"`
}

// postSuggestOptions represents the single required argument, url,
// for suggesting tags for a post.
type postsSuggestOptions struct {
	URL string
}

// PostsSuggestPopular returns a slice of popular tags for a given
// URL. Popular tags are tags used site-wide for the url.
//
// https://pinboard.in/api/#posts_suggest
func PostsSuggestPopular(url string) ([]string, error) {
	resp, err := get("postsSuggest", &postsSuggestOptions{URL: url})
	if err != nil {
		return nil, err
	}

	var pr []postsSuggestResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return nil, err
	}

	return pr[0].Popular, nil
}

// PostsSuggestRecommended returns a slice of recommended tags for a
// given URL. Recommended tags are drawn from the user's own tags.
//
// https://pinboard.in/api/#posts_suggest
func PostsSuggestRecommended(url string) ([]string, error) {
	resp, err := get("postsSuggest", &postsSuggestOptions{URL: url})
	if err != nil {
		return nil, err
	}

	var pr []postsSuggestResponse
	err = json.Unmarshal(resp, &pr)
	if err != nil {
		return nil, err
	}

	return pr[1].Recommended, nil
}
