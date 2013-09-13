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
	usermdl "github.com/santiaago/purple-wing/models/user"
)

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get team id
	teamId , _ := strconv.ParseInt(r.FormValue("TeamId"), 10, 64)
	c.Infof("pw: TeamRel.Show, r.Method = %s", r.Method)	
	if r.Method == "POST" {
		if err := usermdl.Join(r, teamId, auth.CurrentUser(r).Id); err != nil {
			c.Errorf("pw: teamRels.Show: %v", err)
		}
	} else if r.Method == "DELETE" {
		if err := usermdl.Leave(r, teamId, auth.CurrentUser(r).Id); err != nil {
			c.Errorf("pw: teamRels.Show: %v", err)
		}
	}
	
	http.Redirect(w,r, "/m/teams/"+r.FormValue("TeamId"), http.StatusFound)
}

