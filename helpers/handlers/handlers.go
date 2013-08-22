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
 
package handlers

import (
	"net/http"
	"strings"
	"strconv"

	"appengine"	
	
	"github.com/santiaago/purple-wing/helpers/auth"
)
// is it a user?
func User(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsUser(r) {
			http.Redirect(w, r, "/m", http.StatusFound)
		} else{
			f(w, r)
		}
	}
}
// is it an admin?
func Admin(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAdmin(r) {
			http.Redirect(w, r, "/m", http.StatusFound)
		} else{
			f(w, r)
		}
	}
}
// parse permalink id from URL  and return it
func PermalinkID(r *http.Request, level int64)(int64, error){

	c := appengine.NewContext(r)
	path := strings.Split(r.URL.String(), "/")
	intID, err := strconv.ParseInt(path[level],0,64)
	if err != nil{
		c.Errorf("pw: error when calling PermalinkID with %v.Error: %v",path[level], err)
	}
	return intID, err
}
