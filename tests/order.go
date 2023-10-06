package tests

import "github.com/paveldanilin/ginx/requestbody"

type order struct {
	// A marker for binding a request body
	requestbody.JSON

	// Resolve value by a request query string.
	Extra string `ginx:"query=extra" json:"-"`

	// These values will be resolver by a request body.

	ID      int    `json:"id"`
	Name    string `json:"name"`
	Product string `json:"product"`
}
