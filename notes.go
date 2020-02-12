package pinboard

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"
)

// Note represents a Pinboard note.
type Note struct {
	// Unique ID of the note.
	ID string

	// Title of the note.
	Title string

	// 20 character long sha1 hash of the note text.
	Hash []byte

	// Time the note was created.
	CreatedAt time.Time

	// Time the note was updated.
	UpdatedAt time.Time

	// Character length of the note.
	Length int

	// Body text of the note.
	//
	// Note: only /notes/ID returns body text.
	Text []byte
}

// note holds intermediate data for preprocessing types because JSON
// is so cool.
type note struct {
	ID        string      `json:"id,omitempty"`
	Title     string      `json:"title,omitempty"`
	Hash      string      `json:"hash,omitempty"`
	CreatedAt string      `json:"created_at,omitempty"`
	UpdatedAt string      `json:"updated_at,omitempty"`
	Length    interface{} `json:"length,omitempty"`
	Text      string      `json:"text,omitempty"`
}

// notesResponse holds notes responses from the Pinboard API.
type notesResponse struct {
	Count int
	Notes []note
}

// parseNote takes the note as JSON data returned from the Pinboard
// API and translates them into Notes with proper types.
func parseNote(n note) (*Note, error) {
	var note Note

	note.ID = n.ID
	note.Title = n.Title
	note.Hash = []byte(n.Hash)
	note.Text = []byte(n.Text)

	layout := "2006-01-02 15:04:05"
	created, err := time.Parse(layout, n.CreatedAt)
	if err != nil {
		return nil, err
	}
	note.CreatedAt = created

	updated, err := time.Parse(layout, n.UpdatedAt)
	if err != nil {
		return nil, err
	}
	note.UpdatedAt = updated

	switch v := reflect.ValueOf(n.Length); v.Kind() {
	case reflect.String:
		length, err := strconv.Atoi(n.Length.(string))
		if err != nil {
			return nil, err
		}
		note.Length = length
	case reflect.Float64:
		note.Length = int(n.Length.(float64))
	}

	return &note, nil
}

// NotesList returns a list of the user's notes.
//
// https://pinboard.in/api/#notes_list
func NotesList() ([]*Note, error) {
	resp, err := get("notesList", nil)
	if err != nil {
		return nil, err
	}

	var nr notesResponse
	err = json.Unmarshal(resp, &nr)
	if err != nil {
		return nil, err
	}

	// Parse returned untyped notes into Notes.
	var notes []*Note
	for _, n := range nr.Notes {
		note, err := parseNote(n)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

// notesIDOptions represents the single required argument for
// /notes/ID.
type notesIDOptions struct {
	ID string
}

// NotesID returns an individual user note. The hash property is a 20
// character long sha1 hash of the note text.
//
// https://pinboard.in/api/#notes_get
func NotesID(id string) (*Note, error) {
	resp, err := get("notesID", &notesIDOptions{ID: id})
	if err != nil {
		return nil, err
	}

	var n note
	err = json.Unmarshal(resp, &n)
	if err != nil {
		return nil, err
	}

	note, err := parseNote(n)
	if err != nil {
		return nil, err
	}

	return note, nil
}
