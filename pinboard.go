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

type Post struct {
	Url         string
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

func UnmarshalResponse(body []byte) Response {
	var r Response
	err := json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
	}

	return r
}

func (p *Post) Encode() {
	u, err := url.Parse(api)
	if err != nil {
		log.Println(err)
	}

	q := u.Query()
	q.Set("auth_token", p.Token)
	q.Set("format", "json")

	if p.Url != "" {
		q.Set("url", p.Url)
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

func (p *Post) Add() error {
	p.Encoded.Path = ver + "/posts/add"
	json := Get(p.Encoded.String())
	res := UnmarshalResponse(json)

	if res.Result != "done" {
		return errors.New(res.Result)
	}

	return nil
}

func (p *Post) Delete() error {
	p.Encoded.Path = ver + "/posts/delete"
	json := Get(p.Encoded.String())
	res := UnmarshalResponse(json)

	if res.Result != "done" {
		return errors.New(res.Result)
	}

	return nil
}

func (p *Post) ShowRecent() Response {
	p.Encoded.Path = ver + "/posts/recent"
	json := Get(p.Encoded.String())
	res := UnmarshalResponse(json)

	return res
}
