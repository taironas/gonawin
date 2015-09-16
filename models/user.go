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

package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"
	gaeuser "appengine/user"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
)

// Accuracy of Tournament
type ScoreOfTournament struct {
	ScoreId      int64 // id of score entity
	TournamentId int64 // id of tournament
}

type User struct {
	Id                    int64
	Email                 string
	Username              string
	Name                  string
	Alias                 string              // name to display chosen by user if requested.
	IsAdmin               bool                // is user gonawin admin.
	Auth                  string              // authentication auth token
	PredictIds            []int64             // current user predicts.
	ArchivedPredictInds   []int64             // archived user predicts.
	TournamentIds         []int64             // current tournament ids of user <=> tournaments user subscribed.
	ArchivedTournamentIds []int64             // archived tournament ids of user <=> finnished tournametns user subscribed.
	TeamIds               []int64             // current team ids of user <=> teams user belongs to.
	Score                 int64               // overall user score.
	ScoreOfTournaments    []ScoreOfTournament // ids of Scores for each tournament the user is participating on.
	ActivityIds           []int64             // ids of user's activities
	Created               time.Time
}

type UserJson struct {
	Id                    *int64               `json:",omitempty"`
	Email                 *string              `json:",omitempty"`
	Username              *string              `json:",omitempty"`
	Name                  *string              `json:",omitempty"`
	Alias                 *string              `json:",omitempty"`
	IsAdmin               *bool                `json:",omitempty"`
	Auth                  *string              `json:",omitempty"`
	PredictIds            *[]int64             `json:",omitempty"`
	ArchivedPredictInds   *[]int64             `json:",omitempty"`
	TournamentIds         *[]int64             `json:",omitempty"`
	ArchivedTournamentIds *[]int64             `json:",omitempty"`
	TeamIds               *[]int64             `json:",omitempty"`
	Score                 *int64               `json:",omitempty"`
	ScoreOfTournaments    *[]ScoreOfTournament `json:",omitempty"`
	ActivityIds           *[]int64             `json:",omitempty"`
	Created               *time.Time           `json:",omitempty"`
}

// CreateUser let you create a user entity.
//
func CreateUser(c appengine.Context, email, username, name, alias string, isAdmin bool, auth string) (*User, error) {

	userId, _, err := datastore.AllocateIDs(c, "User", nil, 1)
	if err != nil {
		log.Errorf(c, " User.Create: %v", err)
	}

	key := datastore.NewKey(c, "User", "", userId, nil)

	emptyArray := make([]int64, 0)
	emptyScores := make([]ScoreOfTournament, 0)
	user := &User{
		Id:                    userId,
		Email:                 email,
		Username:              username,
		Name:                  name,
		Alias:                 alias,
		IsAdmin:               isAdmin,
		Auth:                  auth,
		PredictIds:            emptyArray,
		ArchivedPredictInds:   emptyArray,
		TournamentIds:         emptyArray,
		ArchivedTournamentIds: emptyArray,
		TeamIds:               emptyArray,
		Score:                 int64(0),
		ScoreOfTournaments:    emptyScores,
		ActivityIds:           emptyArray,
		Created:               time.Now(),
	}

	if _, err = datastore.Put(c, key, user); err != nil {
		log.Errorf(c, "User.Create: %v", err)
		return nil, errors.New("model/user: Unable to put user in Datastore")
	}

	// add name to inverted index
	// as name and username can have the same words.
	// We build a string with a set of words between these two strings
	allnames := name + " " + username
	setOfNames := helpers.SetOfStrings(allnames)
	names := ""
	for _, n := range setOfNames {
		names = names + " " + n
	}
	AddToUserInvertedIndex(c, names, user.Id)

	return user, nil
}

