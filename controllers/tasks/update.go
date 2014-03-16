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

package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"appengine"
	//"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"

	mdl "github.com/santiaago/purple-wing/models"
)

// update score  handler:
//
// Use this handler to ...
//	GET	/a/update/scores/	Description..
//
// The response is ...
func UpdateScores(w http.ResponseWriter, r *http.Request /*, u *mdl.User*/) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Update Scores Handler:"
	log.Infof(c, "%s task called, processing...", desc)
	if r.Method == "POST" {
		// all or nothing here :)
		tournamentBlob := []byte(r.FormValue("tournament"))
		matchBlob := []byte(r.FormValue("match"))

		var t mdl.Tournament
		err1 := json.Unmarshal(tournamentBlob, &t)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract tournament from data, %v", desc, err1)
		}

		var m mdl.Tmatch
		err2 := json.Unmarshal(matchBlob, &m)
		if err2 != nil {
			log.Errorf(c, "%s unable to extract match from data, %v", desc, err2)
		}

		log.Infof(c, "%s value of tournament id: %v", desc, t.Id)
		log.Infof(c, "%s value of match id: %v", desc, m.Id)

		users := t.Participants(c)
		usersToUpdate := make([]*mdl.User, 0)
		for i, u := range users {
			if score, err := u.ScoreForMatch(c, &m); err != nil {
				log.Errorf(c, "%s unable udpate user %v score: %v", desc, u.Id, err)
			} else {
				// update user overall score
				users[i].Score += score
				usersToUpdate = append(usersToUpdate, users[i])
				// update score entity for user, tournament pair.
				// if does not exist, create it and update it
				// else update it
				if scoreEntity, _ := u.TournamentScore(c, &t); scoreEntity == nil {
					log.Infof(c, "%s create score entity as it does not exist", desc)
					if scoreEntity1, err := mdl.CreateScore(c, u.Id, t.Id); err != nil {
						log.Errorf(c, "%s unable to create score entity", desc)
						return err
					} else {
						log.Infof(c, "%s score ready add it to tournament %v", desc, scoreEntity1)
						u.AddTournamentScore(c, scoreEntity1.Id, t.Id)
						log.Infof(c, "%s score entity exists now, lets update it", desc)
						var err error
						if err = scoreEntity1.Add(c, score); err != nil {
							log.Errorf(c, "%s unable to add score of user %v, ", desc, u.Id, err)
						}
					}
				} else {
					log.Infof(c, "%s score entity exists, lets update it", desc)
					var err error
					if err = scoreEntity.Add(c, score); err != nil {
						log.Errorf(c, "%s unable to add score of user %v, ", desc, u.Id, err)
					}
				}
			}
		}

		if err := mdl.UpdateUsers(c, usersToUpdate); err != nil {
			log.Errorf(c, "%s unable udpate users scores: %v", desc, err)
			return errors.New(helpers.ErrorCodeUsersCannotUpdate)
		}
		log.Infof(c, "%s task done!", desc)
		return nil
	}
	log.Infof(c, "%s something went wrong...")
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
