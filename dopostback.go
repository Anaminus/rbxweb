package rbxweb

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"rbxweb/core"
)

func findInput(inputs *goquery.Selection, inputValues *url.Values, name string) {
	s := inputs.Filter(`input[name="` + name + `"]`)
	if v, e := s.Attr(`value`); e {
		inputValues.Set(name, v)
	}
}

// DoPostBack is a lower-level function that allows the user to perform a
// context-based HTML form POST for a given webpage. On Roblox, many older
// types of requests require a validation token, which is sent along with the
// page when viewing it. Many requests are also context-based. That is,
// manipulation of an item requires knowledge of how the item is displayed on
// the page.
//
// When DoPostBack is called, a GET request is first sent to the page, then
// the validation tokens are parsed from the response and added as POST
// parameters automatically.
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
func DoPostBack(client *http.Client, page string, params url.Values) (err error) {
	// Get form data from URL
	resp, err := client.Get(page)
	if err = core.AssertResp(resp, err); err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	inputs := doc.Find(`input`)

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
	if err = core.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
