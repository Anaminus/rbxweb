// The set package deals with ROBLOX sets.
package set

import (
	"github.com/anaminus/rbxweb/util"
	"net/http"
	"net/url"
)

// Add adds an asset to a set. The set must belong to the current user, and
// the asset must be addable to sets.
//
// This function requires the client to be logged in.
func Add(client *http.Client, assetId int64, setId int32) (err error) {
	query := url.Values{
		"rqtype":  {"addtoset"},
		"assetId": {util.I64toa(assetId)},
		"setId":   {util.I32toa(setId)},
	}

	resp, err := client.Post(util.GetURL(`www`, `/Sets/SetHandler.ashx`, query), "", nil)
	if err = util.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
