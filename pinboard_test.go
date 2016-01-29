package pinboard

import (
	"testing"
)

// A valid token is needed to make calls to the Pinboard API.
var token = ""

func TestAdd(t *testing.T) {

	p := new(Post)

	p.Token = token
	p.Description = "TESTING PINBOARD CLI CLIENT"
	p.URL = "https://github.com/imwally/pinboard/testing"
	p.Tags = "pin pinboard test testing"
	p.Extended = "This is a test from imwally's golang pinboard package. For more information please refer to the pinned URL."
	p.Toread = "yes"
	p.Shared = "yes"

	t.Log("Adding: ", p.URL)
	err := p.Add()
	if err != nil {
		t.Error(err)
	}
}

func TestShow(t *testing.T) {
	p := new(Post)

	p.Token = token
	p.Count = 2

	recent := p.ShowRecent()
	for _, r := range recent.Posts {
		t.Log("Description:", r.Description)
		t.Log("URL:", r.Href)
		t.Log("Tags:", r.Tags)
		t.Log("Extended:", r.Extended)
		t.Log("To Read:", r.Toread)
		t.Log("Shared:", r.Shared)
		t.Log("----")
	}
}

func TestDelete(t *testing.T) {
	p := new(Post)

	p.Token = token
	p.URL = "https://github.com/imwally/pinboard/testing"

	t.Log("Deleting: ", p.URL)
	err := p.Delete()
	if err != nil {
		t.Error(err)
	}
}
