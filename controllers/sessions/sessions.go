/*
 * Copyright (c) 2014 Santiago Arias | Remy Jourde
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

// Package sessions provides the JSON handlers to handle connections to gonawin app.
package sessions

import (
	"errors"
	golog "log"
	"net/http"
	"net/url"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"appengine/user"

	oauth "github.com/garyburd/go-oauth/oauth"

	"github.com/taironas/gonawin/helpers"
	authhlp "github.com/taironas/gonawin/helpers/auth"
	"github.com/taironas/gonawin/helpers/log"
	"github.com/taironas/gonawin/helpers/memcache"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	gwconfig "github.com/taironas/gonawin/config"
	mdl "github.com/taironas/gonawin/models"
)

var (
	config                 *gwconfig.GwConfig
	twitterConfig          oauth.Client
	twitterCallbackURL     string
	googleVerifyTokenURL   string
	facebookVerifyTokenURL string
)

func init() {
	// read config file.
	var err error
	if config, err = gwconfig.ReadConfig(""); err != nil {
		golog.Printf("Error: unable to read config file; %v", err)
	}
	// Set up a configuration for twitter.
	twitterConfig = oauth.Client{
		Credentials:                   oauth.Credentials{Token: config.Twitter.Token, Secret: config.Twitter.Secret},
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}
	twitterCallbackURL = "/j/auth/twitter/callback"
	googleVerifyTokenURL = "https://www.google.com/accounts/AuthSubTokenInfo?bearer_token"
	facebookVerifyTokenURL = "https://graph.facebook.com/me?access_token"
}

// Authenticate handler, use it to authenticate a user.
// It returns the JSON data of the requested user.
func Authenticate(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)

	userInfo := authhlp.UserInfo{ID: r.FormValue("id"), Email: r.FormValue("email"), Name: r.FormValue("name")}

	var verifyURL string
	if r.FormValue("provider") == "google" {
		verifyURL = googleVerifyTokenURL
	} else if r.FormValue("provider") == "facebook" {
		verifyURL = facebookVerifyTokenURL
	}

	if !authhlp.CheckUserValidity(r, verifyURL, r.FormValue("access_token")) {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsAccessTokenNotValid)}
	}

	var user *mdl.User
	var err error
	if user, err = mdl.SigninUser(c, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsUnableToSignin)}
	}

	imageURL := helpers.UserImageURL(user.Username, user.ID)

	userData := struct {
		User     *mdl.User
		ImageURL string
	}{
		user,
		imageURL,
	}

	return templateshlp.RenderJson(w, c, userData)
}

// TwitterAuth handler, use it to authenticate via twitter.
func TwitterAuth(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Twitter Auth handler:"

	credentials, err := twitterConfig.RequestTemporaryCredentials(urlfetch.Client(c), "http://"+r.Host+twitterCallbackURL, nil)
	if err != nil {
		c.Errorf("JsonTwitterAuth, error = %v", err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetTempCredentials)}
	}

	if err = memcache.Set(c, "secret", credentials.Secret); err != nil {
		// store secret in datastore
		secretId, _, err := datastore.AllocateIDs(c, "Secret", nil, 1)
		if err != nil {
			log.Errorf(c, "%s Cannot allocate ID for secret. %v", desc, err)
		}

		key := datastore.NewKey(c, "Secret", "", secretId, nil)

		_, err = datastore.Put(c, key, credentials.Secret)
		if err != nil {
			log.Errorf(c, "%s Cannot put secret in Datastore. %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotSetSecretValue)}
		}
	}

	// return OAuth token
	oAuthToken := struct {
		OAuthToken string
	}{
		credentials.Token,
	}

	return templateshlp.RenderJson(w, c, oAuthToken)
}

// TwitterAuthCallback handler, use it to make a callback for Twitter Authentication.
func TwitterAuthCallback(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	http.Redirect(w, r, "http://"+r.Host+"/#/auth/twitter/callback?oauth_token="+r.FormValue("oauth_token")+"&oauth_verifier="+r.FormValue("oauth_verifier"), http.StatusFound)
	return nil
}

// TwitterUser handler, use it to get the Twitter user data.
func TwitterUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Twitter User handler:"

	var user *mdl.User

	// get the request token
	requestToken := r.FormValue("oauth_token")

	// update credentials with request token and 'secret cookie value'
	var cred oauth.Credentials
	cred.Token = requestToken
	if secret, err := memcache.Get(c, "secret"); secret != nil {
		cred.Secret = string(secret.([]byte))
	} else {
		log.Errorf(c, "%s cannot get secret value from memcache: %v", desc, err)
		// try to get secret from datastore
		q := datastore.NewQuery("Secret")
		var secrets []string
		if keys, err := q.GetAll(c, &secrets); err == nil && len(secrets) > 0 {
			// delete secret from datastore
			if err = datastore.Delete(c, keys[0]); err != nil {
				log.Errorf(c, "%s Error when trying to delete 'secret' key in Datastore: %v", desc, err)
			}

		} else if err != nil || len(secrets) == 0 {
			log.Errorf(c, "%s cannot get secret value from Datastore: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetSecretValue)}
		}
	}

	if err := memcache.Delete(c, "secret"); err != nil {
		log.Errorf(c, "%s Error when trying to delete memcached 'secret' key: %v", desc, err)
	}

	token, values, err := twitterConfig.RequestToken(urlfetch.Client(c), &cred, r.FormValue("oauth_verifier"))
	if err != nil {
		log.Errorf(c, "%s Error when trying to delete memcached 'secret' key: %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetRequestToken)}
	}

	// get user info
	urlValues := url.Values{}
	urlValues.Set("user_id", values.Get("user_id"))
	resp, err := twitterConfig.Get(urlfetch.Client(c), token, "https://api.twitter.com/1.1/users/show.json", urlValues)
	if err != nil {
		log.Errorf(c, "%s Cannot get user info from twitter. %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetUserInfo)}
	}

	userInfo, err := authhlp.FetchTwitterUserInfo(resp)
	if err != nil {
		log.Errorf(c, "%s Cannot get user info by fetching twitter response. %v", desc, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetUserInfo)}
	}

	if user, err = mdl.SigninUser(c, "Username", "", userInfo.ScreenName, userInfo.Name); err != nil {
		log.Errorf(c, "%s Unable to signin user %s. %v", desc, userInfo.Name, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsUnableToSignin)}
	}

	imageURL := helpers.UserImageURL(user.Username, user.ID)

	userData := struct {
		User     *mdl.User
		ImageURL string
	}{
		user,
		imageURL,
	}

	return templateshlp.RenderJson(w, c, userData)
}

// GoogleAccountsLoginURL handler, use it to get Google accounts login URL.
func GoogleAccountsLoginURL(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Google Accounts Login URL Handler:"

	var url string
	var err error
	url, err = user.LoginURL(c, "/j/auth/google/callback/")
	if err != nil {
		log.Errorf(c, "%s error when getting Google accounts login URL", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetGoogleLoginURL)}
	}

	loginData := struct {
		Url string
	}{
		url,
	}

	return templateshlp.RenderJson(w, c, loginData)
}

// GoogleAuthCallback handler, use it to make a callback for Google Authentication.
func GoogleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Google Accounts Auth Callback Handler:"

	u := user.Current(c)
	if u == nil {
		log.Errorf(c, "%s user cannot be nil", desc)
		return &helpers.InternalServerError{Err: errors.New("user cannot be nil")}
	}

	http.Redirect(w, r, "http://"+r.Host+"/#/auth/google/callback?auth_token="+u.ID, http.StatusFound)
	return nil
}

// GoogleUser handler, use it to get Google accounts user.
func GoogleUser(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Google Accounts User Handler:"

	u := user.Current(c)
	if u == nil {
		log.Errorf(c, "%s user cannot be nil", desc)
		return &helpers.InternalServerError{Err: errors.New("user cannot be nil")}
	}

	if u.ID != r.FormValue("auth_token") {
		log.Errorf(c, "%s Auth token doesn't match user ID", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsAccessTokenNotValid)}
	}

	userInfo := authhlp.GetUserGoogleInfo(u)

	var user *mdl.User
	var err error
	if user, err = mdl.SigninUser(c, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsUnableToSignin)}
	}

	imageURL := helpers.UserImageURL(user.Username, user.ID)

	userData := struct {
		User     *mdl.User
		ImageURL string
	}{
		user,
		imageURL,
	}

	return templateshlp.RenderJson(w, c, userData)
}

// GoogleDeleteCookie handler, use it to delete cookie created by Google account.
func GoogleDeleteCookie(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)

	cookieName := "ACSID"
	if appengine.IsDevAppServer() {
		cookieName = "dev_appserver_login"
	}
	cookie := http.Cookie{Name: cookieName, Path: "/", MaxAge: -1}
	http.SetCookie(w, &cookie)

	return templateshlp.RenderJson(w, c, "Google user has been logged out")
}

// AuthServiceIds handler, use it to get the identifiers of Gonawin
func AuthServiceIds(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}
	c := appengine.NewContext(r)

	data := struct {
		GooglePlusClientId string
		FacebookAppId      string
	}{
		config.GooglePlus.ClientId,
		config.Facebook.AppId,
	}
	return templateshlp.RenderJson(w, c, data)
}
