package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"rbxweb/core"
)

// Info is used to contain information about a user.
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
func GetInfo(client *http.Client) (info Info, err error) {
	resp, err := client.Get(core.GetURL(`www`, `/MobileAPI/UserInfo`, nil))
	if err = core.AssertResp(resp, err); err != nil {
		return info, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(&info); err != nil {
		return info, errors.New("JSON decode failed: " + err.Error())
	}
	return info, nil
}
