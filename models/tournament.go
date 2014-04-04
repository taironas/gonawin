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
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
)

type Tournament struct {
	Id              int64
	KeyName         string
	Name            string
	Description     string
	Start           time.Time
	End             time.Time
	AdminId         int64
	Created         time.Time
	GroupIds        []int64
	Matches1stStage []int64
	Matches2ndStage []int64
	UserIds         []int64
	TeamIds         []int64
}

type TournamentJson struct {
	Id              *int64     `json:",omitempty"`
	KeyName         *string    `json:",omitempty"`
	Name            *string    `json:",omitempty"`
	Description     *string    `json:",omitempty"`
	Start           *time.Time `json:",omitempty"`
	End             *time.Time `json:",omitempty"`
	AdminId         *int64     `json:",omitempty"`
	Created         *time.Time `json:",omitempty"`
	GroupIds        *[]int64   `json:",omitempty"`
	Matches1stStage *[]int64   `json:",omitempty"`
	Matches2ndStage *[]int64   `json:",omitempty"`
	UserIds         *[]int64   `json:",omitempty"`
	TeamIds         *[]int64   `json:",omitempty"`
}

// Create tournament entity given a name and description.
func CreateTournament(c appengine.Context, name string, description string, start time.Time, end time.Time, adminId int64) (*Tournament, error) {

	tournamentID, _, err := datastore.AllocateIDs(c, "Tournament", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Tournament", "", tournamentID, nil)

	// empty groups and tournaments for now
	emptyArray := make([]int64, 0)

	tournament := &Tournament{tournamentID, helpers.TrimLower(name), name, description, start, end, adminId, time.Now(), emptyArray, emptyArray, emptyArray, emptyArray, emptyArray}

	_, err = datastore.Put(c, key, tournament)
	if err != nil {
		return nil, err
	}

	AddToTournamentInvertedIndex(c, helpers.TrimLower(name), tournamentID)
	return tournament, nil
}

// Destroy a tournament entity given a tournament id.
func (t *Tournament) Destroy(c appengine.Context) error {

	if _, err := TournamentById(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf("Cannot find tournament with Id=%d", t.Id))
	} else {
		key := datastore.NewKey(c, "Tournament", "", t.Id, nil)
		return datastore.Delete(c, key)
	}
}

