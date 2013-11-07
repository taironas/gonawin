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
	"fmt"
	"net/http"
	"net/url"
	"html/template"
	"strconv"
	"time"
	
	"appengine"
	"appengine/urlfetch"
	
	oauth "github.com/garyburd/go-oauth/oauth"
	oauth2 "code.google.com/p/goauth2/oauth"
	
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	userhlp "github.com/santiaago/purple-wing/helpers/user"
	
	"github.com/santiaago/purple-wing/helpers/auth"
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

func Authenticate(w http.ResponseWriter, r *http.Request){
	if !auth.LoggedIn(r) {
		
		funcs := template.FuncMap{}
		
		t := template.Must(template.New("tmpl_auth").
			Funcs(funcs).
			ParseFiles("templates/session/auth.html"))
		// no data needed
		templateshlp.Render_with_data(w, r, t, nil, funcs, "renderAuth")
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}
// Google 

func AuthenticateWithGoogle(w http.ResponseWriter, r *http.Request){
	if !auth.LoggedIn(r) {
		url := googleConfig(r.Host).AuthCodeURL(r.URL.RawQuery)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}

func GoogleAuthCallback(w http.ResponseWriter, r *http.Request){
	// Exchange code for an access token at OAuth provider.
	code := r.FormValue("code")
	t := &oauth2.Transport{
		Config: googleConfig(r.Host),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(r),
		},
	}
	
	var userInfo *userhlp.GPlusUserInfo
	
	if _, err := t.Exchange(code); err == nil {
		userInfo, _ = userhlp.FetchGPlusUserInfo(r, t.Client())
	}
	if auth.IsAuthorizedWithGoogle(userInfo) {
		var user *usermdl.User
		// find user
		if user = usermdl.Find(r, "Email", userInfo.Email); user == nil {
			// create user if it does not exist
			user = usermdl.Create(r, userInfo.Email, userInfo.Name, userInfo.Name, auth.GenerateAuthKey())
		}
		// set 'auth' cookie
		auth.SetAuthCookie(w, user.Auth)
		// store in memcache auth key in memcaches
		auth.StoreAuthKey(r, user.Id, user.Auth)
	}

	http.Redirect(w, r, root, http.StatusFound)
}

// Twitter

func AuthenticateWithTwitter(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	if !auth.LoggedIn(r) {
		credentials, err := twitterConfig.RequestTemporaryCredentials(urlfetch.Client(c), "http://" + r.Host + twitterCallbackURL, nil)
		if err != nil {
			http.Error(w, "Error getting temp cred, "+err.Error(), 500)
			return
		}
		http.SetCookie(w, &http.Cookie{ Name: "secret", Value: credentials.Secret, Path: "/m", })
		http.Redirect(w, r, twitterConfig.AuthorizationURL(credentials, nil), 302)
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}

func TwitterAuthCallback(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	// get the request token
	requestToken := r.FormValue("oauth_token")
	// update credentials with request token and 'secret cookie value'
	var cred oauth.Credentials
	cred.Token = requestToken
	if cookie, err := r.Cookie("secret"); err == nil {
		cred.Secret = cookie.Value
	}
	// clear 'secret' cookie
	http.SetCookie(w, &http.Cookie{ Name: "secret", Path: "/m", Expires: time.Now(), })
	
	token, values, err := twitterConfig.RequestToken(urlfetch.Client(c), &cred, r.FormValue("oauth_verifier"))
	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}
	
	// get user info
	urlValues := url.Values{}
	urlValues.Set("user_id", values.Get("user_id"))
	resp, err := twitterConfig.Get(urlfetch.Client(c), token, "https://api.twitter.com/1.1/users/show.json", urlValues)
	if err != nil {
		c.Debugf("pw: error getting user info from twitter: %v", err)
	}

	userInfo, _ := userhlp.FetchTwitterUserInfo(resp)

	if auth.IsAuthorizedWithTwitter(userInfo) {
		var user *usermdl.User
		// find user
		if user = usermdl.Find(r, "Username", userInfo.Screen_name); user == nil {
			// create user if it does not exist
			user = usermdl.Create(r, "", userInfo.Screen_name, userInfo.Name, auth.GenerateAuthKey())
		}
		// set 'auth' cookie
		auth.SetAuthCookie(w, user.Auth)
		// store in memcache auth key in memcaches
		auth.StoreAuthKey(r, user.Id, user.Auth)
	}

	http.Redirect(w, r, root, http.StatusFound)
}

// Facebook

func AuthenticateWithFacebook(w http.ResponseWriter, r *http.Request){
	if !auth.LoggedIn(r) {
		url := facebookConfig(r.Host).AuthCodeURL(r.URL.RawQuery)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}

func FacebookAuthCallback(w http.ResponseWriter, r *http.Request){
	
	c := appengine.NewContext(r)
	c.Infof("pw: FacebookAuthCallback: Enter!")
	// Exchange code for an access token at OAuth provider.
	code := r.FormValue("code")
	t := &oauth2.Transport{
		Config: facebookConfig(r.Host),
		Transport: &urlfetch.Transport{
			Context: appengine.NewContext(r),
		},
	}	
	if v, err := t.Exchange(code); err != nil {
		c.Errorf("pw: Error in Exchange: %v, returned value: %v", err, v)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}

	c.Infof("pw: FacebookAuthCallback Exchange ok!")
	// build facebook url from debug_token url and get response
	urlFacebook := fmt.Sprintf(kUrlFacebookDebugToken, t.Token.AccessToken, FACEBOOK_CLIENT_ID, FACEBOOK_CLIENT_SECRET)
	accessTokenResponse, err := t.Client().Get(urlFacebook)
	if err != nil {
		c.Errorf("pw: Client Get error calling %s : %v",urlFacebook, err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	// verify if information from facebook is valid for current user.
	if isValid, err := isFacebookTokenValid(accessTokenResponse); (err != nil || !isValid){
		c.Errorf("pw: isFacebookTokenValid: Is valid: %v, Error: %v", isValid, err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	// ask facebook.com/me for user profile information
	graphResponse, err := t.Client().Get(kUrlFacebookMe)
	if err != nil {
		c.Errorf("pw: Failure on Get request %s: %v", kUrlFacebookMe, err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	userInfo, err := userhlp.FetchFacebookUserInfo(graphResponse)
	if err != nil{
		c.Errorf("pw: FetchFacebookUserInfo: %v", err)
		http.Redirect(w, r, root, http.StatusFound)
		return
	}
	if auth.IsAuthorizedWithFacebook(userInfo){
		var user *usermdl.User
		// find user
		if user = usermdl.Find(r, "Email", userInfo.Email); user == nil {
			// create user if it does not exist
			user = usermdl.Create(r, userInfo.Email, userInfo.Name, userInfo.Name, auth.GenerateAuthKey())
		}
		// set 'auth' cookie
		auth.SetAuthCookie(w, user.Auth)
		// store in memcache auth key in memcaches
		auth.StoreAuthKey(r, user.Id, user.Auth)
	}
	http.Redirect(w, r, root, http.StatusFound)
}

func isFacebookTokenValid(response *http.Response) (bool, error){

	tokeData, err := userhlp.FetchFacebookTokenData(response)
	if err == nil{
		if tokeData.Data.Is_valid && 
			(strconv.Itoa(tokeData.Data.App_id) == FACEBOOK_CLIENT_ID) &&
			(tokeData.Data.Application == "purple-wing"){
			return true, err
		}
	}
	return false, err
}

func SessionLogout(w http.ResponseWriter, r *http.Request){
	auth.ClearAuthCookie(w)
	
	http.Redirect(w, r, root, http.StatusFound)
}
