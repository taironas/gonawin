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

// Package handlers provides a set of functions to manipulate http.HandlerFunc
// at a route level. cf github.com/taironas/gonawin/gonawin/main.go for details
// on how to use it.
//
package handlers

import (
	"errors"
	"net/http"

	"appengine"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/auth"
	"github.com/taironas/gonawin/helpers/log"

	mdl "github.com/taironas/gonawin/models"
)

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request) error

// Error handler returns the proper error handler function with respecto to the error rised by the function called.
func ErrorHandler(f func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err == nil {
			return
		}
		switch err.(type) {
		case *helpers.BadRequest:
			http.Error(w, err.Error(), http.StatusBadRequest)
		case *helpers.NotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case *helpers.Forbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
		case *helpers.Unauthorized:
			http.Error(w, err.Error(), http.StatusUnauthorized)
		case *helpers.InternalServerError:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		default:
			c := appengine.NewContext(r)
			log.Errorf(c, "%v", err)
			http.Error(w, "Sorry, something went wrong.", http.StatusInternalServerError)
		}
	}
}

// Authorized runs the function pass by parameter and checks authentication data prior to any call.
// Will rise a bad request error handler if authentication fails.
func Authorized(f func(w http.ResponseWriter, r *http.Request, u *mdl.User) error) ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var user *mdl.User
		if auth.KOfflineMode {
			user = auth.CurrentOfflineUser(r, appengine.NewContext(r))
		} else {
			user = auth.CheckAuthenticationData(r)
		}

		if user == nil {
			return &helpers.BadRequest{Err: errors.New("Bad Authentication data")}
		} else {
			return f(w, r, user)
		}
	}
}

// Admin Authorized runs the function pass by parameter and checks authentication data prior to any call.
// Will rise a bad request error handler if authentication fails. User should be a gonawin admin .
func AdminAuthorized(f func(w http.ResponseWriter, r *http.Request, u *mdl.User) error) ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var user *mdl.User
		if auth.KOfflineMode {
			user = auth.CurrentOfflineUser(r, appengine.NewContext(r))
		} else {
			user = auth.CheckAuthenticationData(r)
			if user == nil {
				return &helpers.BadRequest{Err: errors.New("Bad Authentication data")}
			}
		}
		if !auth.IsGonawinAdmin(appengine.NewContext(r)) { //user) {
			return &helpers.Forbidden{Err: errors.New(helpers.ErrorCodeSessionsForbiden)}
		}
		return f(w, r, user)
	}
}
