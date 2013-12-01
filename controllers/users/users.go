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

	"appengine"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/handlers"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	teamrelshlp "github.com/santiaago/purple-wing/helpers/teamrels"
	tournamentrelshlp "github.com/santiaago/purple-wing/helpers/tournamentrels"
	
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	teamrequestmdl "github.com/santiaago/purple-wing/models/teamrequest"
	
)

type Form struct {
	Username string
	Name string
	Email string
	ErrorUsername string
	ErrorName string
	ErrorEmail string
}

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	userId, err := handlers.PermalinkID(r, c, 3)
	if err != nil{
		http.Redirect(w,r, "/m/users/", http.StatusFound)
		return
	}
	
	funcs := template.FuncMap{
		"Profile": func() bool {return true},
	}
	
	t := template.Must(template.New("tmpl_user_show").
		Funcs(funcs).
		ParseFiles("templates/user/show.html", 
		"templates/user/info.html",
		"templates/user/teams.html",
		"templates/user/tournaments.html",
		"templates/user/requests.html"))

	var user *usermdl.User
	user, err = usermdl.ById(c,userId)
	if err != nil{
		helpers.Error404(w)
		return
	}
	
	teams := teamrelshlp.Teams(c, userId)
	tournaments := tournamentrelshlp.Tournaments(c, userId)
	teamRequests := teamrelshlp.TeamsRequests(c, teams)
	
	userData := struct { 
		User *usermdl.User
		Teams []*teammdl.Team
		Tournaments []*tournamentmdl.Tournament
		TeamRequests []*teamrequestmdl.TeamRequest
	}{
		user,
		teams,
		tournaments,
		teamRequests,
	}
	
	templateshlp.RenderWithData(w, r, c, t, userData, funcs, "renderUserShow")
}
