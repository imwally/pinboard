# Package pinboard

import "github.com/imwally/pinboard.go"

## Overview

Package pinboard implements a golang wrapper for the
pinboard [api](https://pinboard.in/api/).

## Example

```go
// Initialise a new pinboard post.
p := new(pinboard.Post)

// Set pinboard authentication token.
p.Token = "username:TOKEN"

// Createa new bookmark.
p.URL = "http://golang.org"                     // required
p.Description = "The Go Programming Language"   // required

// Add the bookmark.
err := p.Add()
if err != nil {
	fmt.Println(err)
}

// Delete the bookmark.
err := p.Delete()
if err != nil {
	fmt.Println(err)
}
```
