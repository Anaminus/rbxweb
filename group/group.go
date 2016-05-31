// Deals with services related to ROBLOX groups.
package group

import (
	"github.com/anaminus/rbxweb"
	"net/url"
	"strconv"
)

// Shout sets the status message of a given group. The account the client is
// logged into must have a group role that has permission to set the group's
// status.
//
// This function requires the client to be logged in.
func Shout(client *rbxweb.Client, groupID int32, message string) (success bool) {
	page := client.GetURL(`www`, `/My/Groups.aspx`, url.Values{"gid": {strconv.FormatInt(int64(groupId), 10)}})
	err := client.DoRawPost(page, url.Values{
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$GroupStatusPane$StatusTextBox":               {message},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$GroupStatusPane$StatusSubmitButton":          {},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$rbxGroupRoleSetMembersPane$currentRoleSetID": {},
	})
	if err != nil {
		return false
	}
	return true
}

// Wall posts a message to a group wall. The account the client is logged into
// must have a group role that has permission to post to the group's wall.
//
// This function requires the client to be logged in.
func Wall(client *rbxweb.Client, groupID int32, message string) (success bool) {
	page := client.GetURL(`www`, `/My/Groups.aspx`, url.Values{"gid": {strconv.FormatInt(int64(groupID), 10)}})
	err := client.DoRawPost(page, url.Values{
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$GroupWallPane$NewPost":                       {message},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$GroupWallPane$NewPostButton":                 {},
		"ctl00$ctl00$cphRoblox$cphMyRobloxContent$rbxGroupRoleSetMembersPane$currentRoleSetID": {},
	})
	if err != nil {
		return false
	}
	return true
}
