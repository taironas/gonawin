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

package tournamentteamrels

import (
	"net/http"
	"strconv"
	
	"appengine"
	
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
	"github.com/santiaago/purple-wing/helpers/log"
)

func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get tournament id
	tournamentId , err := strconv.ParseInt(r.FormValue("TournamentId"), 10, 64)
	if err != nil {
		log.Errorf(c, " tournamentteamrels.Show, string value could not be parsed: %v", err)
	}
	// get team id
	teamId , err := strconv.ParseInt(r.FormValue("TeamIdButton"), 10, 64)
	if err != nil {
		log.Errorf(c, " tournamentteamrels.Show, string value could not be parsed: %v", err)
	}

	if r.Method == "POST" && r.FormValue("Action_" + r.FormValue("TeamIdButton")) == "post_action" {
		if err := tournamentmdl.TeamJoin(r, tournamentId, teamId); err != nil {
			log.Errorf(c, " tournamentteamrels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action_" + r.FormValue("TeamIdButton")) == "delete_action" {
		if err := tournamentmdl.TeamLeave(r, tournamentId, teamId); err != nil {
			log.Errorf(c, " tournamentteamrels.Show: %v", err)
		}
	}
	
	http.Redirect(w,r, "/m/tournaments/"+r.FormValue("TournamentId"), http.StatusFound)
}

