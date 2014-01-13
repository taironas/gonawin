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

package sessions

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"
	
	"appengine"
	"appengine/urlfetch"
	
	oauth "github.com/garyburd/go-oauth/oauth"
	oauth2 "code.google.com/p/goauth2/oauth"

	"github.com/santiaago/purple-wing/helpers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	userhlp "github.com/santiaago/purple-wing/helpers/user"
	
	authhlp "github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/log"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
)

const root string = "/m"
// Set up a configuration for google.
func googleConfig(host string) *oauth2.Config{
	return &oauth2.Config{
		ClientId:     GOOGLE_CLIENT_ID,
		ClientSecret: GOOGLE_CLIENT_SECRET,
		Scope:        "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
		RedirectURL:  fmt.Sprintf("http://%s%s/auth/google/callback", host, root),
	}
}
// Set up a configuration for twitter.
var twitterConfig = oauth.Client{
	Credentials: oauth.Credentials{ Token:	CONSUMER_KEY, Secret: CONSUMER_SECRET },
	TemporaryCredentialRequestURI: "http://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "http://api.twitter.com/oauth/authorize",
	TokenRequestURI:               "http://api.twitter.com/oauth/access_token",
}
var twitterCallbackURL string = "/m/auth/twitter/callback"

const kUrlFacebookMe = "https://graph.facebook.com/me"
const kUrlFacebookDebugToken = "https://graph.facebook.com/debug_token?input_token=%s&access_token=%s|%s"

// Set up a configuration for google.
func facebookConfig(host string) *oauth2.Config{
	return &oauth2.Config{
		ClientId:     FACEBOOK_CLIENT_ID,
		ClientSecret: FACEBOOK_CLIENT_SECRET,
		Scope:"email",
		AuthURL:      "https://graph.facebook.com/oauth/authorize",
		TokenURL:     "https://graph.facebook.com/oauth/access_token",
		RedirectURL:  fmt.Sprintf("http://%s%s/auth/facebook/callback", host, root),
	}
}

// Google
func AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	if !authhlp.LoggedIn(r, c) {
		url := googleConfig(r.Host).AuthCodeURL(r.URL.RawQuery)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}
// Google Authentication Callback
func GoogleAuthCallback(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// Exchange code for an access token at OAuth provider.
	code := r.FormValue("code")
	
	t := &oauth2.Transport{
		Config: googleConfig(r.Host),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(r),
		},
	}
	
	var err error
	var userInfo *userhlp.GPlusUserInfo
	var user *usermdl.User
	
	if _, err = t.Exchange(code); err == nil {
		userInfo, _ = userhlp.FetchGPlusUserInfo(r, t.Client())
	}
	if authhlp.IsAuthorizedWithGoogle(userInfo) {
		if user, err = authhlp.SigninUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil{
			log.Errorf(c, " SigninUser: %v", err)
			http.Redirect(w, r, root, http.StatusFound)
			return
		}
		// store in memcache auth key in memcache
		authhlp.StoreAuthKey(c, user.Id, user.Auth)
		// set 'auth' cookie
		authhlp.SetAuthCookie(w, user.Auth)
	}

	http.Redirect(w, r, root, http.StatusFound)
}

