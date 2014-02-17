package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anaminus/rbxweb/util"
	"net/http"
)

// Login logs the client into a user account on the website. This is
// neccessary for many API functions to properly execute.
func Login(client *http.Client, username string, password string) (err error) {
	bd, err := json.Marshal(map[string]interface{}{
		"userName":        username,
		"password":        password,
		"isCaptchaOn":     false,
		"challenge":       "",
		"captchaResponse": "",
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", util.GetSecureURL(`www`, `/Services/Secure/LoginService.asmx/ValidateLogin`, nil), bytes.NewReader(bd))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err = util.AssertResp(resp, err); err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response data
	// {"d":{"sl_translate":"Message","IsValid":true,"Message":"","ErrorCode":""}}
	respData := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return errors.New("Login failed. JSON decode failed. " + err.Error())
	}
	d := respData["d"].(map[string]interface{})
	if !d["IsValid"].(bool) {
		return errors.New("Login failed. Error code " + d["ErrorCode"].(string) + ": \"" + d["Message"].(string) + "\"")
	}
	resp.Body.Close()
	return nil
}

// Logout logs the client of out of the current user account.
func Logout(client *http.Client) (err error) {
	req, _ := http.NewRequest("POST", util.GetSecureURL(`www`, `/authentication/logout`, nil), nil)
	resp, err := client.Do(req)
	if err = util.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
