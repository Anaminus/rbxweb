package rbxweb

import (
	"code.google.com/p/go.net/html"
	"net/http"
	"net/url"
	"rbxweb/util"
)

func findInput(inputs []*html.Node, inputValues *url.Values, name string) {
	for _, node := range inputs {
		var match bool
		var value string
		for _, attr := range node.Attr {
			if attr.Key == "name" && attr.Val == name {
				match = true
				if value != "" {
					break
				}
			} else if attr.Key == "value" {
				value = attr.Val
				if match {
					break
				}
			}
		}
		if match && value != "" {
			inputValues.Set(name, value)
			break
		}
	}
}

func recurseNode(node *html.Node, f func(*html.Node) bool) {
	for node != nil {
		if f(node) && node.FirstChild != nil {
			recurseNode(node.FirstChild, f)
		}
		node = node.NextSibling
	}
}

// DoRawPost is a lower-level function that allows the user to perform a
// context-based HTML form POST for a given webpage. On Roblox, many older
// types of requests require a validation token, which is sent along with the
// page when viewing it. Many requests are also context-based. That is,
// manipulation of an item requires knowledge of how the item is displayed on
// the page.
//
// When DoRawPost is called, a GET request is first sent to the page, then the
// validation tokens are parsed from the response and added as POST parameters
// automatically.
//
// The values in `params` are used as the parameters for the POST request. If
// a value is empty, then the value is retrieved from the original GET
// response body, in the same way as the validation tokens. Values are
// retrieved by matching the key with the name attribute of any input tags
// found in the GET response body.
//
// `page` must be a full URL, including query values, if required.
//
// Whether the client needs to be logged in varies depending on the request.
func DoRawPost(client *http.Client, page string, params url.Values) (err error) {
	// Get form data from URL
	resp, err := client.Get(page)
	if err = util.AssertResp(resp, err); err != nil {
		return err
	}
	defer resp.Body.Close()

	// Search for all input tags by parsing the response body
	root, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	inputs := make([]*html.Node, 0)
	recurseNode(root.FirstChild, func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			if node.Data == "input" {
				inputs = append(inputs, node)
			}
			return true
		}
		return false
	})

	// Look up validation parameters
	findInput(inputs, &params, `__VIEWSTATE`)
	findInput(inputs, &params, `__EVENTVALIDATION`)

	for name, value := range params {
		if len(value) == 0 || value[0] == "" {
			findInput(inputs, &params, name)
		}
	}

	// Post to URL with parameters
	resp, err = client.PostForm(page, params)
	if err = util.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
