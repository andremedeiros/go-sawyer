// Package httpcache provides facilities for caching HTTP responses and
// hypermedia for REST resources.  The saved responses respect HTTP caching
// policies.  Hypermedia shouldn't change, so is stored for as long as possible.
package httpcache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/lostisland/go-sawyer"
	"github.com/lostisland/go-sawyer/hypermedia"
	"net/http"
)

var (
	DefaultAdapter  Adapter
	NoResponseError = errors.New("No Response")
)

type Adapter interface {
	// Get retrieves a Response for a REST resource by its URL.  The URL should be
	// the full canonical URL for the resource.  The response will be nil if it is
	// expired.
	Get(*http.Request, interface{}) *sawyer.Response

	// Set caches a Response for a resource by its URL.
	Set(*http.Request, *sawyer.Response, interface{}) error

	// Rels retrieves the hypmermedia for a REST resource by its URL.  The relations
	// cache doesn't expire with the assumption that web services will provide
	// redirects as URLs change.
	Rels(*http.Request) hypermedia.Relations
}

func ResponseError(err error) *sawyer.Response {
	return sawyer.ResponseError(err)
}

func EmptyResponse() *sawyer.Response {
	return ResponseError(NoResponseError)
}

func RequestKey(r *http.Request) string {
	return r.Header.Get(keyHeader) + keySep + r.URL.String()
}

func RequestSha(r *http.Request) string {
	key := RequestKey(r)
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[0:32])
}

const (
	keySep    = ":"
	keyHeader = "Accept"
)
