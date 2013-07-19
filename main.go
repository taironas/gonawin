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

	"github.com/santiaago/purple-wing/controllers"
	"github.com/santiaago/purple-wing/helpers"
)

func init(){
	h := new(helpers.RegexpHandler)
	/* usual pages*/
	h.HandleFunc("/", controllers.Home)
	//h.HandleFunc("/", controllers.About)
	//h.HandleFunc("/", controllers.Contact)
	/* session */
	h.HandleFunc("/auth/?", controllers.SessionAuth)
	h.HandleFunc("/oauth2callback/?", controllers.SessionAuthCallback)
	h.HandleFunc("/logout/?", controllers.SessionLogout)	
	/* user */
	h.HandleFunc("/users/[0-9]+/?", controllers.UserShow)
	h.HandleFunc("/users/[0-9]+/edit/?", controllers.UserEdit)

	http.Handle("/", h)
}
