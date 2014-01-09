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

package teamrels

import (
	"net/http"
	"strconv"
	
	"appengine"

	"github.com/santiaago/purple-wing/helpers"	
	"github.com/santiaago/purple-wing/helpers/log"	
	"github.com/santiaago/purple-wing/helpers/auth"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
	usermdl "github.com/santiaago/purple-wing/models/user"
	teammdl "github.com/santiaago/purple-wing/models/team"
)

// create handler for team relations
func Create(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	if err != nil {
		log.Errorf(c, " teamRels.Create, string value could not be parsed: %v", err)
	}
	
	if r.Method == "POST" {
		if err := teammdl.Join(c, teamId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " teamRels.Create: %v", err)
		}
	}
	
	http.Redirect(w,r, "/m/teams/"+r.FormValue("TeamId"), http.StatusFound)
}

// destroy handler for team relations
func Destroy(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	if err != nil {
		log.Errorf(c, " teamRels.Destroy, string value could not be parsed: %v", err)
	}
	
	if r.Method == "POST" {
		if !teammdl.IsTeamAdmin(c, teamId, auth.CurrentUser(r, c).Id) {
			if err := teammdl.Leave(c, teamId, auth.CurrentUser(r, c).Id); err != nil {
				log.Errorf(c, " teamRels.Destroy: %v", err)
			}
		} else {
			log.Errorf(c, " teamRels.Drestroy, Team administrator cannot leave the team")
		}
	}
	
	http.Redirect(w,r, "/m/teams/"+r.FormValue("TeamId"), http.StatusFound)
}

// json create handler for team relations
func CreateJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	if err != nil {
		log.Errorf(c, " teamRels.Create, string value could not be parsed: %v", err)
		return helpers.NotFound{err}
	}
	
	if r.Method == "POST" {
		if err := teammdl.Join(c, teamId, u.Id); err != nil {
			log.Errorf(c, " teamRels.Create: %v", err)
			return helpers.InternalServerError{err}
		}
	}
	
	// return the joined team
	var team *teammdl.Team
	if team, err = teammdl.ById(c, teamId); err != nil{
		return helpers.NotFound{err}
	}
	return templateshlp.RenderJson(w, c, team)
}

// json destroy handler for team relations
func DestroyJson(w http.ResponseWriter, r *http.Request, u *usermdl.User) error{
	c := appengine.NewContext(r)
	
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	if err != nil {
		log.Errorf(c, " teamRels.Destroy, string value could not be parsed: %v", err)
		return helpers.NotFound{err}
	}
	
	if r.Method == "POST" {
		if !teammdl.IsTeamAdmin(c, teamId, u.Id) {
			if err := teammdl.Leave(c, teamId, u.Id); err != nil {
				log.Errorf(c, " teamRels.Destroy: %v", err)
				return helpers.InternalServerError{err}
			}
		} else {
			log.Errorf(c, " teamRels.Destroy, Team administrator cannot leave the team")
			return helpers.BadRequest{err}
		}
	}
	
	// return the left team
	var team *teammdl.Team
	if team, err = teammdl.ById(c, teamId); err != nil{
		return helpers.NotFound{err}
	}
	return templateshlp.RenderJson(w, c, team)
}