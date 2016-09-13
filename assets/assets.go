// +build dev

package assets

import "net/http"

// Assets is a HTTP filesystem containing all of the HTML/JS/CSS.
var Assets = http.Dir("assets/static")
