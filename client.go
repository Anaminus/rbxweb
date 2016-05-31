package rbxweb

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Client embeds a http.Client, and is used with every function that makes a
// request.
//
// BaseDomain is the URL domain to which all requests will be sent. Subdomains
// are handled automatically as a part of API requests. Alternative domains,
// such as gametest, follow a scheme that makes switching domains easier:
//
//     BaseDomain:                  With subdomain:
//     roblox.com               --> www.roblox.com
//     gametest.robloxlabs.com  --> www.gametest.robloxlabs.com
type Client struct {
	http.Client
	BaseDomain string
}

func NewClient() *Client {
	return &Client{
		BaseDomain: "roblox.com",
	}
}

// GetURL constructs a URL using BaseDomain and the given arguments, with HTTP
// as the protocol.
//
// If `subdomain` is not empty, then it is added as the subdomain before the
// base domain. `path` is the part of the URL that appears after the base
// domain. If `query` is not nil, then it is encoded into query parameters and
// added to the end of the URL.
func (client *Client) GetURL(subdomain string, path string, query url.Values) (url string) {
	url = `http://`
	if subdomain != `` {
		url = url + subdomain + `.`
	}
	url = url + client.BaseDomain + path
	if query != nil {
		url = url + `?` + query.Encode()
	}
	return
}

// GetSecureURL is similar to GetURL, but it uses HTTPS instead of HTTP.
func (client *Client) GetSecureURL(subdomain string, path string, query url.Values) (url string) {
	url = `https://`
	if subdomain != `` {
		url = url + subdomain + `.`
	}
	url = url + client.BaseDomain + path
	if query != nil {
		url = url + `?` + query.Encode()
	}
	return
}

// AssertResp checks whether a HTTP response errored. Also errors if the
// response has a non-2XX status code.
func (client *Client) AssertResp(resp *http.Response, err error) error {
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 && resp.StatusCode >= 300 {
		resp.Body.Close()
		return errors.New(strconv.Itoa(resp.StatusCode) + ": " + http.StatusText(resp.StatusCode))
	}
	return nil
}
