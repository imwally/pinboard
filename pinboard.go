package pinboard

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// The Post struct holds all of the values needed to construct a valid URL that
// is used to make a GET request to the pinboard API.
type Post struct {
	URL         string
	Href        string
	Title       string
	Description string
	Extended    string
	Tags        string
	Tag         string
	Dt          string
	Replace     string
	Shared      string
	Toread      string
	Token       string
	Count       int
	Encoded     *url.URL
}

// Response struct holds the response of the pinboard API GET requests.
type Response struct {
	Result string `json:"result_code"`
	Date   string `json:"date"`
	User   string `json:"user"`
	Posts  []Post `json:"posts"`
}

const (
	api string = "https://api.pinboard.in/"
	ver string = "v1"
)

// Get shortens an http.Get and returns the body.
func Get(u string) []byte {
	res, err := http.Get(u)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	return body
}

// UnmarshalResponse will unmarshal the json response from the API into the
// Response struct.
func UnmarshalResponse(body []byte) Response {
	var r Response
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
	}

	return r
}

// Encode is where the magic happens. It takes the field values from a Post
// and constructs the URL needed to make the GET request to the pinboard API. It
// saves the encoded URL in the Post itself as Post.Encoded.
func (p *Post) Encode() {
	u, err := url.Parse(api)
	if err != nil {
		log.Println(err)
	}

	q := u.Query()
	q.Set("auth_token", p.Token)
	q.Set("format", "json")

	if p.URL != "" {
		q.Set("url", p.URL)
	}

	if p.Count > 0 {
		q.Set("count", strconv.Itoa(p.Count))
	}

	if p.Title != "" {
		q.Set("title", p.Title)
	}

	if p.Tag != "" {
		q.Set("tag", p.Tag)
	}

	if p.Tags != "" {
		q.Set("tags", p.Tags)
	}

	if p.Description != "" {
		q.Set("description", p.Description)
	}

	if p.Extended != "" {
		q.Set("extended", p.Extended)
	}

	if p.Dt != "" {
		q.Set("dt", p.Dt)
	}

	if p.Replace != "" {
		q.Set("replace", p.Replace)
	}

	if p.Shared != "" {
		q.Set("shared", p.Shared)
	}

	if p.Toread != "" {
		q.Set("toread", p.Toread)
	}

	u.RawQuery = q.Encode()
	p.Encoded = u
}

// Add adds a new bookmark. It sets the constructed GET request URL's path to
// /posts/add.
func (p *Post) Add() error {
	p.Encoded.Path = ver + "/posts/add"
	json := Get(p.Encoded.String())
	res := UnmarshalResponse(json)

	if res.Result != "done" {
		return errors.New(res.Result)
	}

	return nil
}

// Delete deletes a bookrmark. It sets the constructed GET request URL's path to
// /posts/delete.
func (p *Post) Delete() error {
	p.Encoded.Path = ver + "/posts/delete"
	json := Get(p.Encoded.String())
	res := UnmarshalResponse(json)

	if res.Result != "done" {
		return errors.New(res.Result)
	}

	return nil
}

// ShowRecent will show the most recent bookmarks. It sets the constructed GET
// request URL's path to /posts/recent.
func (p *Post) ShowRecent() Response {
	p.Encoded.Path = ver + "/posts/recent"
	json := Get(p.Encoded.String())
	res := UnmarshalResponse(json)

	return res
}
