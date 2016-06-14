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
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
)

// Tournament holds tournament entity data.
//
type Tournament struct {
	Id                   int64
	KeyName              string
	Name                 string
	Description          string
	Start                time.Time
	End                  time.Time
	AdminIds             []int64 // ids of User that are admins of the team
	Created              time.Time
	GroupIds             []int64
	Matches1stStage      []int64
	Matches2ndStage      []int64
	UserIds              []int64
	TeamIds              []int64
	TwoLegged            bool
	IsFirstStageComplete bool
	Official             bool
}

// TournamentJSON is the JSON version of the Tournament struct.
//
type TournamentJSON struct {
	Id                   *int64     `json:",omitempty"`
	KeyName              *string    `json:",omitempty"`
	Name                 *string    `json:",omitempty"`
	Description          *string    `json:",omitempty"`
	Start                *time.Time `json:",omitempty"`
	End                  *time.Time `json:",omitempty"`
	AdminIds             *[]int64   `json:",omitempty"`
	Created              *time.Time `json:",omitempty"`
	GroupIds             *[]int64   `json:",omitempty"`
	Matches1stStage      *[]int64   `json:",omitempty"`
	Matches2ndStage      *[]int64   `json:",omitempty"`
	UserIds              *[]int64   `json:",omitempty"`
	TeamIds              *[]int64   `json:",omitempty"`
	TwoLegged            *bool      `json:",omitempty"`
	IsFirstStageComplete *bool      `json:",omitempty"`
	Official             *bool      `json:",omitempty"`
}

// TournamentBuilder is interface used to build a tournament
//
type TournamentBuilder interface {
	MapOfTeamCodes() map[string]string
	ArrayOfPhases() []string
	MapOfGroups() map[string][]string
	MapOfGroupMatches() map[string][][]string
	MapOf2ndRoundMatches() map[string][][]string
	MapOfPhaseIntervals() map[string][]int64
	MapOfIDTeams(c appengine.Context, tournament *Tournament) map[int64]string
}

