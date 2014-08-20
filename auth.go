package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var ErrLoggedIn = errors.New("client is already logged in")

// Login logs the client into a user account on the website. This is
// neccessary for many API functions to properly execute.
func (client *Client) Login(username string, password string) (err error) {
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

	// Ensure the client has a cookiejar
	if client.Jar == nil {
		client.Jar, _ = cookiejar.New(&cookiejar.Options{})
	}
	// Check if the client is already logged in
	domain, _ := url.Parse(client.GetURL(`www`, ``, nil))
	cookies := client.Jar.Cookies(domain)
	for _, cookie := range cookies {
		if cookie.Name == ".ROBLOSECURITY" {
			// Client is already logged in
			return ErrLoggedIn
		}
	}

	req, _ := http.NewRequest("POST", client.GetSecureURL(`www`, `/Services/Secure/LoginService.asmx/ValidateLogin`, nil), bytes.NewReader(bd))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
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
func (client *Client) Logout() (err error) {
	req, _ := http.NewRequest("POST", client.GetSecureURL(`www`, `/authentication/logout`, nil), nil)
	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
