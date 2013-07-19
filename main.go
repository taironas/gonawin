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

package pw

import (
	"net/http"
	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/controllers"
)

func init(){
	h := new(helpers.RegexpHandler)

	h.HandleFunc("/", controllers.Home)
	/* session */
	h.HandleFunc("/auth/?", controllers.Auth)
	h.HandleFunc("/oauth2callback/?", controllers.AuthCallback)
	h.HandleFunc("/logout/?", controllers.Logout)	
	/* user */
	h.HandleFunc("/users/[0-9]+/?", controllers.Show)
	h.HandleFunc("/users/[0-9]+/edit/?", controllers.Edit)

	http.Handle("/", h)
}
