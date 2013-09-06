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
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"html/template"
	"time"
	
	"appengine"
	"appengine/urlfetch"
	
	oauth "github.com/garyburd/go-oauth/oauth"
	oauth2 "code.google.com/p/goauth2/oauth"
	
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	"github.com/santiaago/purple-wing/helpers/auth"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

const root string = "/m"
// Set up a configuration for google.
func googleConfig(host string) *oauth2.Config{
	return &oauth2.Config{
		ClientId:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
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

func Authenticate(w http.ResponseWriter, r *http.Request){
	if !auth.LoggedIn(r) {
		c := appengine.NewContext(r)
		
		funcs := template.FuncMap{}
		
		t := template.Must(template.New("tmpl_auth").
			Funcs(funcs).
			ParseFiles("templates/session/auth.html"))
		
		var buf bytes.Buffer
		err := t.ExecuteTemplate(&buf,"tmpl_auth", nil)
		main := buf.Bytes()
		
		if err != nil{
			c.Errorf("pw: error executing template auth: %v", err)
		}
		err = templateshlp.Render(w, r, main, &funcs, "renderAuth")
		
		if err != nil{
			c.Errorf("pw: error when calling Render from helpers in Authenticate Handler: %v", err)
		}
	} else {
		//redirect to home page
		http.Redirect(w, r, root, http.StatusFound)
	}
}

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
	
	var userInfo *usermdl.GPlusUserInfo
	
	if _, err := t.Exchange(code); err == nil {
		userInfo, _ = usermdl.FetchGPlusUserInfo(r, t.Client())
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

	userInfo, _ := usermdl.FetchTwitterUserInfo(resp)

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

func SessionLogout(w http.ResponseWriter, r *http.Request){
	auth.ClearAuthCookie(w)
	
	http.Redirect(w, r, root, http.StatusFound)
}
