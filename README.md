# Package pinboard

import "github.com/imwally/pinboard.go"

## Overview

Package pinboard implements a golang wrapper for the
pinboard [api](https://pinboard.in/api/).

## Example

```
// initialise a new pinboard post
var p pinboard.Post

// set pinboard authentication token
p.Token = "username:TOKEN"

// add a new bookmark
p.Url = "http://golang.org"                     // required
p.Description = "The Go Programming Language"   // required

// encode the post into an api GET request URL.
p.Encode()

// add the bookmark
err := p.Add()
if err != nil {
    fmt.Println(err)
}
```
