package pinboard

import (
	"encoding/json"
	"errors"
)

// Tags maps a tag name to the number of bookmarks that use that tag.
type Tags map[string]string

// tagsResponse holds the response result from deleting or renaming
// tags.
type tagsResponse struct {
	Result string `json:"result"`
}

// TagsGet returns a full list of the user's tags along with the
// number of times they were used.
func TagsGet() (Tags, error) {
	resp, err := get("tagsGet", nil)
	if err != nil {
		return nil, err
	}

	var tags Tags
	err = json.Unmarshal(resp, &tags)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

// tagsDeleteOptions holds the single required argument to delete a
// tag.
type tagsDeleteOptions struct {
	Tag string
}

// TagsDelete deletes an existing tag.
func TagsDelete(tag string) error {
	resp, err := get("tagsDelete", &tagsDeleteOptions{Tag: tag})
	if err != nil {
		return err
	}

	var tr tagsResponse
	err = json.Unmarshal(resp, &tr)
	if err != nil {
		return err
	}

	if tr.Result != "done" {
		return errors.New(tr.Result)
	}

	return nil
}

// tagsRenameOptions holds the required arguments needed to rename a
// tag.
type tagsRenameOptions struct {
	Old string
	New string
}

// TagsRename renames a tag, or folds it in to an existing tag.
func TagsRename(old, new string) error {
	resp, err := get("tagsRename", &tagsRenameOptions{
		Old: old,
		New: new,
	})
	if err != nil {
		return err
	}

	var tr tagsResponse
	err = json.Unmarshal(resp, &tr)
	if err != nil {
		return err
	}

	if tr.Result != "done" {
		return errors.New(tr.Result)
	}

	return nil
}