// Find all entity tournaments with respect of a filter and value.
func FindTournaments(c appengine.Context, filter string, value interface{}) []*Tournament {

	q := datastore.NewQuery("Tournament").Filter(filter+" =", value)
	var tournaments []*Tournament
	if _, err := q.GetAll(c, &tournaments); err == nil {
		return tournaments
	} else {
		log.Errorf(c, " Tournament.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// Get a pointer to a tournament given a tournament id.
func TournamentById(c appengine.Context, id int64) (*Tournament, error) {

	var t Tournament
	key := datastore.NewKey(c, "Tournament", "", id, nil)
	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " tournament not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// Get a pointer to a tournament key given a tournament id.
func TournamentKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Tournament", "", id, nil)
	return key
}

// Update a tournament given a tournament id and a tournament pointer.
func (t *Tournament) Update(c appengine.Context) error {

	// update key name
	t.KeyName = helpers.TrimLower(t.Name)
	k := TournamentKeyById(c, t.Id)
	oldTournament := new(Tournament)
	if err := datastore.Get(c, k, oldTournament); err == nil {
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		// use name with trim lower as tournament inverted index stores lower key names.
		UpdateTournamentInvertedIndex(c, oldTournament.KeyName, t.KeyName, t.Id)
	}
	return nil
}

// Find all tournaments in the datastore.
func FindAllTournaments(c appengine.Context) []*Tournament {

	q := datastore.NewQuery("Tournament")
	var tournaments []*Tournament
	if _, err := q.GetAll(c, &tournaments); err != nil {
		log.Errorf(c, "FindAllTournaments, error occurred during GetAll call: %v", err)
	}
	return tournaments
}

// Find all tournaments with respect to array of ids.
func TournamentsByIds(c appengine.Context, ids []int64) []*Tournament {

	var tournaments []*Tournament
	for _, id := range ids {
		if tournament, err := TournamentById(c, id); err == nil {
			tournaments = append(tournaments, tournament)
		} else {
			log.Errorf(c, " Tournament.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return tournaments
}

// Checks if a user has joined a tournament.
func (t *Tournament) Joined(c appengine.Context, u *User) bool {
	// change in contains
	hasTournament, _ := u.ContainsTournamentId(t.Id)
	return hasTournament
}

// Makes a user join a tournament.
func (t *Tournament) Join(c appengine.Context, u *User) error {
	// add
	if err := u.AddTournamentId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.Join, error joining tournament for user:%v Error: %v", u.Id, err))
	}
	if err := t.AddUserId(c, u.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.Join, error joining tournament for user:%v Error: %v", u.Id, err))
	}

	return nil
}

// Makes a user leave a tournament.
// Todo: should we check that user is indeed a member of the tournament?
func (t *Tournament) Leave(c appengine.Context, u *User) error {
	// find and remove
	if err := u.RemoveTournamentId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.Leave, error leaving tournament for user:%v Error: %v", u.Id, err))
	}
	if err := t.RemoveUserId(c, u.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.Leave, error joining tournament for user:%v Error: %v", u.Id, err))
	}

	return nil
}

// Checks if user is admin of tournament with id 'tournamentId'.
func IsTournamentAdmin(c appengine.Context, tournamentId int64, userId int64) bool {
	if tournament, err := TournamentById(c, tournamentId); err == nil {
		return tournament.AdminId == userId
	}
	return false
}

// Check if a Team has joined the tournament.
func (t *Tournament) TeamJoined(c appengine.Context, team *Team) bool {
	// change in contains
	hasTournament, _ := team.ContainsTournamentId(t.Id)
	return hasTournament
}

// Team joins the Tournament.
func (t *Tournament) TeamJoin(c appengine.Context, team *Team) error {
	// add
	if err := team.AddTournamentId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error adding tournament id to team entity:%v Error: %v", team.Id, err))
	}
	if err := t.AddTeamId(c, team.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error adding team id to tournament entity:%v Error: %v", t.Id, err))
	}
	if err := t.AddUserIds(c, team.UserIds); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error adding user ids to tournament entity:%v Error: %v", t.Id, err))
	}
	if p, errp := CreatePrice(c, team.Id, t.Id, ""); errp != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error creating price for team entity:%v Error: %v", t.Id, errp))
	} else {
		if err := team.AddPriceId(c, p.Id); err != nil {
			return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error adding price id to team entity:%v Error: %v", team.Id, err))
		}
	}
	return nil
}

// Team leaves the Tournament.
func (t *Tournament) TeamLeave(c appengine.Context, team *Team) error {
	// find and remove
	if err := team.RemoveTournamentId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamLeave, error leaving tournament for team:%v Error: %v", team.Id, err))
	}
	if err := t.RemoveTeamId(c, team.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamLeave, error removing team from tournament. For team:%v Error: %v", team.Id, err))
	}
	if err := team.RemovePriceByTournamentId(c, t.Id); err != nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error removing price for team entity:%v Error: %v", team.Id, err))
	}
	return nil
}

// Get the frequency of given word with respect to tournament id.
func GetWordFrequencyForTournament(c appengine.Context, id int64, word string) int64 {

	if tournaments := FindTournaments(c, "Id", id); tournaments != nil {
		return helpers.CountTerm(strings.Split(tournaments[0].KeyName, " "), word)
	}
	return 0
}

// Reset tournament values: Points, GoalsF, GoalsA to zero.
func (t *Tournament) Reset(c appengine.Context) error {
	groups := Groups(c, t.GroupIds)
	for _, g := range groups {
		g.Points = make([]int64, len(g.Teams))
		g.GoalsF = make([]int64, len(g.Teams))
		g.GoalsA = make([]int64, len(g.Teams))
		for _, m := range g.Matches {
			m.Result1 = 0
			m.Result2 = 0
			if err := UpdateMatch(c, &m); err != nil {
				return err
			}
		}
		if err := UpdateGroup(c, g); err != nil {
			return err
		}
	}
	// reset all match rules
	mapMatches2ndRound := MapOf2ndRoundMatches()

	const (
		cMatchId       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	// build matches 2nd phase
	const shortForm = "Jan/02/2006"
	for _, roundMatches := range mapMatches2ndRound {
		for _, matchData := range roundMatches {
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])
			m := GetMatchByIdNumber(c, *t, int64(matchInternalId))
			rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])
			m.Rule = rule
			m.Result1 = 0
			m.Result2 = 0
			if err := UpdateMatch(c, m); err != nil {
				log.Errorf(c, "Reset: unable to reset rule on match: %v", err)
				return err
			}
		}
	}
	return nil
}

// From a tournament returns an array of the users that participate in it.
func (t *Tournament) Participants(c appengine.Context) []*User {
	var users []*User

	for _, uId := range t.UserIds {
		user, err := UserById(c, uId)
		if err != nil {
			log.Errorf(c, " Participants, cannot find user with ID=%", uId)
		} else {
			users = append(users, user)
		}
	}

	return users
}

