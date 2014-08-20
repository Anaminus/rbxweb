package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

// UserInfo contains information about the current user.
type UserInfo struct {
	UserID                  int32
	UserName                string
	RobuxBalance            int64
	TicketsBalance          int64
	ThumbnailUrl            string
	IsAnyBuildersClubMember bool
}

// GetUserInfo returns information about the current user.
//
// This function requires the client to be logged in.
func GetUserInfo(client *Client) (info UserInfo, err error) {
	resp, err := client.Get(client.GetURL(`www`, `/MobileAPI/UserInfo`, nil))
	if err = client.AssertResp(resp, err); err != nil {
		return info, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&info); err != nil {
		return info, errors.New("JSON decode failed: " + err.Error())
	}
	return info, nil
}

// GetCurrentUserId returns the id of the user currently logged in.
//
// This function requires the client to be logged in.
func GetCurrentUserId(client *Client) (id int32, err error) {
	resp, err := client.Get(client.GetURL(`www`, `/Game/GetCurrentUser.ashx`, nil))
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r bytes.Buffer
	r.ReadFrom(resp.Body)
	id, err = client.Atoi32(r.String())
	if err != nil {
		return 0, errors.New("user is not authorized")
	}
	return id, nil
}

// GetIdFromName returns a user id from a user name.
func GetIdFromName(client *Client, name string) (id int32, err error) {
	if name == "" {
		return 0, errors.New("name not specified")
	}
	query := url.Values{
		"UserName": {name},
	}
	req, _ := http.NewRequest("HEAD", client.GetURL(`www`, `/User.aspx`, query), nil)
	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}
	resp.Body.Close()
	values, err := url.ParseQuery(resp.Header.Get("Location"))
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}
	return client.Atoi32(values.Get("ID"))
}

// GetNameFromId returns a user name from a user id.
func GetNameFromId(client *Client, id int32) (name string, err error) {
	if id == 0 {
		return "", errors.New("id not specified")
	}
	resp, err := client.Get(client.GetURL(`api`, `/users/`+client.I32toa(id), nil))
	if err = client.AssertResp(resp, err); err != nil {
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
