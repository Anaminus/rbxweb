// Data Sizes
// User ID:  int32
// Asset ID: int64
// Group ID: int32
// Set ID:   int32
// Currency: int64

package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"rbxweb/core"
)

// GetCurrentUserID returns just the ID of the user currently logged in.
//
// This function requires the client to be logged in.
func GetCurrentUserID(client *http.Client) (id int32, err error) {
	resp, err := client.Get(core.GetURL(`www`, `/Game/GetCurrentUser.ashx`, nil))
	if err = core.AssertResp(resp, err); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r bytes.Buffer
	r.ReadFrom(resp.Body)
	id, err = core.Atoi32(r.String())
	if err != nil {
		return 0, errors.New("user is not authorized")
	}
	return id, nil
}

// GetIDFromName returns a user ID from a user name. If `name` is empty, then
// the ID of the current user is returned.
func GetIDFromName(client *http.Client, name string) (id int32, err error) {
	if name == "" {
		return GetCurrentUserID(client)
	}
	query := url.Values{
		"UserName": {name},
	}
	req, _ := http.NewRequest("HEAD", core.GetURL(`www`, `/User.aspx`, query), nil)
	resp, err := client.Do(req)
	if err = core.AssertResp(resp, err); err != nil {
		return 0, err
	}
	resp.Body.Close()
	values, err := url.ParseQuery(resp.Header.Get("Location"))
	if err = core.AssertResp(resp, err); err != nil {
		return 0, err
	}
	return core.Atoi32(values.Get("ID"))
}

// GetNameFromID returns a user name from a user ID. If `id` is 0, then the
// name of the current user will be returned.
func GetNameFromID(client *http.Client, id int32) (name string, err error) {
	if id == 0 {
		id, _ = GetCurrentUserID(client)
	}
	resp, err := client.Get(core.GetURL(`api`, `/users/`+core.I32toa(id), nil))
	if err = core.AssertResp(resp, err); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	var user struct {
		Username string
	}
	if err = dec.Decode(&user); err != nil {
		return "", errors.New("JSON decode failed: " + err.Error())
	}
	return user.Username, nil
}
