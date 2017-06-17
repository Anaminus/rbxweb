// Deals with services related to users.
package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"github.com/anaminus/rbxweb"
)

// Info contains information about the current user.
type Info struct {
	UserID                  int32
	UserName                string
	RobuxBalance            int64
	TicketsBalance          int64
	ThumbnailUrl            string
	IsAnyBuildersClubMember bool
}

// GetInfo returns information about the current user.
//
// This function requires the client to be logged in.
func GetInfo(client *rbxweb.Client) (info Info, err error) {
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

// GetCurrentId returns the id of the user currently logged in.
//
// This function requires the client to be logged in.
func GetCurrentId(client *rbxweb.Client) (id int32, err error) {
	resp, err := client.Get(client.GetURL(`www`, `/Game/GetCurrentUser.ashx`, nil))
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r bytes.Buffer
	r.ReadFrom(resp.Body)
	n, err := strconv.ParseInt(r.String(), 10, 32)
	if err != nil {
		return 0, errors.New("user is not authorized")
	}
	return int32(n), nil
}

// GetIdFromName returns a user id from a user name.
func GetIdFromName(client *rbxweb.Client, name string) (id int32, err error) {
	if name == "" {
		return 0, errors.New("name not specified")
	}
	query := url.Values{
		"username": {name},
	}
	req, _ := http.NewRequest("GET", client.GetSecureURL(`api`, `/users/get-by-username`, query), nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}

	var userIdResult struct {
		Id float64
		Username string
	}
	err = json.NewDecoder(resp.Body).Decode(&userIdResult)
	if err != nil {
		return 0, errors.New("JSON decode failed: " + err.Error())
	}

	if userIdResult.Id == 0 {
		return 0, errors.New("user doesn't exist")
	}
	return int32(userIdResult.Id), nil
}

// GetNameFromId returns a user name from a user id.
func GetNameFromId(client *rbxweb.Client, id int32) (name string, err error) {
	if id == 0 {
		return "", errors.New("id not specified")
	}
	resp, err := client.Get(client.GetURL(`api`, `/users/`+strconv.FormatInt(int64(id), 10), nil))
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
