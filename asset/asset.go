// Deals with services related to ROBLOX assets.
package asset

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anaminus/rbxweb"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Asset types.
const (
	TypeImage        byte = 1
	TypeTShirt       byte = 2
	TypeAudio        byte = 3
	TypeMesh         byte = 4
	TypeLua          byte = 5
	TypeHTML         byte = 6
	TypeText         byte = 7
	TypeHat          byte = 8
	TypePlace        byte = 9
	TypeModel        byte = 10
	TypeShirt        byte = 11
	TypePants        byte = 12
	TypeDecal        byte = 13
	TypeAvatar       byte = 16
	TypeHead         byte = 17
	TypeFace         byte = 18
	TypeGear         byte = 19
	TypeBadge        byte = 21
	TypeGroupEmblem  byte = 22
	TypeAnimation    byte = 24
	TypeArms         byte = 25
	TypeLegs         byte = 26
	TypeTorso        byte = 27
	TypeRightArm     byte = 28
	TypeLeftArm      byte = 29
	TypeLeftLeg      byte = 30
	TypeRightLeg     byte = 31
	TypePackage      byte = 32
	TypeYouTubeVideo byte = 33
	TypeGamePass     byte = 34
	TypeApp          byte = 35
	TypeCode         byte = 37
	TypePlugin       byte = 38
)

// GetLatestModel returns the asset id of the latest model for a given user.
// While this is useful for retrieving the asset id of a newly created model
// that was just uploaded, it is not necessarily reliable for this purpose.
func GetLatestModel(client *rbxweb.Client, userId int32) (assetId int64, err error) {
	if userId == 0 {
		return 0, errors.New("invalid user id")
	}

	query := url.Values{
		"Category":          {"Models"},
		"SortType":          {"RecentlyUpdated"},
		"IncludeNotForSale": {"true"},
		"ResultsPerPage":    {"1"},
		"CreatorID":         {client.I32toa(userId)},
	}
	resp, err := client.Get(client.GetURL(`api`, `/catalog/json`, query))
	if err = client.AssertResp(resp, err); err != nil {
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

// GetIdFromVersion returns an asset id from an asset version id.
func GetIdFromVersion(client *rbxweb.Client, assetVersionId int64) (assetId int64, err error) {
	query := url.Values{
		"avid": {client.I64toa(assetVersionId)},
	}

	// This relies on how asset names are converted to url names. Currently,
	// if an asset name is "_", its url becomes "unnamed".
	req, _ := http.NewRequest("HEAD", client.GetURL(`www`, `/_-item`, query), nil)
	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}
	resp.Body.Close()

	values, err := url.ParseQuery(resp.Header.Get("Location"))
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}

	return client.Atoi64(values.Get("id"))
}

// Upload generically uploads data from `reader` as an asset to the ROBLOX
// website. `info` can be used to specify information about the model. The
// following parameters are known:
//
//     type           - The type of asset.
//     assetid        - The id of the asset to update. 0 uploads a new asset.
//     name           - The name of the asset.
//     description    - The asset's description.
//     genreTypeId    - The asset's genre.
//     isPublic       - Whether the asset can be taken by other users.
//     allowComments  - Whether users can comment on the asset.
//
// The success of this function is highly dependent on these parameters. For
// example, most asset types may only be uploaded by authorized users.
// Parameters that specify information about the asset only apply for new
// assets. That is, updating an asset will only update the contents, but not
// the information about it.
//
// `assetVersionId` is the version id of the uploaded asset. This is unique
// for each upload. This can be used with GetIdFromVersion to get the asset
// id.
//
// This function requires the client to be logged in.
func Upload(client *rbxweb.Client, reader io.Reader, info url.Values) (assetVersionId int64, err error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	req, _ := http.NewRequest("POST", client.GetURL(`www`, `/Data/Upload.ashx`, info), buf)
	req.Header.Set("User-Agent", "roblox/rbxweb")

	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	r := new(bytes.Buffer)
	r.ReadFrom(resp.Body)
	assetVersionId, _ = client.Atoi64(r.String())

	return assetVersionId, err
}

// UploadModel uploads data from `reader` to Roblox as a Model asset. If
// updating an existing model, `modelId` should be the id of the model. If
// `modelId` is 0, then a new model will be uploaded. If uploading a new
// model, `info` can be used to specify information about the model.
//
// This function requires the client to be logged in.
func UploadModel(client *rbxweb.Client, reader io.Reader, modelId int64, info url.Values) (assetVersionId int64, err error) {
	query := url.Values{
		"assetid": {client.I64toa(modelId)},
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

	return assetVersionId, err
}

// UploadModelFile is similar to UploadModel, but gets the data from a file
// name.
//
// This function requires the client to be logged in.
func UploadModelFile(client *rbxweb.Client, filename string, modelId int64, info url.Values) (assetVersionId int64, err error) {
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return 0, err
	}
	defer file.Close()
	return UploadModel(client, file, modelId, info)
}

// UpdatePlace uploads data from `reader` to Roblox as a Place asset.
// `placeId` must be the id of an existing place. This function cannot create
// a new place.
//
// This function requires the client to be logged in.
func UpdatePlace(client *rbxweb.Client, reader io.Reader, placeId int64) (err error) {
	query := url.Values{
		"assetid": {client.I64toa(placeId)},
		"type":    {"Place"},
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	req, _ := http.NewRequest("POST", client.GetURL(`www`, `/Data/Upload.ashx`, query), buf)
	req.Header.Set("User-Agent", "Roblox")

	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// UpdatePlaceFile is similar to UpdatePlace, but gets the data from a file name.
//
// This function requires the client to be logged in.
func UpdatePlaceFile(client *rbxweb.Client, filename string, placeId int64) (err error) {
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return
	}
	defer file.Close()
	return UpdatePlace(client, file, placeId)
}

// Contains information about an asset.
type Info struct {
	AssetId     int64
	ProductId   int64
	Name        string
	Description string
	AssetTypeId int32
	Creator     struct {
		Id   int32
		Name string
	}
	Created                time.Time
	Updated                time.Time
	PriceInRobux           int64
	PriceInTickets         int64
	Sales                  int32
	IsNew                  bool
	IsForSale              bool
	IsPublicDomain         bool
	IsLimited              bool
	IsLimitedUnique        bool
	Remaining              int32
	MinimumMembershipLevel int32
	ContentRatingTypeId    int32
}

// GetInfo returns information about an asset, given an asset id.
func GetInfo(client *rbxweb.Client, id int64) (info Info, err error) {
	query := url.Values{
		"assetId": {client.I64toa(id)},
	}
	resp, err := client.Get(client.GetURL(`api`, `/marketplace/productinfo`, query))
	if err = client.AssertResp(resp, err); err != nil {
		return Info{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	dec.Decode(&info)
	return
}
