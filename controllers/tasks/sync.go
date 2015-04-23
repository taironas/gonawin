/*
 * Copyright (c) 2014 Santiago Arias | Remy Jourde
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

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"

	mdl "github.com/santiaago/gonawin/models"
)

// SyncScores handler, use it to synchronize all the scores.
//
//	GET	/a/sync/scores/", Description..
//
// Go though all particiants of tournament passed by HTTP POST request.
// For each user, compute again scores of each finished match and update Score entity and global score of user.
func SyncScores(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Sync Scores Handler:"

	if r.Method != "POST" {
		log.Infof(c, "%s something went wrong...", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	log.Infof(c, "%s processing...", desc)

	tournamentBlob := []byte(r.FormValue("tournament"))

	var t mdl.Tournament
	err1 := json.Unmarshal(tournamentBlob, &t)
	if err1 != nil {
		log.Errorf(c, "%s unable to extract tournament from data, %v.", desc, err1)
		return err1
	}

	log.Infof(c, "%s value of tournament id: %v.", desc, t.Id)
	log.Infof(c, "%s get tournament participants.", desc)

	users := t.Participants(c)

	// prepare data.
	log.Infof(c, "%s preparing data...", desc)
	log.Infof(c, "%s go through each participant and compute global scores.", desc)
	for _, u := range users {
		// update global score of user
		globalScore := int64(0)
		for _, tid := range u.TournamentIds {
			score := u.ScoreByTournament(c, tid)
			globalScore = globalScore + score
		}
		u.Score = globalScore
		if err := u.Update(c); err != nil {
			log.Errorf(c, "%s unable to update user %v with global score. %v", desc, u.Id, err)
			continue
		}
	}
	log.Infof(c, "%s task done!", desc)
	return nil
}