// Twitter
func AuthenticateWithTwitter(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	if !authhlp.LoggedIn(r, c) {
		credentials, err := twitterConfig.RequestTemporaryCredentials(urlfetch.Client(c), "http://" + r.Host + twitterCallbackURL, nil)
		if err != nil {
			log.Errorf(c, " AuthenticateWithTwitter, error getting temporary credentials: %v", err)
			http.Redirect(w, r, root, http.StatusFound)
			return
		}
		
		http.SetCookie(w, &http.Cookie{ Name: "secret", Value: credentials.Secret, Path: "/m", })
		http.Redirect(w, r, twitterConfig.AuthorizationURL(credentials, nil), 302)
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}

// Twitter Authentication Callback
func TwitterAuthCallback(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// get the request token
	requestToken := r.FormValue("oauth_token")
	// update credentials with request token and 'secret cookie value'
	var cred oauth.Credentials
	cred.Token = requestToken
	if cookie, err := r.Cookie("secret"); err == nil {
		cred.Secret = cookie.Value
	} else {
		log.Errorf(c, " TwitterAuthCallback, error getting 'secret' cookie: %v", err)
	}
	
	// clear 'secret' cookie
	http.SetCookie(w, &http.Cookie{ Name: "secret", Path: "/m", Expires: time.Now(), })
	
	token, values, err := twitterConfig.RequestToken(urlfetch.Client(c), &cred, r.FormValue("oauth_verifier"))
	if err != nil {
		log.Errorf(c, " TwitterAuthCallback, error getting request token: %v", err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	
	// get user info
	urlValues := url.Values{}
	urlValues.Set("user_id", values.Get("user_id"))
	resp, err := twitterConfig.Get(urlfetch.Client(c), token, "https://api.twitter.com/1.1/users/show.json", urlValues)
	if err != nil {
		log.Errorf(c, " TwitterAuthCallback, error getting user info from twitter: %v", err)
	}

	userInfo, err := userhlp.FetchTwitterUserInfo(resp)
	if err != nil {
		log.Errorf(c, " TwitterAuthCallback, error occurred when fetching twitter user info: %v", err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}

	if authhlp.IsAuthorizedWithTwitter(userInfo) {
		var user *usermdl.User
		if user, err = authhlp.SigninUser(w, r, "Username", "", userInfo.Screen_name, userInfo.Name); err != nil{
			log.Errorf(c, " SigninUser: %v", err)
			http.Redirect(w, r, root, http.StatusFound)
			return
		}
		// store in memcache auth key in memcache
		authhlp.StoreAuthKey(c, user.Id, user.Auth)
		// set 'auth' cookie
		authhlp.SetAuthCookie(w, user.Auth)
	}

	http.Redirect(w, r, root, http.StatusFound)
}

// Facebook
func AuthenticateWithFacebook(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	if !authhlp.LoggedIn(r, c) {
		url := facebookConfig(r.Host).AuthCodeURL(r.URL.RawQuery)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}

// Facebook authentication callback
func FacebookAuthCallback(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// Exchange code for an access token at OAuth provider.
	code := r.FormValue("code")
	t := &oauth2.Transport{
		Config: facebookConfig(r.Host),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(r),
		},
	}
	
	if v, err := t.Exchange(code); err != nil {
		log.Errorf(c, " FacebookAuthCallback, error occurred during exchange: %v, returned value: %v", err, v)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}

	// build facebook url from debug_token url and get response
	urlFacebook := fmt.Sprintf(kUrlFacebookDebugToken, t.Token.AccessToken, FACEBOOK_CLIENT_ID, FACEBOOK_CLIENT_SECRET)
	accessTokenResponse, err := t.Client().Get(urlFacebook)
	if err != nil {
		log.Errorf(c, " FacebookAuthCallback, client get error calling %s : %v", urlFacebook, err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	// verify if information from facebook is valid for current user.
	if isValid, err := isFacebookTokenValid(accessTokenResponse); (err != nil || !isValid){
		log.Errorf(c, " FacebookAuthCallback, isFacebookTokenValid: Is valid: %v, Error: %v", isValid, err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	// ask facebook.com/me for user profile information
	graphResponse, err := t.Client().Get(kUrlFacebookMe)
	if err != nil {
		log.Errorf(c, " FacebookAuthCallback, failure on get request %s: %v", kUrlFacebookMe, err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	userInfo, err := userhlp.FetchFacebookUserInfo(graphResponse)
	if err != nil{
		log.Errorf(c, " FacebookAuthCallback, error occurred when fetching facebook user info: %v", err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	if authhlp.IsAuthorizedWithFacebook(userInfo){
		var user *usermdl.User
		if user, err = authhlp.SigninUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil{
			log.Errorf(c, " SigninUser: %v", err)
			http.Redirect(w, r, root, http.StatusFound)
			return
		}
		// store in memcache auth key in memcache
		authhlp.StoreAuthKey(c, user.Id, user.Auth)
		// set 'auth' cookie
		authhlp.SetAuthCookie(w, user.Auth)
	}
	http.Redirect(w, r, root, http.StatusFound)
}

// Verifies if token present in http.Response is valid.
// Valid means: data is valid, app_id match to server app_id and applicatiion name match
func isFacebookTokenValid(response *http.Response) (bool, error){

	tokenData, err := userhlp.FetchFacebookTokenData(response)
	if err == nil{
		if tokenData.Data.Is_valid &&
			(strconv.Itoa(tokenData.Data.App_id) == FACEBOOK_CLIENT_ID) &&
			(tokenData.Data.Application == "purple-wing"){
			return true, err
		}
	}
	return false, err
}

// Logout from session, clear authentication cookie and redirect to root.
func SessionLogout(w http.ResponseWriter, r *http.Request){
	authhlp.ClearAuthCookie(w)
	
	http.Redirect(w, r, root, http.StatusFound)
}

// json Google authentication handler
func JsonGoogleAuth(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	userInfo := userhlp.GPlusUserInfo{r.FormValue("id"), r.FormValue("email"), r.FormValue("name")}
	
	var err error
	var user *usermdl.User
	
	if !authhlp.CheckGoogleUserValidity(r.FormValue("access_token"), r) {
		return helpers.InternalServerError{errors.New("Access token is not valid")}
	}
	if !authhlp.IsAuthorizedWithGoogle(&userInfo) {
		return helpers.Forbidden{errors.New("You are not authorized to log in to purple-wing")}
	}
	if user, err = authhlp.SigninUser(w, r, "Email", userInfo.Email, userInfo.Name, userInfo.Name); err != nil{
		return helpers.InternalServerError{errors.New("Error occurred during signin process")}
	}

	// return user
	return templateshlp.RenderJson(w, c, user)
}
