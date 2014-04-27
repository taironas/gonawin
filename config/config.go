package config

import (
	"encoding/json"
	"log"
	"os"
)

// configuration structure to hold the JSON unmarshalled data.
type GwConfig struct {
	ApiVersion        string     `json:"apiVersion"`
	OfflineMode       bool       `json:"offlineMode"`
	OfflineUsers      []User     `json:"offlineUsers"`
	DevUsers          []User     `json:"devUsers"`
	Admins            []string   `json:"admins"`
	Twitter           Twitter    `json:"twitter"`
	Facebook          Facebook   `json:"facebook"`
	GooglePlus        GooglePlus `json:"googlePlus"`
	AuthorizedGmail   []string   `json:"authorizedGmail"`
	AuthorizedTwitter []string   `json:"authorizedTwitter"`
}

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

type Twitter struct {
	Token  string `json:"token"`
	Secret string `json:"secret"`
}

type Facebook struct {
	AppId string `json:"appId"`
}

type GooglePlus struct {
	ClientId string `json:"clientId"`
}

// Read configuration file and return it.
func ReadConfig(filename string) (*GwConfig, error) {

	c := &GwConfig{}
	var f *os.File
	var err error
	if len(filename) > 0 {
		log.Printf("gonawin: gwConfig.load: reading file %s", filename)
		f, err = os.Open(filename)
	} else {
		env, _ := os.Getwd()
		log.Printf("gonawin: gwConfig %v", env)
		log.Printf("gonawin: gwConfig.load: reading file ./config.json")
		f, err = os.Open("./config.json")
	}
	if nil == err {
		log.Printf("gonawin: gwConfig.load: decoding file.")
		decoder := json.NewDecoder(f)
		err = decoder.Decode(&c)
		if err == nil {
			return c, nil
		}
	}
	return nil, err
}