// Destroy lets you remove a user from the data store given a user id.
//
func (u *User) Destroy(c appengine.Context) error {

	if _, err := UserById(c, u.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find user with Id=%d", u.Id))
	} else {
		key := datastore.NewKey(c, "User", "", u.Id, nil)

		if errd := datastore.Delete(c, key); errd != nil {
			return errd
		} else {
			// remove key name.
			if err = UpdateUserInvertedIndex(c, helpers.TrimLower(u.Name), "", u.Id); err != nil {
				return err
			}
			// remove key username.
			if err = UpdateUserInvertedIndex(c, helpers.TrimLower(u.Username), "", u.Id); err != nil {
				return err
			}
			// remove key alias.
			if err = UpdateUserInvertedIndex(c, helpers.TrimLower(u.Alias), "", u.Id); err != nil {
				return err
			}

			return nil
		}
	}
}

// Search for a user entity given a filter and value.
func FindUser(c appengine.Context, filter string, value interface{}) *User {

	q := datastore.NewQuery("User").Filter(filter+" =", value)
	var users []*User
	if _, err := q.GetAll(c, &users); err == nil && len(users) > 0 {
		return users[0]
	} else if len(users) == 0 {
		log.Infof(c, " User.Find, error occurred during GetAll")
	} else {
		log.Errorf(c, " User.Find, error occurred during GetAll: %v", err)
	}
	return nil
}

// Find all users present in datastore.
func FindAllUsers(c appengine.Context) []*User {
	q := datastore.NewQuery("User")
	var users []*User
	if _, err := q.GetAll(c, &users); err != nil {
		log.Errorf(c, "FindAllUser, error occurred during GetAll call: %v", err)
	}
	return users
}

// UserById finds a user entity by id.
//
func UserById(c appengine.Context, id int64) (*User, error) {

	var u User
	key := datastore.NewKey(c, "User", "", id, nil)
	if err := datastore.Get(c, key, &u); err != nil {
		log.Errorf(c, " user not found : %v", err)
		return nil, err
	}
	return &u, nil
}

// UsersByIds returns an array of pointers to Users with respect to an array of ids.
// It only return the found users.
//
func UsersByIds(c appengine.Context, ids []int64) ([]*User, error) {

	users := make([]User, len(ids))
	keys := UserKeysByIds(c, ids)

	var wrongIndexes []int
	if err := datastore.GetMulti(c, keys, users); err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == datastore.ErrNoSuchEntity {
					log.Errorf(c, "UsersByIds, missing key: %v %v", err, keys[i].IntID())
					wrongIndexes = append(wrongIndexes, i)
				}
			}
		} else {
			return nil, err
		}
	}

	var existingUsers []*User
	for i := range users {
		if !contains(wrongIndexes, i) {
			log.Infof(c, "UsersByIds %v", users[i])
			existingUsers = append(existingUsers, &users[i])
		}
	}
	return existingUsers, nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func UserKeysByIds(c appengine.Context, ids []int64) []*datastore.Key {
	keys := make([]*datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = UserKeyById(c, id)
	}
	return keys
}

// Get key pointer given a user id.
func UserKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "User", "", id, nil)
	return key
}

// Update user given a user pointer.
func (u *User) Update(c appengine.Context) error {
	k := UserKeyById(c, u.Id)
	oldUser := new(User)
	if err := datastore.Get(c, k, oldUser); err == nil {
		if _, err := datastore.Put(c, k, u); err != nil {
			return err
		}
		// use lower trim names for alias as user inverted index store them like this.
		// alias is the only field that can be changed.
		UpdateUserInvertedIndex(c, helpers.TrimLower(oldUser.Alias), helpers.TrimLower(u.Alias), u.Id)
	}
	return nil
}

// SigninUser saves a user from given parameters in the datastore and return a pointer to it.
func SigninUser(c appengine.Context, queryName string, email string, username string, name string) (*User, error) {

	var user *User

	queryValue := ""
	if queryName == "Email" {
		queryValue = strings.ToLower(email)
	} else if queryName == "Username" {
		queryValue = username
	} else {
		return nil, errors.New("no valid query name.")
	}

	// find user
	if user = FindUser(c, queryName, queryValue); user == nil {
		// create user if it does not exist

		isAdmin := gaeuser.IsAdmin(c)

		// start with an empty alias.
		alias := ""
		if userCreate, err := CreateUser(c, email, username, name, alias, isAdmin, GenerateAuthKey()); err != nil {
			log.Errorf(c, "Signup: %v", err)
			return nil, errors.New("unable to create user.")
		} else {
			user = userCreate
		}
		// publish new activity
		user.Publish(c, "welcome", "joined gonawin", ActivityEntity{}, ActivityEntity{})
	}

	return user, nil
}

