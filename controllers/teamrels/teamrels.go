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
	teammdl "github.com/santiaago/purple-wing/models/team"
)

// show handler for team relations
func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	if err != nil {
		log.Errorf(c, " teamRels.Show, string value could not be parsed: %v", err)
	}
	
	if r.Method == "POST" && r.FormValue("Action") == "post_action" {
		if err := teammdl.Join(c, teamId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " teamRels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action") == "delete_action" {
		if !teammdl.IsTeamAdmin(c, teamId, auth.CurrentUser(r, c).Id) {
			if err := teammdl.Leave(c, teamId, auth.CurrentUser(r, c).Id); err != nil {
				log.Errorf(c, " teamRels.Show: %v", err)
			}
		} else {
			log.Errorf(c, " teamRels.Show, Team administrator cannot leave the team")
		}
	}
	
	http.Redirect(w,r, "/m/teams/"+r.FormValue("TeamId"), http.StatusFound)
}

// json show handler for team relations
func ShowJson(w http.ResponseWriter, r *http.Request) error{
	c := appengine.NewContext(r)
	
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	if err != nil {
		log.Errorf(c, " teamRels.Show, string value could not be parsed: %v", err)
		return helpers.NotFound{err}
	}
	
	if r.Method == "POST" && r.FormValue("Action") == "post_action" {
		if err := teammdl.Join(c, teamId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " teamRels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action") == "delete_action" {
		if !teammdl.IsTeamAdmin(c, teamId, auth.CurrentUser(r, c).Id) {
			if err := teammdl.Leave(c, teamId, auth.CurrentUser(r, c).Id); err != nil {
				log.Errorf(c, " teamRels.Show: %v", err)
			}
		} else {
			log.Errorf(c, " teamRels.Show, Team administrator cannot leave the team")
		}
	}
	
	http.Redirect(w,r, "/j/teams/"+r.FormValue("TeamId"), http.StatusFound)
	return nil
}

