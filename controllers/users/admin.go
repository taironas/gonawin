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

package users

import (
	"html/template"
	"net/http"
	"time"

	"appengine"

	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	usermdl "github.com/santiaago/purple-wing/models/user"
)

// Admin show handler
func AdminShow(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	funcs := template.FuncMap{}

	t := template.Must(template.New("tmpl_admin_show").
		Funcs(funcs).
		ParseFiles("templates/admin/show.html"))

	templateshlp.RenderWithData(w, r, c, t, nil, funcs, "renderAdminShow")
}

// Admin users handler
func AdminUsers(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// sample of users
	user1 := usermdl.User{1, "test1@example.com", "jdoe1", "John Doe 1", false, "", time.Now()}
	user2 := usermdl.User{1, "test2@example.com", "jdoe2", "John Doe 2", false, "", time.Now()}
	user3 := usermdl.User{1, "test3@example.com", "jdoe3", "John Doe 3", false, "", time.Now()}
	users := []usermdl.User{user1, user2, user3}
	// end samlpe of users

	funcs := template.FuncMap{}

	t := template.Must(template.New("tmpl_admin_users_show").
		Funcs(funcs).
		ParseFiles("templates/admin/users.html",
		"templates/user/info.html"))

	templateshlp.RenderWithData(w, r, c, t, users, funcs, "renderAdminUsersShow")
}
