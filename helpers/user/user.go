/*
 * Copyright (c) 2013 Santiago Arias | Remy Jourde | Carlos Bernal
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"appengine"

	"github.com/santiaago/purple-wing/helpers/log"
)

type GPlusUserInfo struct {
	Id    string
	Email string
	Name  string
}

type TwitterUserInfo struct {
	Id          int64
	Name        string
	Screen_name string
}

type FacebookUserInfo struct {
	Name, Email string
}

type FacebookTokenData struct {
	Data DataType
}

type DataType struct {
	App_id      int
	Application string
	Expires_at  int
	Is_valid    bool
	Issued_at   int
	Metadata    MetadataType
	Scopes      []string
	User_id     int
}

type MetadataType struct {
	Sso string
}

// unmarshal google plus response for googleapi.com request
func FetchGPlusUserInfo(r *http.Request, c *http.Client) (*GPlusUserInfo, error) {
	co := appengine.NewContext(r)
	// Make the request.
	request, err := c.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	log.Infof(co, "FetchGPlusUserInfo request: %q", request)

	if err != nil {
		return nil, err
	}

	defer request.Body.Close()

	if body, err := ioutil.ReadAll(request.Body); err == nil {
		var ui *GPlusUserInfo
		log.Infof(co, "FetchGPlusUserInfo body: %q", body)
		if err := json.Unmarshal(body, &ui); err == nil {
			return ui, err
		}
	}

	log.Infof(co, "FetchGPlusUserInfo err: %v", err)

	return nil, err
}

// unmarshal twitter response
func FetchTwitterUserInfo(r *http.Response) (*TwitterUserInfo, error) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err == nil {
		var ui *TwitterUserInfo

		if err = json.Unmarshal(body, &ui); err == nil {
			return ui, err
		}
	}

	return nil, err
}

// unmarshal facebook response for graph.facebook.com request
func FetchFacebookTokenData(r *http.Response) (*FacebookTokenData, error) {
	defer r.Body.Close()
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	} else {
		var data *FacebookTokenData
		if err = json.Unmarshal(body, &data); err != nil {
			return nil, err
		} else {
			return data, err
		}
	}
}

// unmarshal facebook response from facebook.com/me
func FetchFacebookUserInfo(graphResponse *http.Response) (*FacebookUserInfo, error) {
	defer graphResponse.Body.Close()
	if graphBody, err := ioutil.ReadAll(graphResponse.Body); err != nil {
		return nil, err
	} else {
		var userInfo *FacebookUserInfo
		if err = json.Unmarshal(graphBody, &userInfo); err != nil {
			return nil, err
		} else {
			return userInfo, err
		}
	}
}
