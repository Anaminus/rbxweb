// Deals with services related to ROBLOX sets.
package set

import (
	"github.com/anaminus/rbxweb"
	"net/url"
)

// Add adds an asset to a set. The set must belong to the current user, and
// the asset must be addable to sets.
//
// This function requires the client to be logged in.
func Add(client *rbxweb.Client, assetId int64, setId int32) (err error) {
	query := url.Values{
		"rqtype":  {"addtoset"},
		"assetId": {client.I64toa(assetId)},
		"setId":   {client.I32toa(setId)},
	}

	resp, err := client.Post(client.GetURL(`www`, `/Sets/SetHandler.ashx`, query), "", nil)
	if err = client.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
