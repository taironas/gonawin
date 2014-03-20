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
	"net/url"

	"appengine"
	//"appengine/datastore"
	"appengine/taskqueue"

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

		// prepare data.
		log.Infof(c, "%s preparing data...", desc)

		scores := make([]int64, 0)
		userIds := make([]int64, 0)
		userIdsToCreateSE := make([]int64, 0)
		tournamentId := t.Id

		for _, u := range users {
			if score, err := u.ScoreForMatch(c, &m); err != nil {
				log.Errorf(c, "%s unable udpate user %v score: %v", desc, u.Id, err)
			} else {
				scores = append(scores, score)
				userIds = append(userIds, u.Id)
			}
			if scoreEntity, _ := u.TournamentScore(c, &t); scoreEntity == nil {
				userIdsToCreateSE = append(userIdsToCreateSE, u.Id)
			}
		}
		log.Infof(c, "%s the data is ready.", desc)

		// task queue for updating scores of users.
		log.Infof(c, "%s task queue for updating scores of users: -->", desc)

		bscores, errm11 := json.Marshal(scores)
		if errm11 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm11)
		}
		buserIds, errm12 := json.Marshal(userIds)
		if errm12 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm12)
		}
		btournamentId, errm13 := json.Marshal(tournamentId)
		if errm13 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm13)
		}
		task1 := taskqueue.NewPOSTTask("/a/update/users/scores/", url.Values{
			"userIds":      []string{string(buserIds)},
			"scores":       []string{string(bscores)},
			"tournamentId": []string{string(btournamentId)},
		})
		if _, err := taskqueue.Add(c, task1, ""); err != nil {
			log.Errorf(c, "%s unable to add task to taskqueue.", desc)
			return err
		} else {
			log.Infof(c, "%s add task to taskqueue successfully", desc)
		}
		log.Infof(c, "%s task queue for updating scores of users: <--", desc)

		// task queue for adding necessary score entities.
		log.Infof(c, "%s task queue for adding necessary score entities.: -->", desc)

		buserIdsToCreateSE, errm21 := json.Marshal(userIdsToCreateSE)
		if errm21 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm21)
		}
		bscores, errm22 := json.Marshal(scores)
		if errm22 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm22)
		}

		task2 := taskqueue.NewPOSTTask("/a/create/scoreentities/", url.Values{
			"userIds":      []string{string(buserIdsToCreateSE)},
			"scores":       []string{string(bscores)},
			"tournamentId": []string{string(btournamentId)},
		})
		if _, err := taskqueue.Add(c, task2, ""); err != nil {
			log.Errorf(c, "%s unable to add task to taskqueue.", desc)
			return err
		} else {
			log.Infof(c, "%s add task to taskqueue successfully", desc)
		}
		log.Infof(c, "%s task queue for adding necessary score entities.: <--", desc)

		// task queue for adding the score to the score entity.
		log.Infof(c, "%s task queue for adding the score to the score entity: -->", desc)

		bscores, errm31 := json.Marshal(scores)
		if errm31 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm31)
		}
		buserIds, errm2 := json.Marshal(userIds)
		if errm2 != nil {
			log.Errorf(c, "%s Error marshaling", desc, errm2)
		}
		task := taskqueue.NewPOSTTask("/a/add/scoreentities/score/", url.Values{
			"userIds":    []string{string(buserIds)},
			"scores":     []string{string(bscores)},
			"tournament": []string{string(tournamentBlob)},
		})
		if _, err := taskqueue.Add(c, task, ""); err != nil {
			log.Errorf(c, "%s unable to add task to taskqueue.", desc)
			return err
		} else {
			log.Infof(c, "%s add task to taskqueue successfully", desc)
		}
		log.Infof(c, "%s task queue for adding the score to the score entity: <--", desc)

		// users := t.Participants(c)
		// usersToUpdate := make([]*mdl.User, 0)
		// for i, u := range users {
		// 	if score, err := u.ScoreForMatch(c, &m); err != nil {
		// 		log.Errorf(c, "%s unable udpate user %v score: %v", desc, u.Id, err)
		// 	} else {
		// 		// update user overall score
		// 		users[i].Score += score
		// 		usersToUpdate = append(usersToUpdate, users[i])
		// 		// update score entity for user, tournament pair.
		// 		// if does not exist, create it and update it
		// 		// else update it
		// 		if scoreEntity, _ := u.TournamentScore(c, &t); scoreEntity == nil {
		// 			log.Infof(c, "%s create score entity as it does not exist", desc)
		// 			if scoreEntity1, err := mdl.CreateScore(c, u.Id, t.Id); err != nil {
		// 				log.Errorf(c, "%s unable to create score entity", desc)
		// 				return err
		// 			} else {
		// 				log.Infof(c, "%s score ready add it to tournament %v", desc, scoreEntity1)
		// 				u.AddTournamentScore(c, scoreEntity1.Id, t.Id)
		// 				log.Infof(c, "%s score entity exists now, lets update it", desc)
		// 				var err error
		// 				if err = scoreEntity1.Add(c, score); err != nil {
		// 					log.Errorf(c, "%s unable to add score of user %v, ", desc, u.Id, err)
		// 				}
		// 			}
		// 		} else {
		// 			log.Infof(c, "%s score entity exists, lets update it", desc)
		// 			var err error
		// 			if err = scoreEntity.Add(c, score); err != nil {
		// 				log.Errorf(c, "%s unable to add score of user %v, ", desc, u.Id, err)
		// 			}
		// 		}
		// 	}
		// }

		// if err := mdl.UpdateUsers(c, usersToUpdate); err != nil {
		// 	log.Errorf(c, "%s unable udpate users scores: %v", desc, err)
		// 	return errors.New(helpers.ErrorCodeUsersCannotUpdate)
		// }
		log.Infof(c, "%s task done!", desc)
		return nil
	}
	log.Infof(c, "%s something went wrong...")
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Update users scores.
func UpdateUsersScores(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Update Users Scores Handler:"
	log.Infof(c, "%s task called, processing...", desc)
	if r.Method == "POST" {
		log.Infof(c, "%s reading data...", desc)

		userIdsBlob := []byte(r.FormValue("userIds"))
		scoresBlob := []byte(r.FormValue("scores"))

		var userIds []int64
		err1 := json.Unmarshal(userIdsBlob, &userIds)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract userIds from data, %v", desc, err1)
		}

		var scores []int64
		err2 := json.Unmarshal(scoresBlob, &scores)
		if err2 != nil {
			log.Errorf(c, "%s unable to extract scores from data, %v", desc, err2)
		}

		log.Infof(c, "%s value of user ids: %v", desc, userIds)
		log.Infof(c, "%s value of scores: %v", desc, scores)

		log.Infof(c, "%s crunching data...", desc)
		usersToUpdate := make([]*mdl.User, 0)
		for i, id := range userIds {
			if u, err := mdl.UserById(c, id); err != nil {
				log.Errorf(c, "cannot find user with id=%", id)
			} else {
				u.Score += scores[i]
				usersToUpdate = append(usersToUpdate, u)
			}
		}
		log.Infof(c, "%s update users", desc)

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

// Create the score entities.
func CreateScoreEntities(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Create score entities Handler:"
	log.Infof(c, "%s task called, processing...", desc)
	if r.Method == "POST" {
		log.Infof(c, "%s preparing data", desc)

		userIdsBlob := []byte(r.FormValue("userIds"))
		tournamentIdBlob := []byte(r.FormValue("tournamentId"))

		var userIds []int64
		err1 := json.Unmarshal(userIdsBlob, &userIds)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract userIds from data, %v", desc, err1)
		}

		log.Infof(c, "%s value of user ids: %v", desc, userIds)

		var tournamentId int64
		err1 = json.Unmarshal(tournamentIdBlob, &tournamentId)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract tournamentId from data, %v", desc, err1)
		}
		log.Infof(c, "%s crunching data...", desc)
		for _, id := range userIds {
			if u, err := mdl.UserById(c, id); err != nil {
				log.Errorf(c, "%s cannot find user with id=%", desc, id)
			} else {
				log.Infof(c, "%s create score entity as it does not exist", desc)
				if se, err := mdl.CreateScore(c, u.Id, tournamentId); err != nil {
					log.Errorf(c, "%s unable to create score entity", desc)
					return err
				} else {
					log.Infof(c, "%s score ready add it to tournament %v", desc, se)
					u.AddTournamentScore(c, se.Id, se.TournamentId)
					log.Infof(c, "%s score entity exists now, lets update it", desc)
				}
			}
		}
		log.Infof(c, "%s task done!", desc)
		return nil
	}
	log.Infof(c, "%s something went wrong...")
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}

