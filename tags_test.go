package pinboard

import (
	"testing"
)

func contains(haystack []string, needle string) bool {
	for _, a := range haystack {
		if a == needle {
			return true
		}

	}

	return false
}

func TestTagsGet(t *testing.T) {
	tags, err := TagsGet()
	if err != nil {
		t.Error(err)
	}

	// Crappy test but I'm not sure of another way right now
	if len(tags) < 100 {
		t.Error("expected more tags")
	}
}

func TestTagsRename(t *testing.T) {
	err := PostsAdd(optAdd)
	if err != nil {
		t.Error(err)
	}

	err = TagsRename("pinboard_1_testing", "pinboard_2_testing")
	if err != nil {
		t.Error(err)
	}

	p, err := PostsGet(&PostsGetOptions{
		URL: optAdd.URL,
	})
	if err != nil {
		t.Error(err)
	}

	if !contains(p[0].Tags, "pinboard_2_testing") {
		t.Error("error: tag rename failed")
	}

	err = TagsRename("pinboard_2_testing", "pinboard_1_testing")
	if err != nil {
		t.Error(err)
	}

	p, err = PostsGet(&PostsGetOptions{
		URL: optAdd.URL,
	})
	if err != nil {
		t.Error(err)
	}

	if !contains(p[0].Tags, "pinboard_1_testing") {
		t.Error("error: tag rename failed")
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}
}

func TestTagsDelete(t *testing.T) {
	err := PostsAdd(optAdd)
	if err != nil {
		t.Error(err)
	}

	err = TagsDelete("pinboard_1_testing")
	if err != nil {
		t.Error(err)
	}

	p, err := PostsGet(&PostsGetOptions{
		URL: optAdd.URL,
	})
	if err != nil {
		t.Error(err)
	}

	if contains(p[0].Tags, "pinboard_1_testing") {
		t.Error("error: TagsDelete failed")
	}

	// Clean up by removing test post
	err = PostsDelete(optAdd.URL)
	if err != nil {
		t.Errorf("error: clean up failed: %s", err)
	}
}
