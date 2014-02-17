package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anaminus/rbxweb/util"
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
func GetUserInfo(client *http.Client) (info UserInfo, err error) {
	resp, err := client.Get(util.GetURL(`www`, `/MobileAPI/UserInfo`, nil))
	if err = util.AssertResp(resp, err); err != nil {
		return info, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&info); err != nil {
		return info, errors.New("JSON decode failed: " + err.Error())
	}
	return info, nil
}

// GetCurrentUserID returns the ID of the user currently logged in.
//
// This function requires the client to be logged in.
func GetCurrentUserID(client *http.Client) (id int32, err error) {
	resp, err := client.Get(util.GetURL(`www`, `/Game/GetCurrentUser.ashx`, nil))
	if err = util.AssertResp(resp, err); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r bytes.Buffer
	r.ReadFrom(resp.Body)
	id, err = util.Atoi32(r.String())
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
	req, _ := http.NewRequest("HEAD", util.GetURL(`www`, `/User.aspx`, query), nil)
	resp, err := client.Do(req)
	if err = util.AssertResp(resp, err); err != nil {
		return 0, err
	}
	resp.Body.Close()
	values, err := url.ParseQuery(resp.Header.Get("Location"))
	if err = util.AssertResp(resp, err); err != nil {
		return 0, err
	}
	return util.Atoi32(values.Get("ID"))
}

// GetNameFromID returns a user name from a user ID. If `id` is 0, then the
// name of the current user will be returned.
func GetNameFromID(client *http.Client, id int32) (name string, err error) {
	if id == 0 {
		id, _ = GetCurrentUserID(client)
	}
	resp, err := client.Get(util.GetURL(`api`, `/users/`+util.I32toa(id), nil))
	if err = util.AssertResp(resp, err); err != nil {
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
