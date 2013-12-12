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

	"github.com/santiaago/purple-wing/helpers"	
	"github.com/santiaago/purple-wing/helpers/log"	
	"github.com/santiaago/purple-wing/helpers/auth"
	tournamentmdl "github.com/santiaago/purple-wing/models/tournament"
)

// show handler for tournament relationships
func Show(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	
	// get tournament id
	tournamentId , err := strconv.ParseInt(r.FormValue("TournamentId"), 10, 64)
	if err != nil {
		log.Errorf(c, " tournaments.Show, string value could not be parsed: %v", err)
	}

	if r.Method == "POST" && r.FormValue("Action") == "post_action" {
		if err := tournamentmdl.Join(c, tournamentId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " tournamentrels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action") == "delete_action" {
		if err := tournamentmdl.Leave(c, tournamentId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " tournamentrels.Show: %v", err)
		}
	}
	
	http.Redirect(w,r, "/m/tournaments/"+r.FormValue("TournamentId"), http.StatusFound)
}

// json show handler for tournament relationships
func ShowJson(w http.ResponseWriter, r *http.Request) error{
	c := appengine.NewContext(r)
	
	// get tournament id
	tournamentId , err := strconv.ParseInt(r.FormValue("TournamentId"), 10, 64)
	if err != nil {
		log.Errorf(c, " tournaments.Show, string value could not be parsed: %v", err)
		return helpers.NotFound{err}
	}

	if r.Method == "POST" && r.FormValue("Action") == "post_action" {
		if err := tournamentmdl.Join(c, tournamentId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " tournamentrels.Show: %v", err)
		}
	} else if r.Method == "POST" && r.FormValue("Action") == "delete_action" {
		if err := tournamentmdl.Leave(c, tournamentId, auth.CurrentUser(r, c).Id); err != nil {
			log.Errorf(c, " tournamentrels.Show: %v", err)
		}
	}
	
	http.Redirect(w,r, "/j/tournaments/"+r.FormValue("TournamentId"), http.StatusFound)
	return nil
}
