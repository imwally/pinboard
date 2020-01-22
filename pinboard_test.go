package pinboard

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

var token string

func TestMain(m *testing.M) {
	tokenEnv, ok := os.LookupEnv("PINBOARD_TOKEN")
	if !ok {
		fmt.Println("pinboard token env not set")
		os.Exit(1)
	}

	token = tokenEnv
	os.Exit(m.Run())
}

func PrintStruct(s interface{}) {
	structType := reflect.TypeOf(s)
	structValue := reflect.ValueOf(s)
	structNumFields := structType.NumField()

	for i := 0; i < structNumFields; i++ {
		structField := structType.Field(i)
		structValue := structValue.Field(i)
		fmt.Println(structField.Name+":", structValue)
	}
}

func TestAdd(t *testing.T) {
	p := Post{
		Token:       token,
		Description: "TESTING PINBOARD CLI CLIENT",
		URL:         "https://github.com/imwally/pinboard",
		Tags:        "pin pinboard test testing",
		Extended:    "This is a test from imwally's golang pinboard package. For more information please refer to the pinned URL.",
		Toread:      "yes",
		Shared:      "yes",
	}

	t.Log("Adding: ", p.URL)
	err := p.Add()
	if err != nil {
		t.Error(err)
	}
}

func TestShow(t *testing.T) {
	p := Post{
		Token: token,
		Count: 2,
	}

	recent := p.ShowRecent()
	for _, r := range recent.Posts {
		PrintStruct(r)
		fmt.Println("---")
	}
}

func TestDelete(t *testing.T) {
	p := Post{
		Token: token,
		URL:   "https://github.com/imwally/pinboard",
	}

	t.Log("Deleting: ", p.URL)
	err := p.Delete()
	if err != nil {
		t.Error(err)
	}
}
