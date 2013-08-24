package rbxweb

import (
	"net/http"
	"net/url"
	"rbxweb/core"
)

// AddToSet adds an asset to a set. The set must belong to the current user,
// and the asset must be addable to sets.
//
// This function requires the client to be logged in.
func AddToSet(client *http.Client, assetId int64, setId int32) (err error) {
	query := url.Values{
		"rqtype":  {"addtoset"},
		"assetId": {core.I64toa(assetId)},
		"setId":   {core.I32toa(setId)},
	}

	resp, err := client.Post(core.GetURL(`www`, `/Sets/SetHandler.ashx`, query), "", nil)
	if err = core.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
