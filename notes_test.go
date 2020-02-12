package pinboard

import (
	"bytes"
	"testing"
	"time"
)

// Unfortunately, these tests only work with my (wally) personal
// account. There's no way to add and remove a note through the API.
func TestNotesList(t *testing.T) {
	notes, err := NotesList()
	if err != nil {
		t.Error(err)
	}

	expected := 2
	got := len(notes)

	if got != expected {
		t.Errorf("error: got %v, expected %v notes", got, expected)
	}
}

func TestNotesID(t *testing.T) {
	note, err := NotesID("0eefe8bbf5f69c3595e4")
	if err != nil {
		t.Error(err)
	}

	expectedID := "0eefe8bbf5f69c3595e4"
	gotID := note.ID

	if gotID != expectedID {
		t.Errorf("error: got %v, expected %v notes", gotID, expectedID)
	}

	expectedTitle := "Pinboard Testing"
	gotTitle := note.Title

	if gotTitle != expectedTitle {
		t.Errorf("error: got %v, expected %v notes", gotTitle, expectedTitle)
	}

	expectedHash := []byte("40ab7b7ab0a5448d9b49")
	gotHash := note.Hash

	if bytes.Compare(gotHash, expectedHash) != 0 {
		t.Errorf("error: got %v, expected %v notes", gotHash, expectedHash)
	}

	layout := "2006-01-02 15:04:05"

	expectedCreated, _ := time.Parse(layout, "2020-01-28 20:31:41")
	gotCreated := note.CreatedAt

	if gotCreated != expectedCreated {
		t.Errorf("error: got %v, expected %v notes", gotCreated, expectedCreated)
	}

	expectedUpdated, _ := time.Parse(layout, "2020-01-29 04:40:53")
	gotUpdated := note.UpdatedAt

	if gotUpdated != expectedUpdated {
		t.Errorf("error: got %v, expected %v notes", gotUpdated, expectedUpdated)
	}

	expectedLength := 193
	gotLength := note.Length

	if gotLength != expectedLength {
		t.Errorf("error: got %v, expected %v notes", gotLength, expectedLength)
	}

	expectedText := []byte("This is a note used strictly for testing purposes.\r\n\r\nIt's not a very fancy note.\r\n\r\nBut it does have line breaks.\r\n\r\n## Next Section\r\n\r\nAnd even some __crazy__ markdown.\r\n\r\nNeat.\r\n\r\nAn edit...")
	gotText := note.Text

	if bytes.Compare(gotText, expectedText) != 0 {
		t.Errorf("error: got %v, expected %v notes", gotText, expectedText)
	}
}
