package rbxweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
)

var ErrLoggedIn = errors.New("client is already logged in")

// Login logs the client into a user account on the website. This is
// neccessary for many API functions to properly execute.
func (client *Client) Login(username string, password string) (err error) {
	Values := url.Values{}
	Values.Set("username", username)
	Values.Set("password", password)
	bd := Values.Encode()

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

	req, _ := http.NewRequest("POST", client.GetSecureURL(`api`, `/v2/login`, nil), bytes.NewReader([]byte(bd)))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response data
	// {"userId":12345} on success,
	// {"errors":[{"code":1,"message":"Incorrect password or username. Please try again."}]} on failure
	respData := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return errors.New("Login failed. JSON decode failed. " + err.Error())
	}

	if _, ok := respData["userId"].(float64); ok == false {
		var loginError map[string]interface{}
		if respData["errors"] != nil {
			loginError = respData["errors"].([]interface{})[0].(map[string]interface{})
		} else {
			loginError = respData
		}
		return errors.New("Login failed. Error code " + strconv.FormatInt(int64(loginError["code"].(float64)), 10) + ": \"" + loginError["message"].(string) + "\"")
	}
	resp.Body.Close()
	return nil
}

// Logout logs the client of out of the current user account.
func (client *Client) Logout() (err error) {
	req, _ := http.NewRequest("POST", client.GetSecureURL(`api`, `/sign-out/v1`, nil), nil)
	resp, err := client.Do(req)
	if err = client.AssertResp(resp, err); err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