// Generate authentication string key.
// We use this function to create an authentication token for a user entity.
func GenerateAuthKey() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

// From a user id returns an array of teams the user iq involved participates.
func (u *User) Teams(c appengine.Context) []*Team {

	var teams []*Team

	for _, tId := range u.TeamIds {
		t, err := TeamById(c, tId)
		if err != nil {
			log.Errorf(c, "Teams, cannot find team with ID=%", tId)
		} else {
			teams = append(teams, t)
		}
	}

	return teams
}

// From a user id returns an array of teams the user is involved participates.
func (u *User) TeamsByPage(c appengine.Context, count, page int64) []*Team {
	desc := "User.TeamsByPage"
	teams := u.Teams(c)
	// loop backward on all of these ids to fetch the teams
	log.Infof(c, "%s calculateStartAndEnd(%v, %v, %v)", desc, int64(len(teams)), count, page)
	start, end := calculateStartAndEnd(int64(len(teams)), count, page)
	log.Infof(c, "%s start = %d, end = %d", desc, start, end)
	var paged []*Team
	for i := start; i >= end; i-- {
		paged = append(paged, teams[i])
	}
	return paged
}

// From a user id returns an array of tournaments the user is involved participates.
func (u *User) TournamentsByPage(c appengine.Context, count, page int64) []*Tournament {
	desc := "User.TournamentsByPage"
	tournaments := u.Tournaments(c)
	// loop backward on all of these ids to fetch the teams
	log.Infof(c, "%s calculateStartAndEnd(%v, %v, %v)", desc, int64(len(tournaments)), count, page)
	start, end := calculateStartAndEnd(int64(len(tournaments)), count, page)
	log.Infof(c, "%s start = %d, end = %d", desc, start, end)
	var paged []*Tournament
	for i := start; i >= end; i-- {
		paged = append(paged, tournaments[i])
	}
	return paged
}

