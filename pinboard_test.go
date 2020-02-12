package pinboard

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var optAdd *PostsAddOptions

// Can't test anything without proper authentication.
func TestMain(m *testing.M) {
	tokenEnv, ok := os.LookupEnv("PINBOARD_TOKEN")
	if !ok {
		fmt.Println("PINBOARD_TOKEN env not set")
		os.Exit(1)
	}

	SetToken(tokenEnv)

	opt, err := testPostAddOptions()
	if err != nil {
		fmt.Println("could not create test post options")
		os.Exit(1)
	}

	optAdd = opt

	os.Exit(m.Run())
}

func testPostAddOptions() (*PostsAddOptions, error) {
	dt, err := time.Parse(time.RFC3339, "2010-12-11T19:48:02Z")
	if err != nil {
		return nil, err
	}

	// Post for testing functions. Usually this will be added and
	// removed within each test. Make sure Replace is true
	// otherwise "item already exists" errors will ensue.
	testPost := PostsAddOptions{
		URL:         "https://github.com/imwally/pinboard",
		Description: "Testing Pinboard Go Package",
		Extended:    []byte("This is a test from imwally's golang pinboard package. For more information please refer to the pinned URL."),
		Tags:        []string{"pin", "pinboard", "test", "testing", "pinboard_1_testing", "pinboard_testing"},
		Dt:          dt,
		Toread:      true,
		Shared:      true,
		Replace:     true,
	}

	return &testPost, nil
}