// CreateTournament creates tournament entity given a name and description.
//
func CreateTournament(c appengine.Context, name string, description string, start time.Time, end time.Time, adminID int64) (*Tournament, error) {

	tournamentId, _, err := datastore.AllocateIDs(c, "Tournament", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Tournament", "", tournamentId, nil)

	// empty groups and tournaments for now
	var emptyArray []int64
	admins := make([]int64, 1)
	admins[0] = adminID
	twoLegged := false
	official := false

	tournament := &Tournament{tournamentId, helpers.TrimLower(name), name, description, start, end, admins, time.Now(), emptyArray, emptyArray, emptyArray, emptyArray, emptyArray, twoLegged, false, official}

	_, err = datastore.Put(c, key, tournament)
	if err != nil {
		return nil, err
	}

	AddToTournamentInvertedIndex(c, helpers.TrimLower(name), tournamentId)
	return tournament, nil
}

// Destroy a tournament entity given a tournament id.
//
func (t *Tournament) Destroy(c appengine.Context) error {

	if _, err := TournamentByID(c, t.Id); err != nil {
		return fmt.Errorf("Cannot find tournament with Id=%d", t.Id)
	}

	key := datastore.NewKey(c, "Tournament", "", t.Id, nil)
	if errd := datastore.Delete(c, key); errd != nil {
		return errd
	}

	// remove key name.
	return UpdateTournamentInvertedIndex(c, t.KeyName, "", t.Id)
}

// FindTournaments finds all entity tournaments with respect of a filter and value.
//
func FindTournaments(c appengine.Context, filter string, value interface{}) []*Tournament {

	q := datastore.NewQuery("Tournament").Filter(filter+" =", value)
	var tournaments []*Tournament
	if _, err := q.GetAll(c, &tournaments); err != nil {
		log.Errorf(c, " Tournament.Find, error occurred during GetAll: %v", err)
		return nil
	}

	return tournaments
}

// TournamentByID gets a pointer to a tournament given a tournament id.
//
func TournamentByID(c appengine.Context, id int64) (*Tournament, error) {

	var t Tournament
	key := datastore.NewKey(c, "Tournament", "", id, nil)
	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " tournament not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// TournamentKeyByID gets a pointer to a tournament key given a tournament id.
//
func TournamentKeyByID(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Tournament", "", id, nil)
	return key
}

// Update a tournament given a tournament id and a tournament pointer.
//
func (t *Tournament) Update(c appengine.Context) error {

	// update key name
	t.KeyName = helpers.TrimLower(t.Name)
	k := TournamentKeyByID(c, t.Id)
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

// FindAllTournaments finds all tournaments in the datastore.
//
func FindAllTournaments(c appengine.Context, count, page int64) []*Tournament {
	desc := "tournament.FindAllTournaments"
	q := datastore.NewQuery("Tournament")
	var tournaments []*Tournament
	if _, err := q.GetAll(c, &tournaments); err != nil {
		log.Errorf(c, "%s error occurred during GetAll call: %v", desc, err)
	}

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

// TournamentsByIds finds all tournaments with respect to array of ids.
//
func TournamentsByIds(c appengine.Context, ids []int64) ([]*Tournament, error) {

	tournaments := make([]Tournament, len(ids))
	keys := TournamentKeysByIds(c, ids)

	var wrongIndexes []int

	if err := datastore.GetMulti(c, keys, tournaments); err != nil {
		if me, ok := err.(appengine.MultiError); ok {
			for i, merr := range me {
				if merr == datastore.ErrNoSuchEntity {
					log.Errorf(c, "TournamentsByIds, missing key: %v %v", err, keys[i].IntID())

					wrongIndexes = append(wrongIndexes, i)
				}
			}
		} else {
			return nil, err
		}
	}

	var existingTournaments []*Tournament
	for i := range tournaments {
		if !contains(wrongIndexes, i) {
			log.Infof(c, "TournamentsByIds %v", tournaments[i])
			existingTournaments = append(existingTournaments, &tournaments[i])
		}
	}
	return existingTournaments, nil
}

// TournamentKeysByIds finds tournaments keys from an array of ids.
//
func TournamentKeysByIds(c appengine.Context, ids []int64) []*datastore.Key {
	keys := make([]*datastore.Key, len(ids))
	for i, id := range ids {
		keys[i] = TournamentKeyByID(c, id)
	}
	return keys
}

// Joined checks if a user has joined a tournament.
//
func (t *Tournament) Joined(c appengine.Context, u *User) bool {
	// change in contains
	hasTournament, _ := u.ContainsTournamentID(t.Id)
	return hasTournament
}

// Join let a user join a tournament.
//
func (t *Tournament) Join(c appengine.Context, u *User) error {
	// add
	if err := u.AddTournamentID(c, t.Id); err != nil {
		return fmt.Errorf(" Tournament.Join, error joining tournament for user:%v Error: %v", u.Id, err)
	}
	if err := t.AddUserID(c, u.Id); err != nil {
		return fmt.Errorf(" Tournament.Join, error joining tournament for user:%v Error: %v", u.Id, err)
	}

	return nil
}

// IsTournamentAdmin checks if user is admin of tournament with id 'tournamentId'.
//
func IsTournamentAdmin(c appengine.Context, tournamentId int64, userID int64) bool {
	if tournament, err := TournamentByID(c, tournamentId); err == nil {
		for _, aid := range tournament.AdminIds {
			if aid == userID {
				return true
			}
		}
	}
	return false
}

// AddAdmin sdds user to admins of current tournament.
// In order to be an admin of a tournament you should first be a member of the tournament.
//
func (t *Tournament) AddAdmin(c appengine.Context, id int64) error {

	if ismember, _ := t.ContainsUserID(id); ismember {
		if isadmin, _ := t.ContainsAdminID(id); isadmin {
			return fmt.Errorf("User with %v is already an admin of tournament", id)
		}
		t.AdminIds = append(t.AdminIds, id)
		if err := t.Update(c); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("User with %v is not a member of the tournament", id)
}

// RemoveAdmin removes user of admins array in current tournament.
// In order to remove an admin from a team, there should be at least an admin in the array.
//
func (t *Tournament) RemoveAdmin(c appengine.Context, id int64) error {

	if ismember, _ := t.ContainsUserID(id); !ismember {
		return fmt.Errorf("User with %v is not a member of the tournament", id)
	}

	isAdmin := false
	i := 0

	if isAdmin, i = t.ContainsAdminID(id); !isAdmin {
		return fmt.Errorf("User with %v is not admin of the tournament", id)
	}

	if len(t.AdminIds) <= 1 {
		return fmt.Errorf("Cannot remove admin %v as there are no admins left in tournament", id)
	}

	// as the order of index in adminIds is not important,
	// replace elem at index i with last element and resize slice.
	t.AdminIds[i] = t.AdminIds[len(t.AdminIds)-1]
	t.AdminIds = t.AdminIds[0 : len(t.AdminIds)-1]
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// ContainsAdminID checks if user is admin of team.
//
func (t *Tournament) ContainsAdminID(id int64) (bool, int) {

	for i, aID := range t.AdminIds {
		if aID == id {
			return true, i
		}
	}
	return false, -1
}

// TeamJoined checks if a Team has joined the tournament.
//
func (t *Tournament) TeamJoined(c appengine.Context, team *Team) bool {
	// change in contains
	hasTournament, _ := team.ContainsTournamentID(t.Id)
	return hasTournament
}

// TeamJoin joins a teams to the Tournament.
//
func (t *Tournament) TeamJoin(c appengine.Context, team *Team) error {
	var err error

	// add
	if err = team.AddTournamentID(c, t.Id); err != nil {
		return fmt.Errorf(" Tournament.TeamJoin, error adding tournament id to team entity:%d Error: %v", team.Id, err)
	}
	if err = t.AddTeamID(c, team.Id); err != nil {
		return fmt.Errorf(" Tournament.TeamJoin, error adding team id to tournament entity:%d Error: %v", t.Id, err)
	}
	if err = t.AddUserIDs(c, team.UserIds); err != nil {
		return fmt.Errorf(" Tournament.TeamJoin, error adding user ids to tournament entity:%d Error: %v", t.Id, err)
	}

	var price *Price
	if price, err = CreatePrice(c, team.Id, t.Id, t.Name, ""); err != nil {
		return fmt.Errorf(" Tournament.TeamJoin, error creating price for team entity:%d Error: %v", t.Id, err)
	}

	if err = team.AddPriceID(c, price.Id); err != nil {
		return fmt.Errorf(" Tournament.TeamJoin, error adding price id to team entity:%d Error: %v", team.Id, err)
	}

	return nil
}

// TeamLeave leaves a team from the Tournament.
//
func (t *Tournament) TeamLeave(c appengine.Context, team *Team) error {
	// find and remove
	if err := team.RemoveTournamentID(c, t.Id); err != nil {
		return fmt.Errorf(" Tournament.TeamLeave, error leaving tournament for team:%v Error: %v", team.Id, err)
	}
	if err := t.RemoveTeamID(c, team.Id); err != nil {
		return fmt.Errorf(" Tournament.TeamLeave, error removing team from tournament. For team:%v Error: %v", team.Id, err)
	}
	if err := team.RemovePriceByTournamentID(c, t.Id); err != nil {
		return fmt.Errorf(" Tournament.TeamJoin, error removing price for team entity:%v Error: %v", team.Id, err)
	}
	return nil
}

// GetWordFrequencyForTournament gets the frequency of given word with respect to tournament id.
//
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
	var tb TournamentBuilder
	if tb = GetTournamentBuilder(t); tb == nil {
		log.Errorf(c, "TournamentBuilder not found")
		return fmt.Errorf("TournamentBuilder not found")
	}

	mapMatches2ndRound := tb.MapOf2ndRoundMatches()

	const (
		cMatchID       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	// build matches 2nd phase
	const shortForm = "Jan/02/2006"
	for _, roundMatches := range mapMatches2ndRound {
		for _, matchData := range roundMatches {
			matchInternalID, _ := strconv.Atoi(matchData[cMatchID])
			m := GetMatchByIDNumber(c, *t, int64(matchInternalID))
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

// Participants returns an array, from the tournament, of the users that participate in it.
//
func (t *Tournament) Participants(c appengine.Context) []*User {
	var users []*User

	for _, uID := range t.UserIds {
		user, err := UserByID(c, uID)
		if err != nil {
			log.Errorf(c, " Participants, cannot find user with Id=%d", uID)
		} else {
			users = append(users, user)
		}
	}

	return users
}

// Teams returns an array of teams involved in tournament, from a tournament.
//
func (t *Tournament) Teams(c appengine.Context) []*Team {

	var teams []*Team
	for _, tID := range t.TeamIds {
		team, err := TeamByID(c, tID)
		if err != nil {
			log.Errorf(c, " Teams, cannot find team with Id=%", tID)
		} else {
			teams = append(teams, team)
		}
	}
	return teams
}

// RemoveTeamID adds a team Id in the TeamId array.
//
func (t *Tournament) RemoveTeamID(c appengine.Context, tID int64) error {

	hasTeam := false
	i := 0

	if hasTeam, i = t.ContainsTeamID(tID); !hasTeam {
		return fmt.Errorf("RemoveTeamID, not a member")
	}

	// as the order of index in teamsId is not important,
	// replace elem at index i with last element and resize slice.
	t.TeamIds[i] = t.TeamIds[len(t.TeamIds)-1]
	t.TeamIds = t.TeamIds[0 : len(t.TeamIds)-1]

	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// AddTeamID adds a team Id in the TeamId array.
//
func (t *Tournament) AddTeamID(c appengine.Context, tID int64) error {

	if hasTeam, _ := t.ContainsTeamID(tID); hasTeam {
		return fmt.Errorf("AddTeamID, allready a member")
	}

	t.TeamIds = append(t.TeamIds, tID)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// RemoveUserID removes a user Id in the UserId array.
//
func (t *Tournament) RemoveUserID(c appengine.Context, uID int64) error {

	hasUser := false
	i := 0

	if hasUser, i = t.ContainsUserID(uID); !hasUser {
		return fmt.Errorf("RemoveUserID, not a member")
	}

	// as the order of index in usersId is not important,
	// replace elem at index i with last element and resize slice.
	t.UserIds[i] = t.UserIds[len(t.UserIds)-1]
	t.UserIds = t.UserIds[0 : len(t.UserIds)-1]

	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// AddUserID adds a user Id in the UserId array.
//
func (t *Tournament) AddUserID(c appengine.Context, uID int64) error {
	if hasUser, _ := t.ContainsUserID(uID); hasUser {
		return fmt.Errorf("AddUserID, allready a member")
	}

	t.UserIds = append(t.UserIds, uID)
	if err := t.Update(c); err != nil {
		return err
	}
	return nil
}

// AddUserIDs adds user ids in the tournament entity.
//
func (t *Tournament) AddUserIDs(c appengine.Context, uIds []int64) error {
	for _, uID := range uIds {
		if err := t.AddUserID(c, uID); err != nil {
			log.Errorf(c, "Tournament.AddUserIDs, error adding user id to tournament entity: %v", err)
		}
	}
	return nil
}

// ContainsTeamID checks if a team is part of a tournament.
//
func (t *Tournament) ContainsTeamID(id int64) (bool, int) {

	for i, tID := range t.TeamIds {
		if tID == id {
			return true, i
		}
	}
	return false, -1
}

// ContainsUserID checks if user is part of the tournament.
//
func (t *Tournament) ContainsUserID(id int64) (bool, int) {

	for i, tID := range t.UserIds {
		if tID == id {
			return true, i
		}
	}
	return false, -1
}

// RankingByUser ranks users with respect to their score in current tournament.
// Sets the user score to the current tournament score and return array of users
// sorted by that score.
//
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
	return users[len(users)-limit:]
}

// RankingByTeam ranks teams with respect ot their accuracy in the current tournament.
//
func (t *Tournament) RankingByTeam(c appengine.Context, limit int) []*Team {
	if limit < 0 {
		return nil
	}
	teams := t.Teams(c)
	sort.Sort(TeamByAccuracy(teams))
	if len(teams) <= limit {
		return teams
	}

	return teams[len(teams)-limit:]
}

// Publish tournament activity.
//
func (t *Tournament) Publish(c appengine.Context, activityType string, verb string, object ActivityEntity, target ActivityEntity) error {
	var activity Activity
	activity.Type = activityType
	activity.Verb = verb
	activity.Actor = t.Entity()
	activity.Object = object
	activity.Target = target
	activity.Published = time.Now()
	activity.CreatorID = t.Id

	if err := activity.save(c); err != nil {
		return err
	}
	// add new activity id in user activity table for each participant of the tournament
	for _, p := range t.Participants(c) {
		if err := activity.AddNewActivityID(c, p); err != nil {
			log.Errorf(c, "model/tournament, Publish: error occurred during AddNewActivityID call: %v", err)
		} else {
			if err1 := p.Update(c); err1 != nil {
				log.Errorf(c, "model/tournament, Publish: error occurred during update call: %v", err1)
			}

		}
	}

	return nil
}

// Entity is the Activity entity representation of a tournament.
//
func (t *Tournament) Entity() ActivityEntity {
	return ActivityEntity{Id: t.Id, Type: "tournament", DisplayName: t.Name}
}

// Progress is the progression of the tournament.
// The progression is a number between 0 and 1 with the progression of the tournament
// with respect of today's date and start and end date of tournament.
//
func (t *Tournament) Progress(c appengine.Context) float64 {
	now := time.Now()
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

// GetTournamentBuilder gets the tournament builder of a given tournament.
//
func GetTournamentBuilder(t *Tournament) TournamentBuilder {
	var tb TournamentBuilder
	if t.Name == "2014 FIFA World Cup" {
		wct := WorldCupTournament{}
		tb = wct
	} else if t.Name == "2015-2016 UEFA Champions League" {
		clt := ChampionsLeagueTournament{}
		tb = clt
	} else if t.Name == "2015-2016 UEFA Champions League" {
		clt := ChampionsLeagueTournament20152016{}
		tb = clt
	} else if t.Name == "2015 Copa America" {
		cat := CopaAmericaTournament{}
		tb = cat
	} else if t.Name == "2016 UEFA Euro" {
		et := EuroTournament2016{}
		tb = et
	} else if t.Name == "2016 Copa America" {
		cat := CopaAmericaTournament2016{}
		tb = cat
	}
	return tb
}

// MapOfIDTeams is the map of team IDs of a given tournament.
//
func MapOfIDTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	var tb TournamentBuilder

	if tb = GetTournamentBuilder(tournament); tb == nil {
		return nil
	}
	return tb.MapOfIDTeams(c, tournament)
}
