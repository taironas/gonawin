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
	"errors"
	"net/http"
	"strconv"
	"strings"

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/auth"
	"github.com/santiaago/purple-wing/helpers/log"

	usermdl "github.com/santiaago/purple-wing/models/user"
)

// parse permalink id from URL  and return it
func PermalinkID(r *http.Request, c appengine.Context, level int64) (int64, error) {

	path := strings.Split(r.URL.String(), "/")
	// if url has params extract id until the ? character
	var strID string

	if strings.Contains(r.URL.String(), "?") {
		strPath := path[level]
		strID = strPath[0:strings.Index(strPath, "?")]
	} else {
		strID = path[level]
	}
	intID, err := strconv.ParseInt(strID, 0, 64)
	if err != nil {
		log.Errorf(c, " error when calling PermalinkID with %v.Error: %v", path[level], err)
	}
	return intID, err
}

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request) error

// Error handler returns the proper error handler function with respecto to the error rised by the function called.
func ErrorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case helpers.BadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case helpers.NotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case helpers.Forbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case helpers.InternalServerError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			c := appengine.NewContext(r)
			log.Errorf(c, "%v", err)
			http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
		}
	}
}

// Authorized runs the function pass by parameter and checks authentication data prior to any call. Will rise a bad request error hanlder if authentication fails.
func Authorized(f func(w http.ResponseWriter, r *http.Request, u *usermdl.User) error) ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if auth.KOfflineMode {
			return f(w, r, auth.CurrentOfflineUser(r, appengine.NewContext(r)))
		}

		user := auth.CheckAuthenticationData(r)
		if user == nil {
			return helpers.BadRequest{errors.New("Bad Authentication data")}
		} else {
			return f(w, r, user)
		}
	}
}
