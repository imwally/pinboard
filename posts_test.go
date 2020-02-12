package pinboard

import (
	"fmt"
	"log"
	"testing"
	"time"
)

// Test TestPostsAdd first as almost all subsequent tests rely on
// adding a post.
//
// go test -v -failfast
func TestPostsAdd(t *testing.T) {
	err := PostsAdd(optAdd)
	if err != nil {
		t.Error(err)
	}

	err = PostsAdd(&PostsAddOptions{
		Description: "This should fail",
	})
	if err == nil {
		t.Error("error: expected missing url error")
	}

	err = PostsAdd(&PostsAddOptions{
		URL: "https://github.com/imwally/pinboard",
	})
	if err == nil {
		t.Error("error: expected must provide title error")
	}
}

func ExamplePostsAdd() {
	opt := &PostsAddOptions{
		URL:         "https://github.com/imwally/pinboard",
		Description: "Testing Pinboard Go Package",
	}

	err := PostsAdd(opt)
	if err != nil {
		log.Println("error adding post:", err)
	}
}

// Next, make sure PostsDelete works so we can remove test posts
// created in subsequent tests. This will delete the post that was
// created with TestPostsAdd.
func TestPostsDelete(t *testing.T) {
	err := PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: failed to delete test post: %s", err)
	}

	url := "https://thisisjustatest.com"
	err = PostsDelete(url)
	if err == nil {
		t.Error("error: expected item not found error")
	}
}

func TestPostsUpdate(t *testing.T) {
	timeBeforeAdd := time.Now()

	err := PostsAdd(optAdd)
	if err != nil {
		t.Errorf("error: failed to create test post: %s", err)
	}

	timeAfterAdd := time.Now()

	timeUpdate, err := PostsUpdate()
	if err != nil {
		t.Error(err)
	}

	if timeUpdate.After(timeAfterAdd) && timeUpdate.Before(timeBeforeAdd) {
		t.Errorf("error: expected %s to be between %s and %s",
			timeUpdate,
			timeBeforeAdd,
			timeAfterAdd,
		)
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}
}

func TestPostsGet(t *testing.T) {
	err := PostsAdd(optAdd)
	if err != nil {
		t.Errorf("error: failed to create test post: %s", err)
	}

	// Test PostsGet by URL
	posts, err := PostsGet(&PostsGetOptions{
		URL: optAdd.URL,
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) > 1 {
		t.Error("error: PostsGet: expected only 1 post")
	}

	// Test PostsGet by tag
	dt, _ := time.Parse("2006-01-02", "2010-12-11")

	posts, err = PostsGet(&PostsGetOptions{
		Dt:  dt,
		Tag: []string{"pinboard", "testing"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 1 {
		t.Error("error: PostsGet: expected only 1 post")
	}

	// Test PostsGet by Dt
	posts, err = PostsGet(&PostsGetOptions{
		Dt: dt,
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 1 {
		t.Error("error: PostsGet: expected only 1 post")
	}

	// Did we get the meta signature?
	if posts[0].Meta == nil {
		t.Error("error: PostsGet: expected meta signature")
	}

	posts, err = PostsGet(&PostsGetOptions{
		Tag: []string{"sadklfjsldkfjsdlkfj"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 0 {
		t.Error("error: PostsGet: expected zero posts")
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}
}

func ExamplePostsGet() {
	dt, err := time.Parse("2006-01-02", "2010-12-11")
	if err != nil {
		log.Println(err)
	}

	posts, err := PostsGet(&PostsGetOptions{Dt: dt})
	if err != nil {
		log.Println("error getting posts:", err)
	}

	for _, post := range posts {
		fmt.Println(post.Description)
		fmt.Println(post.Href)
		fmt.Println(post.Time)
	}

	// Output:
	// Testing Pinboard Go Package
	// https://github.com/imwally/pinboard
	// 2010-12-11 19:48:02 +0000 UTC
}

func TestPostsRecent(t *testing.T) {
	// Test Count
	posts, err := PostsRecent(&PostsRecentOptions{
		Count: 100,
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 100 {
		t.Error("error: expected 100 posts")
	}

	err = PostsAdd(optAdd)
	if err != nil {
		t.Errorf("error: failed to create test post: %s", err)
	}

	posts, err = PostsRecent(&PostsRecentOptions{
		Tag: []string{"pinboard_testing"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 1 {
		t.Error("error: expected 1 post")
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}

}

func TestPostsDates(t *testing.T) {
	err := PostsAdd(optAdd)
	if err != nil {
		t.Errorf("error: failed to create test post: %s", err)
	}

	dates, err := PostsDates(nil)
	if err != nil {
		t.Error(err)
	}

	expected := "2010-12-11"
	if _, ok := dates[expected]; !ok {
		t.Errorf("error: expected at least 1 post on %s", expected)
	}

	// Test with tags
	dates, err = PostsDates(&PostsDatesOptions{
		Tag: []string{"pinboard_testing"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(dates) != 1 {
		t.Error("error: expected only 1 post")
	}

	if _, ok := dates[expected]; !ok {
		t.Errorf("error: expected at least 1 post on %s", expected)
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}

	// Make sure that date no longer appears
	dates, err = PostsDates(nil)
	if err != nil {
		t.Error(err)
	}

	if _, ok := dates[expected]; ok {
		t.Errorf("error: expected no posts on %s", expected)
	}
}

func TestPostsAll(t *testing.T) {
	err := PostsAdd(optAdd)
	if err != nil {
		t.Errorf("error: failed to create test post: %s", err)
	}

	posts, err := PostsAll(&PostsAllOptions{
		Results: 1,
		Fromdt:  optAdd.Dt.Add(time.Duration(-1) * time.Second),
		Todt:    optAdd.Dt.Add(time.Second),
		Tag:     []string{"pinboard", "testing"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 1 {
		t.Errorf("error: PostsAll: expected only 1 post")
	}

	posts, err = PostsAll(&PostsAllOptions{
		Results: 1,
		Fromdt:  optAdd.Dt.Add(time.Duration(-1) * time.Second),
		Todt:    optAdd.Dt.Add(time.Second),
		Tag:     []string{"this should fail"},
	})
	if err != nil {
		t.Error(err)
	}

	if len(posts) != 0 {
		t.Errorf("error: PostsAll: expected 0 posts")
	}

	posts, err = PostsAll(nil)
	if err != nil {
		t.Error(err)
	}

	if len(posts) < 200 {
		t.Errorf("error: PostsAll: expected more than 200 posts")
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}

}

// The following two tests don't require posting a bookmark to get
// suggested tags.
func TestPostsSuggestPopular(t *testing.T) {
	got, err := PostsSuggestPopular(optAdd.URL)
	if err != nil {
		t.Error(err)
	}

	if len(got) > 1 || got[0] != "code" {
		t.Error("error: expected single tag code")
	}
}

func TestPostsSuggestRecommended(t *testing.T) {
	got, err := PostsSuggestRecommended(optAdd.URL)
	if err != nil {
		t.Error(err)
	}

	if len(got) != 3 {
		t.Error("error: expected 3 tags")
	}

	if got[0] != "code" {
		t.Error("error: expected code tag")
	}

	if got[1] != "github" {
		t.Error("error: expected github tag")
	}

	if got[2] != "IFTTT" {
		t.Error("error: expected IFTTT tag")
	}
}
