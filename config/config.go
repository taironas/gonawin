// Package config provides a way to configure gonawin app.
// It reads a config.json file similar to github.com/taironas/gonawin/example-config.json
//
package config

import (
	"encoding/json"
	"os"
)

// GwConfig is configuration structure to hold the JSON unmarshalled data.
//
type GwConfig struct {
	APIVersion  string     `json:"apiVersion"`
	OfflineMode bool       `json:"offlineMode"`
	OfflineUser User       `json:"offlineUser"`
	DevUsers    []User     `json:"devUsers"`
	Twitter     Twitter    `json:"twitter"`
	Facebook    Facebook   `json:"facebook"`
	GooglePlus  GooglePlus `json:"googlePlus"`
}

// User is the user sturcture used for authentication.
//
type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// Twitter holds data needed for Twitter authentication.
//
type Twitter struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

// Facebook holds data needed for Facebook authentication.
//
type Facebook struct {
	AppID string `json:"appId"`
}

// GooglePlus holds data needed for GooglePlus authentication.
//
type GooglePlus struct {
	ClientID string `json:"clientId"`
}

// ReadConfig reads configuration file and return it.
//
func ReadConfig(filename string) (*GwConfig, error) {

	c := &GwConfig{}
	var f *os.File
	var err error
	if len(filename) > 0 {
		f, err = os.Open(filename)
	} else {
		f, err = os.Open("./config.json")
	}
	if nil == err {
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&c)
		if err == nil {
			return c, nil
		}
	}
	return nil, err
}
