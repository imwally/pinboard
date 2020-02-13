// Package pinboard provides a wrapper for accessing the Pinboard API.
//
// https://pinboard.in/api/
//
// All Pinboard API methods are fully supported.
//
// Function names mirror the API endpoints. For example:
//
//     PostsAdd() calls the /posts/add method
//     TagsDelete() calls the /tags/delete method
//
// If a method supports optional arguments then a MethodOptions struct
// allows you to specify those options to pass to said method. For
// example:
//
//     PostsAdd(&PostsAddOptions{})
//     PostsGet(&PostsGetOptions{})
//
// Not all endpoints require arguments, in which case just pass nil.
//
//     PostsAll(nil)
package pinboard

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	api    string = "https://api.pinboard.in/"
	ver    string = "v1"
	apiurl string = api + ver
)

var (
	// Supported Pinboard API methods that map to endpoints.
	endpoints = map[string]string{
		"postsUpdate":  "/posts/update",
		"postsAdd":     "/posts/add",
		"postsDelete":  "/posts/delete",
		"postsGet":     "/posts/get",
		"postsRecent":  "/posts/recent",
		"postsDates":   "/posts/dates",
		"postsAll":     "/posts/all",
		"postsSuggest": "/posts/suggest",
		"tagsGet":      "/tags/get",
		"tagsRename":   "/tags/rename",
		"tagsDelete":   "/tags/delete",
		"userSecret":   "/user/secret",
		"userAPIToken": "/user/api_token",
		"notesList":    "/notes/list",
		"notesID":      "/notes/",
	}

	pinboardToken = ""
)

// get checks if endpoint is a valid Pinboard API endpoint and then
// constructs a valid endpoint URL including the required 'auth_token'
// and 'format' values along with any optional arguments found in the
// options interface. It makes a http.Get request, checks HTTP status
// codes and then finally returns the response body.
func get(endpoint string, options interface{}) (body []byte, err error) {
	ep, ok := endpoints[endpoint]
	if !ok {
		return nil, fmt.Errorf("error: %s is not a supported endpoint", endpoint)
	}

	u, err := url.Parse(apiurl + ep)
	if err != nil {
		return nil, err
	}

	// Set URL query parameters based on the MethodOptions only if
	// options is not nil.
	ov := reflect.ValueOf(options)
	if ov.Kind() == reflect.Ptr && !ov.IsNil() {
		// /notes/ID hack
		if endpoint == "notesID" {
			idOptions := reflect.Indirect(reflect.ValueOf(options))
			id := idOptions.Field(0).String()
			u.Path = u.Path + id
		} else {
			v, err := values(options)
			if err != nil {
				return nil, err
			}

			u.RawQuery = v.Encode()
		}
	}

	// Add API token and format parameters before making request.
	q := u.Query()
	q.Add("auth_token", pinboardToken)
	q.Add("format", "json")
	u.RawQuery = q.Encode()

	// Call APImethod with fully constructed URL.
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check the HTTP response status code. This will tell us
	// whether the API token is not set (401) or if we somehow
	// managed to request an invalid endpoint (500).
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: http %d", res.StatusCode)
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// values expects a *MethodOptions struct and encodes the fields into
// url.Values.
func values(i interface{}) (url.Values, error) {
	vt := reflect.Indirect(reflect.ValueOf(i)).Type()
	vv := reflect.Indirect(reflect.ValueOf(i))

	uv := url.Values{}

	for j := 0; j < vv.NumField(); j++ {
		fName := strings.ToLower(vt.Field(j).Name)
		fType := vt.Field(j).Type
		fValue := vv.Field(j)

		switch fType.Kind() {

		// No need to anything special with strings.
		case reflect.String:
			uv.Add(fName, fValue.String())

		case reflect.Int:

			// Check to make sure we don't have the zero
			// value first.
			if fValue.Interface().(int) != 0 {
				uv.Add(fName, strconv.Itoa(fValue.Interface().(int)))
			}

		// Slices may be of type byte or type string, so
		// process accordingly.
		case reflect.Slice:
			if fValue.Len() > 0 {

				// Check what kind of slice we have.
				switch fValue.Index(0).Kind() {

				// byte slice, add as a string
				case reflect.Uint8:
					uv.Add(fName, string(fValue.Interface().([]uint8)))

				// string slice, create single space delimted
				// string
				case reflect.String:
					spaceDelimted := ""
					for si := 0; si < fValue.Len(); si++ {
						spaceDelimted += fValue.Index(si).Interface().(string) + " "
					}
					uv.Add(fName, strings.TrimRight(spaceDelimted, " "))
				}
			}

		// Bool's are represented as yes/no strings.
		case reflect.Bool:
			if fValue.Bool() {
				uv.Add(fName, "yes")
			} else {
				uv.Add(fName, "no")
			}

		// Process various structs according to their
		// underlying type.
		case reflect.Struct:
			if fType.String() == "time.Time" {
				// Even though we hit a time.Time
				// field, make sure we have something
				// other than the zero value before
				// adding it to the url values,
				// otherwise the zero value will be
				// added.
				timeField := fValue.Interface().(time.Time)
				if !timeField.IsZero() {
					dt := timeField.Format(time.RFC3339)
					uv.Add(fName, dt)
				}
			}
		}
	}

	return uv, nil
}

// SetToken sets the API token required to make API calls. The token
// is expected to be the full string "name:random".
func SetToken(token string) {
	pinboardToken = token
}
