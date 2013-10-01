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

package tournamentrels

import (
	"net/http"
	"strconv"
	
	"appengine"
	
	"github.com/santiaago/purple-wing/helpers/auth"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
)

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get tournament id
	tournamentId , _ := strconv.ParseInt(r.FormValue("TournamentId"), 10, 64)

	if r.Method == "POST" && r.FormValue("Action") == "post_action" {
		if err := tournamentmdl.Join(r, tournamentId, auth.CurrentUser(r).Id); err != nil {
			c.Errorf("pw: tournamentrels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action") == "delete_action" {
		if err := tournamentmdl.Leave(r, tournamentId, auth.CurrentUser(r).Id); err != nil {
			c.Errorf("pw: tournamentrels.Show: %v", err)
		}
	}
	
	http.Redirect(w,r, "/m/tournaments/"+r.FormValue("TournamentId"), http.StatusFound)
}

