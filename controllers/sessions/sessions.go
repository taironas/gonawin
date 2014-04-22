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

// Package sessions provides the JSON handlers to handle connections to gonawin app.
package sessions

import (
	"errors"
	"net/http"
	"net/url"

	"appengine"
	"appengine/urlfetch"
	"appengine/user"

	oauth "github.com/garyburd/go-oauth/oauth"

	"github.com/santiaago/purple-wing/helpers"
	authhlp "github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/log"
	"github.com/santiaago/purple-wing/helpers/memcache"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"

	mdl "github.com/santiaago/purple-wing/models"
)

// Set up a configuration for twitter.
var twitterConfig = oauth.Client{
	Credentials:                   oauth.Credentials{Token: "A8vvQmN473iMZONHW8p6Ng", Secret: "P0Z8cGoulSmsI1nSzWXBq2RA8s0rb7GwVfOJeF8gKL0"},
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}
var twitterCallbackURL string = "/j/auth/twitter/callback"

var googleVerifyTokenURL string = "https://www.google.com/accounts/AuthSubTokenInfo?bearer_token"
var facebookVerifyTokenURL string = "https://graph.facebook.com/me?access_token"

// JSON authentication handler
func Authenticate(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	userInfo := authhlp.UserInfo{Id: r.FormValue("id"), Email: r.FormValue("email"), Name: r.FormValue("name")}

	var verifyURL string
	if r.FormValue("provider") == "google" {
		verifyURL = googleVerifyTokenURL
	} else if r.FormValue("provider") == "facebook" {
		verifyURL = facebookVerifyTokenURL
	}

	if !authhlp.CheckUserValidity(r, verifyURL, r.FormValue("access_token")) {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsAccessTokenNotValid)}
	}
	if !authhlp.IsAuthorized(&userInfo) {
		return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeSessionsForbiden)}
	}

	var user *mdl.User
	var err error
	if user, err = mdl.SigninUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsUnableToSignin)}
	}

	// return user
	userData := struct {
		User *mdl.User
	}{
		user,
	}

	return templateshlp.RenderJson(w, c, userData)
}

// JSON authentication for Twitter.
func TwitterAuth(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	credentials, err := twitterConfig.RequestTemporaryCredentials(urlfetch.Client(c), "http://"+r.Host+twitterCallbackURL, nil)
	if err != nil {
		c.Errorf("JsonTwitterAuth, error = %v", err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetTempCredentials)}
	}

	memcache.Set(c, "secret", credentials.Secret)

	// return OAuth token
	oAuthToken := struct {
		OAuthToken string
	}{
		credentials.Token,
	}

	return templateshlp.RenderJson(w, c, oAuthToken)
}

// Twitter Authentication Callback
func TwitterAuthCallback(w http.ResponseWriter, r *http.Request) error {

	http.Redirect(w, r, "http://localhost:8080/ng#/auth/twitter/callback?oauth_token="+r.FormValue("oauth_token")+"&oauth_verifier="+r.FormValue("oauth_verifier"), http.StatusFound)
	return nil
}

// Twitter Authentication Callback
func TwitterUser(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Twitter User handler:"
	var err error
	var user *mdl.User

	log.Infof(c, "%s oauth_verifier = %s", desc, r.FormValue("oauth_verifier"))
	log.Infof(c, "%s oauth_token = %s", desc, r.FormValue("oauth_token"))

	// get the request token
	requestToken := r.FormValue("oauth_token")
	// update credentials with request token and 'secret cookie value'
	var cred oauth.Credentials
	cred.Token = requestToken
	if secret := memcache.Get(c, "secret"); secret != nil {
		cred.Secret = string(secret.([]byte))
	} else {
		log.Errorf(c, "%s cannot get secret value.", desc)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsCannotGetSecretValue)}
	}

	if err = memcache.Delete(c, "secret"); err != nil {
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

	if !authhlp.IsAuthorizedWithTwitter(userInfo) {
		log.Errorf(c, "%s User is not authorized by twitter with the following user information. %v", desc, userInfo)
		return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeSessionsForbiden)}
	}

	if user, err = mdl.SigninUser(w, r, "Username", "", userInfo.Screen_name, userInfo.Name); err != nil {
		log.Errorf(c, "%s Unable to signin user %s. %v", desc, userInfo.Name, err)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsUnableToSignin)}
	}

	// return user
	userData := struct {
		User *mdl.User
	}{
		user,
	}

	return templateshlp.RenderJson(w, c, userData)
}

// JSON authentication for Google Accounts.
func GoogleAccountsAuth(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Google Accounts Authentication handler:"

	u, err := user.CurrentOAuth(c, "")
	if err != nil {
		log.Errorf(c, "%s OAuth Google Authorization header required", desc)
		return &helpers.Unauthorized{Err: errors.New(helpers.ErrorCodeSessionsAuthHeaderRequired)}
	}
	log.Infof(c, "GoogleAuth: user = %v", u)
	userInfo := authhlp.GetUserGoogleInfo(u)

	if !authhlp.IsAuthorized(&userInfo) {
		return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeSessionsForbiden)}
	}
  
  // get OAuthConsumerKey as access token
  var accessToken string
  if accessToken, err = user.OAuthConsumerKey(c); err != nil {
    log.Errorf(c, "%s Cannot get OAuth consumer key: %v", desc, err)
    return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsOAuthConsumerKey)}
  }

	var user *mdl.User
	if user, err = mdl.SigninUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeSessionsUnableToSignin)}
	}

	// return user
	userData := struct {
    AccessToken string
		User *mdl.User
	}{
    accessToken,
		user,
	}

	return templateshlp.RenderJson(w, c, userData)
}
