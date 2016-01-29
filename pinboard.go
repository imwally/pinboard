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

// Post holds values needed to construct a valid URL that is used to
// make a GET request to the Pinboard API.
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

// Response holds the response of the Pinboard API GET requests.
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

var (
	// A map of currently supported methods whose values are
	// Pinboard API endpoints.
	methods = map[string]string{
		"add":    "/posts/add",
		"delete": "/posts/delete",
		"show":   "/posts/recent",
	}
)

// Get shortens an http.Get and returns the body.
func get(u string) (body []byte, err error) {
	res, err := http.Get(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return body, nil
}

// UnmarshalResponse unmarshal's the json response from the Pinboard
// API into the Response struct.
func unmarshalResponse(body []byte) (r Response, err error) {
	var res Response

	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Println(err)
		return res, err
	}

	return res, nil
}

// PinboardMethod expects a valid method and a *Post. A valid URL is
// constructed from the *Post and finally it makes a call to the
// Pinboard API based on the method argument.
func pinboardMethod(method string, p *Post) (r Response, err error) {
	var res Response

	method, ok := methods[method]
	if !ok {
		return res, errors.New(method + " is not a valid pinboard method")
	}

	p.Encode()
	p.Encoded.Path = ver + method

	json, err := get(p.Encoded.String())
	if err != nil {
		return res, err
	}

	res, err = unmarshalResponse(json)
	if err != nil {
		return res, err
	}

	if res.Result != "done" {
		return res, errors.New(res.Result)
	}

	return res, nil
}

// Encode takes the field values from a Post and constructs the URL
// needed to make the GET request to the pinboard API. It saves the
// encoded URL in the Post itself as Post.Encoded.
func (p *Post) Encode() {
	u, err := url.Parse(api)
	if err != nil {
		log.Println(err)
	}

	q := u.Query()
	q.Set("format", "json")
	q.Set("auth_token", p.Token)
	q.Set("url", p.URL)
	q.Set("count", strconv.Itoa(p.Count))
	q.Set("title", p.Title)
	q.Set("tag", p.Tag)
	q.Set("tags", p.Tags)
	q.Set("description", p.Description)
	q.Set("extended", p.Extended)
	q.Set("dt", p.Dt)
	q.Set("replace", p.Replace)
	q.Set("shared", p.Shared)
	q.Set("toread", p.Toread)
	
	u.RawQuery = q.Encode()
	p.Encoded = u
}

// Add calls PinboardMethod to add a new bookmark.
func (p *Post) Add() error {
	_, err := pinboardMethod("add", p)
	return err
}

// Delete calls PinboardMethod to delete a bookmark.
func (p *Post) Delete() error {
	_, err := pinboardMethod("delete", p)
	return err
}

// ShowRecent will show the most recent bookmarks.
func (p *Post) ShowRecent() Response {
	res, _ := pinboardMethod("show", p)
	return res
}
