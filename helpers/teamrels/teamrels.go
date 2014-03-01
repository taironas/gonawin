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
	"appengine"

	"github.com/santiaago/purple-wing/helpers/log"

	teammdl "github.com/santiaago/purple-wing/models/team"
	teamrelmdl "github.com/santiaago/purple-wing/models/teamrel"
	teamrequestmdl "github.com/santiaago/purple-wing/models/teamrequest"
	//mdl "github.com/santiaago/purple-wing/models/user"
	mdl "github.com/santiaago/purple-wing/models"
)

// from a team id return an array of users/ players that participates in it.
func Players(c appengine.Context, teamId int64) []*mdl.User {

	var users []*mdl.User

	teamRels := teamrelmdl.Find(c, "TeamId", teamId)

	for _, teamRel := range teamRels {
		user, err := mdl.UserById(c, teamRel.UserId)
		if err != nil {
			log.Errorf(c, " Players, cannot find user with ID=%", teamRel.UserId)
		} else {
			users = append(users, user)
		}
	}

	return users
}

// build a teamRequest array from an array of teams
func TeamsRequests(c appengine.Context, teams []*teammdl.Team) []*teamrequestmdl.TeamRequest {
	var teamRequests []*teamrequestmdl.TeamRequest

	for _, team := range teams {
		teamRequests = append(teamRequests, teamrequestmdl.Find(c, "TeamId", team.Id)...)
	}

	return teamRequests
}
