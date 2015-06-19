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
	"fmt"
	"net/http"
	"net/url"

	"appengine"
	"appengine/datastore"
	"appengine/taskqueue"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"

	mdl "github.com/santiaago/gonawin/models"
)

// UpdateScores updates the scores of all users in tournaments.
// it does this by dispatching different tasks.
//
//	GET	/a/update/scores/	update
//
// The response is ...
func UpdateScores(w http.ResponseWriter, r *http.Request /*, u *mdl.User*/) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Task queue - Update Scores Handler:"

	// err := datastore.RunInTransaction(c, func(c appengine.Context) error {
	// unable to have a transaction in this task...
	// 2014/03/20 21:03:42 ERROR: pw:  Predict.ByIds, error occurred during ByIds call: API error 1 (datastore_v3: BAD_REQUEST): operating on too many entity groups in a single transaction.
	// 2014/03/20 21:03:42 ERROR: pw: predict not found : API error 1 (datastore_v3: BAD_REQUEST): operating on too many entity groups in a single transaction.
	log.Infof(c, "%s processing...", desc)

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

	if _, err := taskqueue.Add(c, task1, "gw-queue"); err != nil {
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

	if _, err := taskqueue.Add(c, task2, "gw-queue"); err != nil {
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
	task3 := taskqueue.NewPOSTTask("/a/add/scoreentities/score/", url.Values{
		"userIds":    []string{string(buserIds)},
		"scores":     []string{string(bscores)},
		"tournament": []string{string(tournamentBlob)},
	})

	if _, err := taskqueue.Add(c, task3, "gw-queue"); err != nil {
		log.Errorf(c, "%s unable to add task to taskqueue.", desc)
		return err
	} else {
		log.Infof(c, "%s add task to taskqueue successfully", desc)
	}
	log.Infof(c, "%s task queue for adding the score to the score entity: <--", desc)

	// task queue for updating scores of users.
	log.Infof(c, "%s task queue for publishing user score activities: -->", desc)

	task4 := taskqueue.NewPOSTTask("/a/publish/users/scoreactivities/", url.Values{
		"userIds": []string{string(buserIds)},
	})

	if _, err := taskqueue.Add(c, task4, "gw-queue"); err != nil {
		log.Errorf(c, "%s unable to add task to taskqueue.", desc)
		return err
	} else {
		log.Infof(c, "%s add task to taskqueue successfully", desc)
	}
	log.Infof(c, "%s task queue for publishing user score activities: <--", desc)

	log.Infof(c, "%s task done!", desc)
	return nil
}

// UpdateUsersScores handler, use it to update users scores.
func UpdateUsersScores(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Task queue - Update Users Scores Handler:"

	log.Infof(c, "%s processing...", desc)
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
	log.Infof(c, "%s get users", desc)
	usersToUpdate := make([]*mdl.User, 0)
	for i, id := range userIds {
		if u, err := mdl.UserById(c, id); err != nil {
			log.Errorf(c, "%s cannot find user with id=%v", desc, id)
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

// CreateScoreEntities handler, use it to create the score entities.
func CreateScoreEntities(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Create score entities Handler:"

	if r.Method != "POST" {
		log.Infof(c, "%s something went wrong...", desc)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	log.Infof(c, "%s processing...", desc)
	log.Infof(c, "%s preparing data", desc)

	userIdsBlob := []byte(r.FormValue("userIds"))
	tournamentIdBlob := []byte(r.FormValue("tournamentId"))

	var userIds []int64
	errjson := json.Unmarshal(userIdsBlob, &userIds)
	if errjson != nil {
		log.Errorf(c, "%s unable to extract userIds from data, %v", desc, errjson)
	}

	var tournamentId int64
	errjson = json.Unmarshal(tournamentIdBlob, &tournamentId)
	if errjson != nil {
		log.Errorf(c, "%s unable to extract tournamentId from data, %v", desc, errjson)
	}

	log.Infof(c, "%s value of user ids: %v", desc, userIds)
	log.Infof(c, "%s value of tournamentId: %v", desc, tournamentId)

	log.Infof(c, "%s crunching data...", desc)

	users := make([]*mdl.User, 0)
	var scores []*mdl.Score
	var keyScores []*datastore.Key

	var err2 error
	log.Infof(c, "%s create score entities as it does not exist", desc)
	if scores, keyScores, err2 = mdl.CreateScores(c, userIds, tournamentId); err2 != nil {
		log.Errorf(c, "%s unable to create score entities. %v", desc, err2)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	log.Infof(c, "%s save scores", desc)
	if err := mdl.SaveScores(c, scores, keyScores); err != nil {
		log.Errorf(c, "%s unable to save score entities. %v", desc, err)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	log.Infof(c, "%s get users", desc)
	for i, id := range userIds {
		if u, err := mdl.UserById(c, id); err != nil {
			log.Errorf(c, "%s cannot find user with id=%d", desc, id)
		} else {
			log.Infof(c, "%s score ready add it to tournament %v", desc, scores[i])
			u.AddTournamentScore(c, scores[i].Id, scores[i].TournamentId)
			users = append(users, u)
		}
	}
	log.Infof(c, "%s update users", desc)
	if err := mdl.UpdateUsers(c, users); err != nil {
		log.Errorf(c, "%s unable udpate users scores: %v", desc, err)
		return errors.New(helpers.ErrorCodeUsersCannotUpdate)
	}
	log.Infof(c, "%s task done!", desc)
	return nil
}

// AddScoreToScoreEntities handler, use it to add a score to score model.
func AddScoreToScoreEntities(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)
	desc := "Task queue - Add score to score entity Handler:"

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	log.Infof(c, "%s processing...", desc)
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
	users := make([]*mdl.User, len(userIds))
	tournamentScores := make([]*mdl.Score, len(userIds))
	log.Infof(c, "%s get users", desc)
	for i, id := range userIds {
		if u, err := mdl.UserById(c, id); err != nil {
			log.Errorf(c, "%s cannot find user with id=%v", desc, id)
		} else {
			users[i] = u
		}
	}

	log.Infof(c, "%s get tournament score entities", desc)
	for i, _ := range users {
		if users[i] != nil {
			if se, err1 := users[i].TournamentScore(c, &t); se == nil {
				log.Errorf(c, "%s score entity does not exist. %v", desc, err1)
			} else {
				tournamentScores[i] = se
			}
		}
	}

	log.Infof(c, "%s add scores", desc)
	if err := mdl.AddScores(c, tournamentScores, scores); err != nil {
		log.Errorf(c, "%s cannot add scores to score entities. %v", desc, err)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeInternal)}
	}
	log.Infof(c, "%s task done!", desc)
	return nil
}

// PublishUsersScoreActivities published user score activities.
//
func PublishUsersScoreActivities(w http.ResponseWriter, r *http.Request) error {

	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	c := appengine.NewContext(r)
	desc := "Task queue - Publish Users Score Activities Handler:"

	log.Infof(c, "%s processing...", desc)
	log.Infof(c, "%s reading data...", desc)

	userIdsBlob := []byte(r.FormValue("userIds"))

	var userIds []int64
	var err error
	if err = json.Unmarshal(userIdsBlob, &userIds); err != nil {
		log.Errorf(c, "%s unable to extract userIds from data, %v", desc, err)
	}

	log.Infof(c, "%s value of user ids: %v", desc, userIds)
	log.Infof(c, "%s crunching data...", desc)

	log.Infof(c, "%s get users", desc)
	users := make([]*mdl.User, len(userIds))
	for i, id := range userIds {
		if u, err := mdl.UserById(c, id); err != nil {
			log.Errorf(c, "%s cannot find user with id=%v", desc, id)
		} else {
			users[i] = u
		}
	}

	log.Infof(c, "%s build activities", desc)
	activities := make([]*mdl.Activity, len(users))
	for i, _ := range users {
		if users[i] != nil {
			verb := fmt.Sprintf("'s score is now %d", users[i].Score)
			if a := users[i].BuildActivity(c, "score", verb, mdl.ActivityEntity{}, mdl.ActivityEntity{}); a != nil {
				activities[i] = a
			} else {
				c.Errorf("%s error when building activity.", desc)
			}
		}
	}

	log.Infof(c, "%s save activities %v", desc, activities)
	if err := mdl.SaveActivities(c, activities); err != nil {
		c.Errorf("%s something went wrong when saving all activities: error: %v", desc, err)
	}

	log.Infof(c, "%s add activity ids", desc)
	for i, _ := range activities {
		if activities[i] != nil && users[i] != nil {
			activities[i].AddNewActivityId(c, users[i])
		}
	}

	log.Infof(c, "%s update users", desc)
	if err := mdl.UpdateUsers(c, users); err != nil {
		log.Errorf(c, "%s unable udpate users scores: %v", desc, err)
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeUsersCannotUpdate)}
	}
	log.Infof(c, "%s tasks done!", desc)
	return nil
}