// Add score to score entities.
func AddScoreToScoreEntities(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Add score to score entity Handler:"
	log.Infof(c, "%s task called, processing...", desc)
	if r.Method == "POST" {
		log.Infof(c, "%s reading data...", desc)
		userIdsBlob := []byte(r.FormValue("userIds"))
		scoresBlob := []byte(r.FormValue("scores"))
		tournamentBlob := []byte(r.FormValue("tournament"))

		var userIds []int64
		err1 := json.Unmarshal(userIdsBlob, &userIds)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract userIds from data, %v", desc, err1)
		}

		var scores []int64
		err1 = json.Unmarshal(scoresBlob, &scores)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract userIds from data, %v", desc, err1)
		}

		var t mdl.Tournament
		err1 = json.Unmarshal(tournamentBlob, &t)
		if err1 != nil {
			log.Errorf(c, "%s unable to extract userIds from data, %v", desc, err1)
		}

		log.Infof(c, "%s value of user ids: %v", desc, userIds)
		log.Infof(c, "%s value of scores: %v", desc, scores)
		log.Infof(c, "%s value of tournament id: %v", desc, t.Id)

		log.Infof(c, "%s crunching data...", desc)
		for i, id := range userIds {
			if u, err := mdl.UserById(c, id); err != nil {
				log.Errorf(c, "%s cannot find user with id=%", desc, id)
			} else {
				log.Infof(c, "%s create score entity as it does not exist", desc)
				if se, _ := u.TournamentScore(c, &t); se == nil {
					log.Errorf(c, "%s score entity does not exist", desc)
				} else {
					log.Infof(c, "%s score entity exists, lets update it", desc)
					var err error
					if err = se.Add(c, scores[i]); err != nil {
						log.Errorf(c, "%s unable to add score of user %v, ", desc, u.Id, err)
					}
				}
			}
		}
		log.Infof(c, "%s task done!", desc)
		return nil
	}
	log.Infof(c, "%s something went wrong...")
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
