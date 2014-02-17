// The util package provides common utilies for all rbxweb packages.
package util

import (
	"errors"
	"github.com/anaminus/rbxweb"
	"net/http"
	"net/url"
	"strconv"
)

// GetURL constructs a URL using BaseDomain and the given arguments, with HTTP
// as the protocol.
//
// If `subdomain` is not empty, then it is added as the subdomain before the
// base domain. `path` is the part of the URL that appears after the base
// domain. If `query` is not nil, then it is encoded into query parameters and
// added to the end of the URL.
func GetURL(subdomain string, path string, query url.Values) (url string) {
	url = `http://`
	if subdomain != `` {
		url = url + subdomain + `.`
	}
	url = url + rbxweb.BaseDomain + path
	if query != nil {
		url = url + `?` + query.Encode()
	}
	return
}

// GetSecureURL is similar to GetURL, but it uses HTTPS instead of HTTP.
func GetSecureURL(subdomain string, path string, query url.Values) (url string) {
	url = `https://`
	if subdomain != `` {
		url = url + subdomain + `.`
	}
	url = url + rbxweb.BaseDomain + path
	if query != nil {
		url = url + `?` + query.Encode()
	}
	return
}

// AssertResp checks whether a HTTP response errored. Also errors if the
// response has a non-2XX status code.
func AssertResp(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		resp.Body.Close()
		return errors.New(strconv.Itoa(resp.StatusCode) + ": " + http.StatusText(resp.StatusCode))
	}
	return nil
}

// Atoi32 converts a string to an int32.
func Atoi32(s string) (i int32, err error) {
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i64), nil
}

// Atoi64 converts a string to an int64.
func Atoi64(s string) (i int64, err error) {
	return strconv.ParseInt(s, 10, 64)
}

// I32toa converts an int32 to a string.
func I32toa(i int32) (s string) {
	return strconv.FormatInt(int64(i), 10)
}

// I64toa converts an int64 to a string.
func I64toa(i int64) (s string) {
	return strconv.FormatInt(i, 10)
}