// Adds a predict Id in the PredictId array.
func (u *User) AddPredictId(c appengine.Context, pId int64) error {

	u.PredictIds = append(u.PredictIds, pId)
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a tournament Id in the TournamentId array.
func (u *User) AddTournamentId(c appengine.Context, tId int64) error {
	log.Infof(c, "User Add tournament id")
	if hasTournament, _ := u.ContainsTournamentId(tId); hasTournament {
		return errors.New(fmt.Sprintf("AddTournamentId, allready a member."))
	}
	log.Infof(c, "User Add tournament id append")
	u.TournamentIds = append(u.TournamentIds, tId)
	if err := u.Update(c); err != nil {
		log.Infof(c, "User Add tournament id unable to update user.")
		return err
	}
	log.Infof(c, "User Add tournament id all good")
	return nil
}

// Removes a tournament Id in the TournamentId array.
func (u *User) RemoveTournamentId(c appengine.Context, tId int64) error {

	if hasTournament, i := u.ContainsTournamentId(tId); !hasTournament {
		return errors.New(fmt.Sprintf("RemoveTournamentId, not a member."))
	} else {
		// as the order of index in tournamentsId is not important,
		// replace elem at index i with last element and resize slice.
		u.TournamentIds[i] = u.TournamentIds[len(u.TournamentIds)-1]
		u.TournamentIds = u.TournamentIds[0 : len(u.TournamentIds)-1]
	}
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

func (u *User) ContainsTournamentId(id int64) (bool, int) {

	for i, tId := range u.TournamentIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// from a user return an array of tournament the user is involved in.
func (u *User) Tournaments(c appengine.Context) []*Tournament {

	var tournaments []*Tournament
	var err error

	if tournaments, err = TournamentsByIds(c, u.TournamentIds); err != nil {
		log.Errorf(c, "Something failed when calling TournamentsByIds from user.Tournaments: %v", err)
	}

	return tournaments
}

// Adds a team Id in the TeamId array.
func (u *User) AddTeamId(c appengine.Context, tId int64) error {

	if hasTeam, _ := u.ContainsTeamId(tId); hasTeam {
		return errors.New(fmt.Sprintf("AddTeamId, allready a member."))
	}

	u.TeamIds = append(u.TeamIds, tId)
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

// Removes a team Id in the TeamId array.
func (u *User) RemoveTeamId(c appengine.Context, tId int64) error {

	if hasTeam, i := u.ContainsTeamId(tId); !hasTeam {
		return errors.New(fmt.Sprintf("RemoveTeamId, not a member."))
	} else {
		// as the order of index in teamsId is not important,
		// replace elem at index i with last element and resize slice.
		u.TeamIds[i] = u.TeamIds[len(u.TeamIds)-1]
		u.TeamIds = u.TeamIds[0 : len(u.TeamIds)-1]
	}
	if err := u.Update(c); err != nil {
		return err
	}
	return nil
}

func (u *User) ContainsTeamId(id int64) (bool, int) {

	for i, tId := range u.TeamIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// Update an array of users.
func UpdateUsers(c appengine.Context, users []*User) error {
	keys := make([]*datastore.Key, len(users))
	for i, _ := range keys {
		keys[i] = UserKeyById(c, users[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, users); err != nil {
		return err
	}
	return nil
}

// PredictFromMatchId returns the user predictions for a specific match.
//
func (u *User) PredictFromMatchId(c appengine.Context, mId int64) (*Predict, error) {

	var predicts []*Predict
	var err error
	if predicts, err = PredictsByIds(c, u.PredictIds); err != nil {
		return nil, err
	}

	for i, p := range predicts {
		if p.MatchId == mId {
			return predicts[i], nil
		}
	}
	return nil, nil
}

func (u *User) ScoreForMatch(c appengine.Context, m *Tmatch) (int64, error) {
	desc := "Score for match:"
	log.Infof(c, "%s teamA: %v - teamB: %v", desc, m.TeamId1, m.TeamId2)
	log.Infof(c, "%s result: %v - %v", desc, m.Result1, m.Result2)
	var p *Predict
	var err1 error
	if p, err1 = u.PredictFromMatchId(c, m.Id); err1 == nil && p == nil {
		log.Infof(c, "%s no predict for match %v was found in user %v account", desc, m.Id, u.Id)
		return 0, nil
	} else if err1 != nil {
		log.Errorf(c, "%s unable to get predict for current user %v: %v", desc, u.Id, err1)
		return 0, nil
	}
	log.Infof(c, "%s predict found, now computing score", desc)
	return computeScore(c, m, p), nil
}

// Sort users by score
type UserByScore []*User

func (a UserByScore) Len() int           { return len(a) }
func (a UserByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UserByScore) Less(i, j int) bool { return a[i].Score < a[j].Score }

// Publish user activity
func (u *User) Publish(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) error {
	activity := u.BuildActivity(c, activityType, verb, object, target)

	if err := activity.save(c); err != nil {
		return err
	}
	// add new activity id in user activity table
	activity.AddNewActivityId(c, u)

	return u.Update(c)
}

func (u *User) BuildActivity(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) *Activity {
	var activity Activity
	activity.Type = activityType
	activity.Verb = verb
	activity.Actor = u.Entity()
	activity.Object = object
	activity.Target = target
	activity.Published = time.Now()
	activity.CreatorID = u.Id
	id, _, err1 := datastore.AllocateIDs(c, "Activity", nil, 1)
	if err1 != nil {
		log.Errorf(c, " BuildActivity: error occurred during AllocateIDs call: %v", err1)
		return nil
	}
	activity.Id = id
	return &activity
}

// Activity entity representation of an user
func (u *User) Entity() ActivityEntity {
	return ActivityEntity{Id: u.Id, Type: "user", DisplayName: u.Username}
}

func (u *User) TournamentScore(c appengine.Context, tournament *Tournament) (*Score, error) {
	//query score
	for _, s := range u.ScoreOfTournaments {
		if s.TournamentId == tournament.Id {
			log.Infof(c, "User.TournamentScore tournament found in ScoreOfTournaments array")
			return ScoreById(c, s.ScoreId)
		}
	}
	log.Infof(c, "User.TournamentScore score entity not found")
	return nil, errors.New("model/team: score entity not found")
}

// add accuracy to team entity and run update.
func (u *User) AddTournamentScore(c appengine.Context, scoreId int64, tourId int64) error {
	log.Infof(c, "model/user: add tournament score")
	tournamentExist := false
	for _, tid := range u.TournamentIds {
		if tid == tourId {
			tournamentExist = true
			break
		}
	}
	if !tournamentExist {
		log.Infof(c, "model/user: add tournament score, tournament does not exist")
		return errors.New("model/team: not member of tournament")
	}
	scoreExist := false
	for _, s := range u.ScoreOfTournaments {
		if s.ScoreId == scoreId {
			scoreExist = true
			break
		}
	}
	if scoreExist {
		log.Infof(c, "model/user: add tournament score, score entity already member")
		return errors.New("model/team: score allready present")
	}
	log.Infof(c, "model/user: add tournament score, create score entity")

	var s ScoreOfTournament
	s.ScoreId = scoreId
	s.TournamentId = tourId
	u.ScoreOfTournaments = append(u.ScoreOfTournaments, s)
	return nil
}

// Returns an array of score entities group by tournament.
func (u *User) Scores(c appengine.Context) []*Score {

	scores := make([]*Score, 0)
	for _, s := range u.ScoreOfTournaments {
		if score, err := ScoreById(c, s.ScoreId); err != nil {
			log.Errorf(c, "User.Scores: error when calling ScoreById")
		} else {
			scores = append(scores, score)
		}
	}
	return scores
}

// Get the score of user with respect to tournament.
// If tournament not found return 0.
func (u *User) ScoreByTournament(c appengine.Context, tId int64) int64 {
	for _, s := range u.ScoreOfTournaments {
		if s.TournamentId == tId {
			if score, err := ScoreById(c, s.ScoreId); err == nil {
				return sumInt64(&score.Scores)
			}
		}
	}
	return int64(0)
}

// Returns an array of scoreOverall entities group by tournament.
func (u *User) TournamentsScores(c appengine.Context) []*ScoreOverall {

	scores := make([]*ScoreOverall, 0)
	for _, s := range u.ScoreOfTournaments {
		if score, err := ScoreById(c, s.ScoreId); err != nil {
			log.Errorf(c, "User.Scores: error when calling ScoreById")
		} else {
			var so ScoreOverall
			so.Id = score.Id
			so.UserId = score.UserId
			so.TournamentId = score.TournamentId
			so.Score = sumInt64(&score.Scores)
			if len(score.Scores) > 0 {
				so.LastProgression = score.Scores[len(score.Scores)-1]
			}
			scores = append(scores, &so)
		}
	}
	return scores
}

// Get the invitations of a user.
func (u *User) Invitations(c appengine.Context) []*Team {
	urs := FindUserRequests(c, "UserId", u.Id)
	var ids []int64
	for _, ur := range urs {
		ids = append(ids, ur.TeamId)
	}

	var teams []*Team
	var err error
	if teams, err = TeamsByIds(c, ids); err != nil {
		log.Errorf(c, "User.Invitations: something failed when calling TeamsByIds: %v", err)
	}
	return teams
}

// Find all entity users with respect of a filter and value.
func FindUsers(c appengine.Context, filter string, value interface{}) []*User {

	q := datastore.NewQuery("User").Filter(filter+" =", value)
	var users []*User
	if _, err := q.GetAll(c, &users); err == nil {
		return users
	} else {
		log.Errorf(c, "FindUsers, error occurred during GetAll: %v", err)
		return nil
	}
}

// Get the frequency of given word with respect to user id.
func GetWordFrequencyForUser(c appengine.Context, id int64, word string) int64 {

	if users := FindUsers(c, "Id", id); users != nil {
		return helpers.CountTerm(strings.Split(users[0].Name, " "), word)
	}
	return 0
}