// from a tournamentid returns an array of teams involved in tournament
func (t *Tournament) Teams(c appengine.Context) []*Team {

	var teams []*Team
	for _, tId := range t.TeamIds {
		team, err := TeamById(c, tId)
		if err != nil {
			log.Errorf(c, " Teams, cannot find team with ID=%", tId)
		} else {
			teams = append(teams, team)
		}
	}
	return teams
}

// Adds a team Id in the TeamId array.
func (t *Tournament) RemoveTeamId(c appengine.Context, tId int64) error {

	if hasTeam, i := t.ContainsTeamId(tId); !hasTeam {
		return errors.New(fmt.Sprintf("RemoveTeamId, not a member."))
	} else {
		// as the order of index in teamsId is not important,
		// replace elem at index i with last element and resize slice.
		t.TeamIds[i] = t.TeamIds[len(t.TeamIds)-1]
		t.TeamIds = t.TeamIds[0 : len(t.TeamIds)-1]
	}
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Adds a team Id in the TeamId array.
func (t *Tournament) AddTeamId(c appengine.Context, tId int64) error {

	if hasTeam, _ := t.ContainsTeamId(tId); hasTeam {
		return errors.New(fmt.Sprintf("AddTeamId, allready a member."))
	}

	t.TeamIds = append(t.TeamIds, tId)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Remove a user Id in the UserId array.
func (t *Tournament) RemoveUserId(c appengine.Context, uId int64) error {

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

// Adds a user Id in the UserId array.
func (t *Tournament) AddUserId(c appengine.Context, uId int64) error {
	if hasUser, _ := t.ContainsUserId(uId); hasUser {
		return errors.New(fmt.Sprintf("AddUserId, allready a member."))
	}

	t.UserIds = append(t.UserIds, uId)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// Add user ids in the tournament entity.
func (t *Tournament) AddUserIds(c appengine.Context, uIds []int64) error {
	for _, uId := range uIds {
		if err := t.AddUserId(c, uId); err != nil {
			log.Errorf(c, "Tournament.AddUserIds, error adding user id to tournament entity: %v", err)
		}
	}
	return nil
}

func (t *Tournament) ContainsTeamId(id int64) (bool, int) {

	for i, tId := range t.TeamIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

func (t *Tournament) ContainsUserId(id int64) (bool, int) {

	for i, tId := range t.UserIds {
		if tId == id {
			return true, i
		}
	}
	return false, -1
}

// Rank users with respect to their score in current tournament.
// Sets the user score to the current tournament score and return array of users
// sorted by that score.
func (t *Tournament) RankingByUser(c appengine.Context, limit int) []*User {
	if limit < 0 {
		return nil
	}
	users := t.Participants(c)
	// set score of user to score of tournament without persisting it.
	for i, u := range users {
		users[i].Score = u.ScoreByTournament(c, t.Id)
	}

	sort.Sort(UserByScore(users))

	if len(users) <= limit {
		return users
	}
	return users[0:limit]
}

func (t *Tournament) RankingByTeam(c appengine.Context, limit int) []*Team {
	if limit < 0 {
		return nil
	}
	teams := t.Teams(c)
	sort.Sort(TeamByAccuracy(teams))
	if len(teams) <= limit {
		return teams
	} else {
		return teams[0:limit]
	}
}

// Publish tournament activity
func (t *Tournament) Publish(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) error {
	var activity Activity
	activity.Type = activityType
	activity.Verb = verb
	activity.Actor = t.Entity("")
	activity.Object = object
	activity.Target = target
	activity.Published = time.Now()
	activity.CreatorID = t.Id

	if err := activity.save(c); err != nil {
		return err
	}
	// add new activity id in user activity table for each participant of the tournament
	for _, p := range t.Participants(c) {
		if err := activity.addNewActivityId(c, p.Id); err != nil {
			log.Errorf(c, "model/tournament, Publish: error occurred during addNewActivityId call: %v", err)
		}
	}

	return nil
}

// Activity entity representation of a tournament
func (t *Tournament) Entity(name string) ActivityEntity {
	displayName := t.Name
	if name != "" {
		displayName = name
	}
	return ActivityEntity{Id: t.Id, Type: "tournament", DisplayName: displayName}
}

func (t *Tournament) Progress(c appengine.Context) float64 {
	const shortForm = "Jan/02/2006"
	// remove this line when going in production.
	now, _ := time.Parse(shortForm, "Jun/28/2014")
	// add this line when going in production.
	//now := time.Now()
	if now.Before(t.Start) {
		return float64(0)
	}
	if now.After(t.End) {
		return float64(1)
	}
	d := t.Start.Sub(now)
	dt := t.Start.Sub(t.End)
	log.Infof(c, "ratio: %v", d.Seconds()/dt.Seconds())
	return d.Seconds() / dt.Seconds()
}
