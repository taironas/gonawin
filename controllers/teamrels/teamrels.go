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
	
	"github.com/santiaago/purple-wing/helpers/auth"
	teammdl "github.com/santiaago/purple-wing/models/team"
)

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get team id
	teamId , _ := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	
	if r.Method == "POST" && r.FormValue("Action") == "post_action" {
		if err := teammdl.Join(r, teamId, auth.CurrentUser(r).Id); err != nil {
			c.Errorf("pw: teamRels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action") == "delete_action" {
		if !teammdl.IsTeamAdmin(r, teamId, auth.CurrentUser(r).Id) {
			if err := teammdl.Leave(r, teamId, auth.CurrentUser(r).Id); err != nil {
				c.Errorf("pw: teamRels.Show: %v", err)
			}
		} else {
			c.Errorf("pw: teamRels.Show: Team administrator cannot leave the team")
		}
	}
	
	http.Redirect(w,r, "/m/teams/"+r.FormValue("TeamId"), http.StatusFound)
}

