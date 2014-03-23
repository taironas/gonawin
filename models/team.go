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

package models

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
)

// Accuracy of Tournament
type AccOfTournament struct {
	AccuracyId   int64 // id of accuracy entity
	TournamentId int64 // id of tournament
}

type Team struct {
	Id               int64
	KeyName          string
	Name             string
	AdminId          int64
	Private          bool
	Created          time.Time
	UserIds          []int64           // ids of Users <=> members of the team.
	TournamentIds    []int64           // ids of Tournaments <=> Tournaments the team subscribed.
	Accuracy         float64           // Overall Team accuracy.
	AccOfTournaments []AccOfTournament // ids of Accuracies for each tournament the team is participating on .
}

type TeamJson struct {
	Id            *int64             `json:",omitempty"`
	KeyName       *string            `json:",omitempty"`
	Name          *string            `json:",omitempty"`
	AdminId       *int64             `json:",omitempty"`
	Private       *bool              `json:",omitempty"`
	Created       *time.Time         `json:",omitempty"`
	UserIds       *[]int64           `json:",omitempty"`
	TournamentIds *[]int64           `json:",omitempty"`
	Accuracy      *float64           `json:",omitempty"`
	AccuracyIds   *[]AccOfTournament `json:",omitempty"`
}

// Create a team given a name, an admin id and a private mode.
func CreateTeam(c appengine.Context, name string, adminId int64, private bool) (*Team, error) {
	// create new team
	teamId, _, err := datastore.AllocateIDs(c, "Team", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Team", "", teamId, nil)
	emptyArray := make([]int64, 0)
	emtpyArrayOfAccOfTournament := make([]AccOfTournament, 0)
	team := &Team{teamId, helpers.TrimLower(name), name, adminId, private, time.Now(), emptyArray, emptyArray, float64(0), emtpyArrayOfAccOfTournament}

	_, err = datastore.Put(c, key, team)
	if err != nil {
		return nil, err
	}
	// udpate inverted index
	AddToTeamInvertedIndex(c, helpers.TrimLower(name), teamId)

	return team, err
}

// Destroy a team given a team id.
func (t *Team) Destroy(c appengine.Context) error {

	if _, err := TeamById(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find team with Id=%d", t.Id))
	} else {
		key := datastore.NewKey(c, "Team", "", t.Id, nil)

		return datastore.Delete(c, key)
	}
}

// Search for all Team entities with respect of a filter and a value.
func FindTeams(c appengine.Context, filter string, value interface{}) []*Team {

	q := datastore.NewQuery("Team").Filter(filter+" =", value)

	var teams []*Team

	if _, err := q.GetAll(c, &teams); err == nil {
		return teams
	} else {
		log.Errorf(c, " Team.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a team given an id.
func TeamById(c appengine.Context, id int64) (*Team, error) {

	var t Team
	key := datastore.NewKey(c, "Team", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " team not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// Get a team key given a team id.
func TeamKeyById(c appengine.Context, id int64) *datastore.Key {
	key := datastore.NewKey(c, "Team", "", id, nil)
	return key
}

// Update a team given an id and a team pointer.
func (t *Team) Update(c appengine.Context) error {
	// update key name
	t.KeyName = helpers.TrimLower(t.Name)
	k := TeamKeyById(c, t.Id)
	oldTeam := new(Team)
	if err := datastore.Get(c, k, oldTeam); err == nil {
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		// use lower trim names as team inverted index store them like this.
		UpdateTeamInvertedIndex(c, oldTeam.KeyName, t.KeyName, t.Id)
	}
	return nil
}

// Get all teams in datastore.
func FindAllTeams(c appengine.Context) []*Team {
	q := datastore.NewQuery("Team")

	var teams []*Team

	if _, err := q.GetAll(c, &teams); err != nil {
		log.Errorf(c, " Team.FindAll, error occurred during GetAll call: %v", err)
	}

	return teams
}

// Get an array of pointers to Teams with respect to an array of ids.
func TeamsByIds(c appengine.Context, ids []int64) []*Team {

	var teams []*Team
	for _, id := range ids {
		if team, err := TeamById(c, id); err == nil {
			teams = append(teams, team)
		} else {
			log.Errorf(c, " Team.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return teams
}

// Checks if a user has joined a team or not.
func (t *Team) Joined(c appengine.Context, u *User) bool {

	hasTeam, _ := u.ContainsTeamId(t.Id)
	return hasTeam
}

// Make a user join a team.
// TeamId is added to user entity.
// UserId is added to team entity.
// UserId is added to all tournaments joined by the team entity.
func (t *Team) Join(c appengine.Context, u *User) error {
	// add
	log.Infof(c, "Team.Join: user")
	log.Infof(c, "Team.Join: add team id to user entity")
	if err := u.AddTeamId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Team.Join, error joining tournament for user:%v Error: %v", u.Id, err))
	}
	log.Infof(c, "Team.Join: add user id to team entity")
	if err := t.AddUserId(c, u.Id); err != nil {
		return errors.New(fmt.Sprintf(" Team.Join, error joining tournament for user:%v Error: %v", u.Id, err))
	}
	log.Infof(c, "Team.Join: add user id to tournaments")
	if err := t.AddUserToTournaments(c, u.Id); err != nil {
		return errors.New(fmt.Sprintf("Team.Join, error adding user:%v to teams tournaments Error: %v", u.Id, err))
	}
	return nil
}

// make a user leave a team
// Todo: Should we check that the user is indeed a memeber of the team?
func (t *Team) Leave(c appengine.Context, u *User) error {
	if err := u.RemoveTeamId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Team.Leave, error leaving team for user:%v Error: %v", u.Id, err))
	}
	if err := t.RemoveUserId(c, u.Id); err != nil {
		return errors.New(fmt.Sprintf(" Team.Leave, error leaving team for user:%v Error: %v", u.Id, err))
	}
	// sar: 7 mar 2014
	// when a user leaves a team should we unsubscribe him from the tournaments of that team?
	// for now I would say no.

	// if err := t.RemoveUserFromTournaments(c, u); err != nil{
	// 	return errors.New(fmt.Sprintf("Team.Leave, error leaving teams tournaments for user:%v Error: %v", u.Id, err))
	// }
	return nil
}

// Check if user is admin of the team with id 'teamId'.
func IsTeamAdmin(c appengine.Context, teamId int64, userId int64) bool {

	if team, err := TeamById(c, teamId); err == nil {
		return team.AdminId == userId
	} else {
		log.Errorf(c, " Team.IsTeamAdmin, error occurred during ById call: %v", err)
		return false
	}
}

// Given a id, and a word, get the frequency of that word in the team terms.
func GetWordFrequencyForTeam(c appengine.Context, id int64, word string) int64 {

	if teams := FindTeams(c, "Id", id); teams != nil {
		return helpers.CountTerm(strings.Split(teams[0].KeyName, " "), word)
	}
	return 0
}

// from a team id return an array of users/ players that participates in it.
func (t *Team) Players(c appengine.Context) []*User {

	var users []*User

	for _, uId := range t.UserIds {
		user, err := UserById(c, uId)
		if err != nil {
			log.Errorf(c, " Players, cannot find user with ID=%", uId)
		} else {
			users = append(users, user)
		}
	}
	return users
}

func (t *Team) ContainsTournamentId(id int64) (bool, int) {

	for i, tId := range t.TournamentIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// Adds a tournament Id in the TournamentId array.
func (t *Team) AddTournamentId(c appengine.Context, tId int64) error {
	log.Infof(c, "team Add tournament id")
	if hasTournament, _ := t.ContainsTournamentId(tId); hasTournament {
		log.Infof(c, "team Add tournament id all ready member")
		return errors.New(fmt.Sprintf("AddTournamentId, allready a member."))
	}
	log.Infof(c, "team Add tournament id append tournament ids")
	t.TournamentIds = append(t.TournamentIds, tId)
	if err := t.Update(c); err != nil {
		return err
	}
	log.Infof(c, "team Add tournament id loop through all users and update arrays")
	for _, uId := range t.UserIds {
		user, err := UserById(c, uId)
		if err != nil {
			log.Errorf(c, "Team.AddTournamentId, user not found")
		} else {
			log.Infof(c, "team Add tournament id add tournament id%v", user.Id)
			user.AddTournamentId(c, tId)
		}
	}
	return nil
}

// Remove a tournament Id in the TournamentId array.
func (t *Team) RemoveTournamentId(c appengine.Context, tId int64) error {

	if hasTournament, i := t.ContainsTournamentId(tId); !hasTournament {
		return errors.New(fmt.Sprintf("RemoveTournamentId, not a member."))
	} else {
		// as the order of index in tournamentsId is not important,
		// replace elem at index i with last element and resize slice.
		t.TournamentIds[i] = t.TournamentIds[len(t.TournamentIds)-1]
		t.TournamentIds = t.TournamentIds[0 : len(t.TournamentIds)-1]
	}
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// from a team return an array of tournament the user is involved in.
func (t *Team) Tournaments(c appengine.Context) []*Tournament {

	var tournaments []*Tournament

	for _, tId := range t.TournamentIds {
		tournament, err := TournamentById(c, tId)
		if err != nil {
			log.Errorf(c, " Tournaments, cannot find team with ID=%", tId)
		} else {
			tournaments = append(tournaments, tournament)
		}
	}

	return tournaments
}

// from a team return an array of accuracy the user is involved in.
func (t *Team) Accuracies(c appengine.Context) []*Accuracy {

	var accs []*Accuracy

	for _, acc := range t.AccOfTournaments {
		a, err := AccuracyById(c, acc.AccuracyId)
		if err != nil {
			log.Errorf(c, " Accuracies, cannot find accuracy with ID=%", acc.AccuracyId)
		} else {
			accs = append(accs, a)
		}
	}
	return accs
}

// Remove a user Id in the UserId array.
func (t *Team) RemoveUserId(c appengine.Context, uId int64) error {

	if hasUser, i := t.ContainsUserId(uId); !hasUser {
		return errors.New(fmt.Sprintf("RemoveUserId, not a member."))
	} else {
		// as the order of index in usersId is not important,
		// replace elem at index i with last element and resize slice.
		t.UserIds[i] = t.UserIds[len(t.UserIds)-1]
		t.UserIds = t.UserIds[0 : len(t.UserIds)-1]
	}
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a team Id in the UserId array.
func (t *Team) AddUserId(c appengine.Context, uId int64) error {

	if hasUser, _ := t.ContainsUserId(uId); hasUser {
		return errors.New(fmt.Sprintf("AddUserId, allready a member."))
	}

	t.UserIds = append(t.UserIds, uId)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Add user to teams tournaments
func (t *Team) AddUserToTournaments(c appengine.Context, uId int64) error {

	log.Infof(c, "Team.AddUserToTournaments")
	for _, tId := range t.TournamentIds {
		if tournament, err := TournamentById(c, tId); err != nil {
			log.Errorf(c, "Cannot find tournament with Id=%d", t.Id)
		} else {
			log.Infof(c, "Team.AddUserToTournaments add user id to tournament")
			if err := tournament.AddUserId(c, uId); err != nil {
				log.Errorf(c, "Team.AddUserToTournaments: unable to add user:%v to tournament:%v", uId, tId)
			}
			log.Infof(c, "Team.AddUserToTournaments get user")
			if u, err := UserById(c, uId); err != nil {
				log.Errorf(c, "User not found %v", uId)
			} else {
				log.Infof(c, "Team.AddUserToTournaments add tournament id for user %v", u.Id)
				if err1 := u.AddTournamentId(c, tournament.Id); err1 != nil {
					log.Errorf(c, "Team.AddUserToTournaments: unable to add tournament id:%v to user:%v", tId, uId)
				} else {
					log.Infof(c, "Team.AddUserToTournaments add tournament id for user %v successfully", u.Id)
				}
			}
		}
	}
	return nil
}

func (t *Team) ContainsUserId(id int64) (bool, int) {

	for i, tId := range t.UserIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// Update an array of teams.
func UpdateTeams(c appengine.Context, teams []*Team) error {
	keys := make([]*datastore.Key, len(teams))
	for i, _ := range keys {
		keys[i] = TeamKeyById(c, teams[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, teams); err != nil {
		return err
	}
	return nil
}

func (t *Team) RankingByUser(c appengine.Context, limit int) []*User {
	if limit < 0{
		return nil
	}
	users := t.Players(c)
	sort.Sort(UserByScore(users))
	if len(users) <= limit{
		return users
	} else {
		return users[0:limit]
	}
}

// Sort teams by score
type TeamByAccuracy []*Team

func (a TeamByAccuracy) Len() int           { return len(a) }
func (a TeamByAccuracy) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TeamByAccuracy) Less(i, j int) bool { return a[i].Accuracy < a[j].Accuracy }

func (t *Team) TournamentAcc(c appengine.Context, tournament *Tournament) (*Accuracy, error) {
	//query accuracy
	for _, acc := range t.AccOfTournaments {
		if acc.TournamentId == tournament.Id {
			return AccuracyById(c, acc.AccuracyId)
		}
	}
	return nil, errors.New("model/team: accuracy not found")
}

// add accuracy to team entity and run update.
func (t *Team) AddTournamentAcc(c appengine.Context, accId int64, tourId int64) error {

	tournamentExist := false
	for _, tid := range t.TournamentIds {
		if tid == tourId {
			tournamentExist = true
			break
		}
	}
	if !tournamentExist {
		return errors.New("model/team: not member of tournament")
	}
	accExist := false
	for _, acc := range t.AccOfTournaments {
		if acc.AccuracyId == accId {
			accExist = true
			break
		}
	}
	if accExist {
		return errors.New("model/team: accuracy allready present")
	}

	var a AccOfTournament
	a.AccuracyId = accId
	a.TournamentId = tourId
	t.AccOfTournaments = append(t.AccOfTournaments, a)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Update the global accuracy for team, with new accuracy and accuracies of other tournaments.
// From all accuracies of tournaments the team participates on we sum the overall accuracies.
// Overall accuracy is it's last element in the array of accuracies. We then normalize by the number of tournaments the teams participates on.
func (t *Team) UpdateAccuracy(c appengine.Context, tId int64, newAcc float64) error {
	log.Infof(c, "Updating accuracy of team")

	sum := float64(0)
	counter := 0
	for _, accOfTournament := range t.AccOfTournaments {
		if accOfTournament.TournamentId == tId {
			sum += newAcc
			counter++
			continue
		}
		if acc, err := AccuracyById(c, accOfTournament.AccuracyId); err == nil && acc != nil {
			// only take into account tournaments with accuracies
			if len(acc.Accuracies) > 0 {
				sum += acc.Accuracies[len(acc.Accuracies)-1]
				counter++
			}
		} else if err != nil {
			log.Infof(c, "Accuracy not found %v, error:", accOfTournament.AccuracyId, err)
		}
	}
	if counter > 0 {
		t.Accuracy = sum / float64(counter)
		if err := t.Update(c); err != nil {
			log.Infof(c, "Team.UpdateAccuracyAccuracy: unable to update team %v", err)
			return err
		} else {
			// publish new activity
			verb := fmt.Sprintf("has a new accuracy of %.2f%%", newAcc*100)
			t.Publish(c, "accuracy", verb, ActivityEntity{}, ActivityEntity{})
		}
	}
	return nil
}

// Find all activities
func (t *Team) Activities(c appengine.Context) Activities {
	q := datastore.NewQuery("Activity").Filter("CreatorID=", t.Id).Order("-Published")

	var activities []*Activity

	if _, err := q.GetAll(c, &activities); err != nil {
		log.Errorf(c, "model/team, Activities: error occurred during GetAll call: %v", err)
	}

	return activities
}

func (t *Team) Publish(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) error {
	var activity Activity
	activity.Type = activityType
	activity.Verb = verb
	activity.Actor = ActivityEntity{ID: t.Id, Type: "team", DisplayName: t.Name}
	activity.Object = object
	activity.Target = target
	activity.Published = time.Now()
	activity.CreatorID = t.Id

	return activity.save(c)
}

// Get an array of type accuracyOverall which holds the accuracy information and the last 5 progression of each tournament.
func (t *Team) AccuraciesGroupByTournament(c appengine.Context, limit int) *[]AccuracyOverall {
	accs := make([]AccuracyOverall, 0)
	for _, aot := range t.AccOfTournaments {
		if acc, err := AccuracyById(c, aot.AccuracyId); err != nil {
			log.Errorf(c, "Team.AccuraciesByTournament: Unable to retreive accuracy entity from id, ", err)
		} else {
			var a AccuracyOverall
			a.Id = aot.AccuracyId
			a.TournamentId = aot.TournamentId
			if len(acc.Accuracies) > 0 {
				a.Accuracy = acc.Accuracies[len(acc.Accuracies)-1]
			} else {
				a.Accuracy = 1
			}
			a.Progression = make([]Progression, 0)
			counter := 0
			// get most recent accuracies with a 'limit' progression.
			for i := len(acc.Accuracies) - 1; i > -1; i-- {
				cur := acc.Accuracies[i]
				var prog Progression
				prog.Value = cur
				a.Progression = append(a.Progression, prog)
				counter++
				if counter == limit {
					break
				}
			}
			accs = append(accs, a)
		}
	}
	return &accs
}

// Get the overall accuracy of a team in the specified tournament.
// the progression of accuracies is in reverse order to have the most reset accuracy as first element.
func (t *Team) AccuracyByTournament(c appengine.Context, tour *Tournament) *AccuracyOverall {
	for _, aot := range t.AccOfTournaments {
		if aot.TournamentId != tour.Id {
			continue
		}
		if acc, err := AccuracyById(c, aot.AccuracyId); err != nil {
			log.Errorf(c, "Team.AccuraciesByTournament: Unable to retreive accuracy entity from id, ", err)
		} else {
			var a AccuracyOverall
			a.Id = aot.AccuracyId
			a.TournamentId = aot.TournamentId
			if len(acc.Accuracies) > 0 {
				a.Accuracy = acc.Accuracies[len(acc.Accuracies)-1]
			} else {
				a.Accuracy = 1
			}
			a.Progression = make([]Progression, len(acc.Accuracies))
			for i, cur := range acc.Accuracies {
				var prog Progression
				prog.Value = cur
				a.Progression[i] = prog
			}
			return &a
		}
	}
	return nil

}
