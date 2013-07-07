package models

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type GoogleUser struct {
	Id string
	Email string
	Name string
	GivenName string
	FamilyName string
}

var CurrentUser *GoogleUser = nil

func FetchUserInfo(r *http.Request, c *http.Client) (*GoogleUser, error) {
	// Make the request.
	request, err := c.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	
	if err != nil {
		return nil, err
	}

	if userInfo, err := ioutil.ReadAll(request.Body); err == nil {
		var u *GoogleUser

		if err := json.Unmarshal(userInfo, &u); err == nil {
			return u, err
		}	
	}

	return nil, err
}