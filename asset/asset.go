package asset

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"rbxweb"
	"rbxweb/core"
)

// GetLatestModel returns the asset ID of the latest model for a given user.
// Useful for retreiving the asset ID of an uploaded model that was just
// uploaded. If `userId` is 0, then the ID of the current user will be used.
func GetLatestModel(client *http.Client, userId int32) (assetId int64, err error) {
	if userId == 0 {
		userId, _ = rbxweb.GetCurrentUserID(client)
	}

	query := url.Values{
		"Category":          {"Models"},
		"SortType":          {"RecentlyUpdated"},
		"IncludeNotForSale": {"true"},
		"ResultsPerPage":    {"1"},
		"CreatorID":         {core.I32toa(userId)},
	}
	resp, err := client.Get(core.GetURL(`api`, `/catalog/json`, query))
	if err = core.AssertResp(resp, err); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	var asset []struct {
		AssetId int64
	}
	if err = dec.Decode(&asset); err != nil {
		return 0, errors.New("JSON decode failed: " + err.Error())
	}
	return asset[0].AssetId, nil
}

// UploadModel uploads data from `reader` to Roblox as a Model asset. If
// updating an existing model, `modelID` should be the ID of the model. If
// `modelID` is 0, then a new model will be uploaded. If uploading a new
// model, `info` can be used to specify information about the model. The
// following parameters are known:
//
//     name           - The name of the model.
//     description    - The model description.
//     genreTypeId    - The model's genre.
//     isPublic       - Whether the model can be taken by other users.
//     allowComments  - Whether users can comment on the model.
//
// In case the model was newly created, UploadModel returns the ID of the
// model.
//
// This function requires the client to be logged in.
func UploadModel(client *http.Client, reader io.Reader, modelID int64, info url.Values) (assetID int64, err error) {
	query := url.Values{
		"assetid": {core.I64toa(modelID)},
		"type":    {"Model"},
		//	"name":          {"Unnamed Model"},
		//	"description":   {""},
		//	"genreTypeId":   {"1"},
		//	"isPublic":      {"False"},
		//	"allowComments": {"False"},
	}
	if info != nil {
		for key, value := range info {
			query[key] = value
		}
	}
	var buf *bytes.Buffer
	buf.ReadFrom(reader)
	req, _ := http.NewRequest("POST", core.GetURL(`www`, `/Data/Upload.ashx`, query), buf)
	req.Header.Set("User-Agent", "Roblox/WinInet")

	resp, err := client.Do(req)
	if err = core.AssertResp(resp, err); err != nil {
		return 0, err
	}
	if modelID == 0 {
		return GetLatestModel(client, 0)
	} else {
		return modelID, nil
	}
}

// UploadModelFile is similar to UploadModel, but gets the data from a file name.
//
// This function requires the client to be logged in.
func UploadModelFile(client *http.Client, filename string, modelID int64, info url.Values) (assetID int64, err error) {
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return 0, err
	}
	defer file.Close()
	return UploadModel(client, file, modelID, info)
}

// UpdatePlace uploads data from `reader` to Roblox as a Place asset.
// `placeID` must be the ID of an existing place. This function cannot create
// a new place.
//
// This function requires the client to be logged in.
func UpdatePlace(client *http.Client, reader io.Reader, placeID int64) (err error) {
	query := url.Values{
		"assetid": {core.I64toa(placeID)},
		"type":    {"Place"},
	}
	var buf *bytes.Buffer
	buf.ReadFrom(reader)
	req, _ := http.NewRequest("POST", core.GetURL(`www`, `/Data/Upload.ashx`, query), buf)
	req.Header.Set("User-Agent", "Roblox")

	resp, err := client.Do(req)
	if err = core.AssertResp(resp, err); err != nil {
		return err
	}
	return nil
}

// UpdatePlaceFile is similar to UpdatePlace, but gets the data from a file name.
//
// This function requires the client to be logged in.
func UpdatePlaceFile(client *http.Client, filename string, placeID int64) (err error) {
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return
	}
	defer file.Close()
	return UpdatePlace(client, file, placeID)
}
