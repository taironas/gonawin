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

package tournament

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
	tournamentinvidmdl "github.com/santiaago/purple-wing/models/tournamentInvertedIndex"
	tournamentrelmdl "github.com/santiaago/purple-wing/models/tournamentrel"
	tournamentteamrelmdl "github.com/santiaago/purple-wing/models/tournamentteamrel"
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
}

type Tgroup struct {
	Id      int64
	Name    string
	Teams   []Tteam
	Matches []Tmatch
	Points  []int64
	GoalsF  []int64
	GoalsA  []int64
}

type Tteam struct {
	Id   int64
	Name string
}

type Tmatch struct {
	Id       int64
	IdNumber int64
	Date     time.Time
	TeamId1  int64
	TeamId2  int64
	Location string
	Rule     string // we use this field to store a specific match rule.
	Result1  int64
	Result2  int64
}

type Tday struct {
	Date    time.Time
	Matches []Tmatch
}

type Tphase struct {
	Name string
	Days []Tday
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
}

// create tournament entity given a name and description
func Create(c appengine.Context, name string, description string, start time.Time, end time.Time, adminId int64) (*Tournament, error) {
	// create new tournament
	tournamentID, _, err := datastore.AllocateIDs(c, "Tournament", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Tournament", "", tournamentID, nil)

	// empty groups and tournaments for now
	groupIds := make([]int64, 0)
	matches1stStageIds := make([]int64, 0)
	matches2ndStageIds := make([]int64, 0)

	tournament := &Tournament{tournamentID, helpers.TrimLower(name), name, description, start, end, adminId, time.Now(), groupIds, matches1stStageIds, matches2ndStageIds}

	_, err = datastore.Put(c, key, tournament)
	if err != nil {
		return nil, err
	}

	tournamentinvidmdl.Add(c, helpers.TrimLower(name), tournamentID)

	return tournament, nil
}

// destroy a tournament given a tournament id
func Destroy(c appengine.Context, tournamentId int64) error {

	if tournament, err := ById(c, tournamentId); err != nil {
		return errors.New(fmt.Sprintf("Cannot find tournament with tournamentId=%d", tournamentId))
	} else {
		key := datastore.NewKey(c, "Tournament", "", tournament.Id, nil)

		return datastore.Delete(c, key)
	}
}

// return an array of tournaments given a filter and value
func Find(c appengine.Context, filter string, value interface{}) []*Tournament {

	q := datastore.NewQuery("Tournament").Filter(filter+" =", value)

	var tournaments []*Tournament

	if _, err := q.GetAll(c, &tournaments); err == nil {
		return tournaments
	} else {
		log.Errorf(c, " Tournament.Find, error occurred during GetAll: %v", err)
		return nil
	}
}

// returns a pointer to a tournament given a tournament id
func ById(c appengine.Context, id int64) (*Tournament, error) {

	var t Tournament
	key := datastore.NewKey(c, "Tournament", "", id, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, " tournament not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// return a pointer to a tournament key given a tournament id
func KeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Tournament", "", id, nil)

	return key
}

// Update a tournament given a tournament id and a tournament pointer.
func Update(c appengine.Context, id int64, t *Tournament) error {
	// update key name
	t.KeyName = helpers.TrimLower(t.Name)
	k := KeyById(c, id)
	oldTournament := new(Tournament)
	if err := datastore.Get(c, k, oldTournament); err == nil {
		if _, err = datastore.Put(c, k, t); err != nil {
			return err
		}
		tournamentinvidmdl.Update(c, oldTournament.Name, t.Name, id)
	}
	return nil
}

// returns an array of all tournaments in the datastore
func FindAll(c appengine.Context) []*Tournament {

	q := datastore.NewQuery("Tournament")

	var tournaments []*Tournament

	if _, err := q.GetAll(c, &tournaments); err != nil {
		log.Errorf(c, " Tournament.FindAll, error occurred during GetAll call: %v", err)
	}

	return tournaments
}

// find with respect to array of ids
func ByIds(c appengine.Context, ids []int64) []*Tournament {

	var tournaments []*Tournament
	for _, id := range ids {
		if tournament, err := ById(c, id); err == nil {
			tournaments = append(tournaments, tournament)
		} else {
			log.Errorf(c, " Tournament.ByIds, error occurred during ByIds call: %v", err)
		}
	}
	return tournaments
}

// checks if a user has joined a tournament
func Joined(c appengine.Context, tournamentId int64, userId int64) bool {
	tournamentRel := tournamentrelmdl.FindByTournamentIdAndUserId(c, tournamentId, userId)
	return tournamentRel != nil
}

// makes a user join a tournament
func Join(c appengine.Context, tournamentId int64, userId int64) error {
	if tournamentRel, err := tournamentrelmdl.Create(c, tournamentId, userId); tournamentRel == nil {
		return errors.New(fmt.Sprintf(" Tournament.Join, error during tournament relationship creation: %v", err))
	}

	return nil
}

// makes a user leave a tournament.
// Todo: should we check that user is indeed a member of the tournament?
func Leave(c appengine.Context, tournamentId int64, userId int64) error {
	return tournamentrelmdl.Destroy(c, tournamentId, userId)
}

// checks if user is admin of given tournament
func IsTournamentAdmin(c appengine.Context, tournamentId int64, userId int64) bool {
	if tournament, err := ById(c, tournamentId); err == nil {
		return tournament.AdminId == userId
	}

	return false
}

// Check if a Team has joined the tournament
func TeamJoined(c appengine.Context, tournamentId int64, teamId int64) bool {
	tournamentteamRel := tournamentteamrelmdl.FindByTournamentIdAndTeamId(c, tournamentId, teamId)
	return tournamentteamRel != nil
}

// Team joins the Tournament
func TeamJoin(c appengine.Context, tournamentId int64, teamId int64) error {
	if tournamentteamRel, err := tournamentteamrelmdl.Create(c, tournamentId, teamId); tournamentteamRel == nil {
		return errors.New(fmt.Sprintf(" Tournament.TeamJoin, error during tournament team relationship creation: %v", err))
	}

	return nil
}

// Team leaves the Tournament
func TeamLeave(c appengine.Context, tournamentId int64, teamId int64) error {
	return tournamentteamrelmdl.Destroy(c, tournamentId, teamId)
}

// get the frequency of given word with respect to tournament id
func GetWordFrequencyForTournament(c appengine.Context, id int64, word string) int64 {

	if tournaments := Find(c, "Id", id); tournaments != nil {
		return helpers.CountTerm(strings.Split(tournaments[0].KeyName, " "), word)
	}
	return 0
}

// create world cup tournament entity
func CreateWorldCup(c appengine.Context, adminId int64) (*Tournament, error) {
	// create new tournament
	tournamentID, _, err := datastore.AllocateIDs(c, "Tournament", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "Tournament", "", tournamentID, nil)

	log.Infof(c, "World Cup: start")

	// build map of groups
	var mapWCGroups map[string][]string
	mapWCGroups = MapOfGroups()

	// build map of matches
	var mapTeamId map[string]int64
	mapTeamId = make(map[string]int64)

	// mapGroupMatches is a map where the key is a string which represent the group
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// map[string][][]string
	// example: {"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}
	mapGroupMatches := MapOfGroupMatches()
	const (
		cMatchId       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	// matches1stStageIds is an array of  int64
	// where we allocate IDs of the Tmatches entities
	// we will store them in the tournament entity for easy retreival later on.
	matchesA := mapGroupMatches["A"]
	matches1stStageIds := make([]int64, 8*len(matchesA))

	log.Infof(c, "World Cup: maps ready")

	// build groups, and teams
	groups := make([]Tgroup, len(mapWCGroups))
	for groupName, teams := range mapWCGroups {
		log.Infof(c, "---------------------------------------")
		log.Infof(c, "World Cup: working with group: %v", groupName)
		log.Infof(c, "World Cup: teams: %v", teams)

		var group Tgroup
		group.Name = groupName
		groupIndex := int64(groupName[0]) - 65
		group.Teams = make([]Tteam, len(teams))
		group.Points = make([]int64, len(teams))
		group.GoalsF = make([]int64, len(teams))
		group.GoalsA = make([]int64, len(teams))
		for i, teamName := range teams {
			log.Infof(c, "World Cup: team: %v", teamName)

			teamID, _, err := datastore.AllocateIDs(c, "Tteam", nil, 1)
			log.Infof(c, "World Cup: team: %v allocateIDs ok", teamName)

			teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
			log.Infof(c, "World Cup: team: %v NewKey ok", teamName)

			team := &Tteam{teamID, teamName}
			log.Infof(c, "World Cup: team: %v instance of team ok", teamName)

			_, err = datastore.Put(c, teamkey, team)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "World Cup: team: %v put in datastore ok", teamName)
			group.Teams[i] = *team
			mapTeamId[teamName] = teamID
		}

		// build group matches:
		log.Infof(c, "World Cup: building group matches")

		// for date parsing
		const shortForm = "Jan/02/2006"

		groupMatches := mapGroupMatches[groupName]
		group.Matches = make([]Tmatch, len(groupMatches))

		for matchIndex, matchData := range groupMatches {
			log.Infof(c, "World Cup: match data: %v", matchData)

			matchID, _, err := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			log.Infof(c, "World Cup: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "World Cup: match: new key ok")

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])
			emptyrule := ""
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalId),
				matchTime,
				mapTeamId[matchData[cMatchTeam1]],
				mapTeamId[matchData[cMatchTeam2]],
				matchData[cMatchLocation],
				emptyrule,
				emptyresult,
				emptyresult,
			}
			log.Infof(c, "World Cup: match: build match ok")

			_, err = datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "World Cup: match: %v put in datastore ok", matchData)
			group.Matches[matchIndex] = *match

			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches1stStageIds[int64(matchInternalId)-1] = matchID
		}

		groupID, _, err := datastore.AllocateIDs(c, "Tgroup", nil, 1)
		log.Infof(c, "World Cup: Group: %v allocate ID ok", groupName)

		groupkey := datastore.NewKey(c, "Tgroup", "", groupID, nil)
		log.Infof(c, "World Cup: Group: %v New Key ok", groupName)

		group.Id = groupID
		groups[groupIndex] = group
		_, err = datastore.Put(c, groupkey, &group)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "World Cup: Group: %v put in datastore ok", groupName)
	}

	// build array of group ids
	groupIds := make([]int64, 8)
	for i, _ := range groupIds {
		groupIds[i] = groups[i].Id
	}

	log.Infof(c, "World Cup: build of groups ids complete: %v", groupIds)

	// matches 2nd stage
	matches2ndStageIds := make([]int64, 0)

	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: {"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}
	mapMatches2ndRound := MapOf2ndRoundMatches()

	// build matches 2nd phase
	const shortForm = "Jan/02/2006"
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "World Cup: building 2nd round matches: round number %v", roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "World Cup: second phase match data: %v", matchData)

			matchID, _, err := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			log.Infof(c, "World Cup: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "World Cup: match: new key ok")

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])

			rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])
			emptyresult := int64(0)
			match := &Tmatch{
				matchID,
				int64(matchInternalId),
				matchTime,
				0, // second round matches start with ids at 0
				0, // second round matches start with ids at 0
				matchData[cMatchLocation],
				rule,
				emptyresult,
				emptyresult,
			}
			log.Infof(c, "World Cup: match 2nd round: build match ok")

			_, err = datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "World Cup: 2nd round match: %v put in datastore ok", matchData)
			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches2ndStageIds = append(matches2ndStageIds, matchID)
		}
	}

	tournament := &Tournament{
		tournamentID,
		helpers.TrimLower("world cup"),
		"World Cup",
		"FIFA World Cup",
		time.Now(),
		time.Now(),
		adminId,
		time.Now(),
		groupIds,
		matches1stStageIds,
		matches2ndStageIds,
	}
	log.Infof(c, "World Cup: instance of tournament ready")

	_, err = datastore.Put(c, key, tournament)
	if err != nil {
		return nil, err
	}
	log.Infof(c, "World Cup:  tournament put in datastore ok")

	tournamentinvidmdl.Add(c, helpers.TrimLower("world cup"), tournamentID)

	return tournament, nil
}

// Get a Tgroup entity by id.
func GroupById(c appengine.Context, groupId int64) (*Tgroup, error) {
	var g Tgroup
	key := datastore.NewKey(c, "Tgroup", "", groupId, nil)

	if err := datastore.Get(c, key, &g); err != nil {
		log.Errorf(c, "group not found : %v", err)
		return &g, err
	}
	return &g, nil
}

// Get a Tmatch entity by id.
func MatchById(c appengine.Context, matchId int64) (*Tmatch, error) {
	var m Tmatch
	key := datastore.NewKey(c, "Tmatch", "", matchId, nil)

	if err := datastore.Get(c, key, &m); err != nil {
		log.Errorf(c, "match not found : %v", err)
		return &m, err
	}
	return &m, nil
}

// Get a Tteam entity by id.
func TeamById(c appengine.Context, teamId int64) (*Tteam, error) {
	var t Tteam
	key := datastore.NewKey(c, "Tteam", "", teamId, nil)

	if err := datastore.Get(c, key, &t); err != nil {
		log.Errorf(c, "team not found : %v", err)
		return &t, err
	}
	return &t, nil
}

// From an array of groups id return array of Tgroups.
func Groups(c appengine.Context, groupIds []int64) []*Tgroup {

	var groups []*Tgroup

	for _, groupId := range groupIds {

		g, err := GroupById(c, groupId)
		if err != nil {
			log.Errorf(c, " Groups, cannot find group with ID=%", groupId)
		} else {
			groups = append(groups, g)
		}
	}
	return groups
}

// From an array of ids return the corresponding array of matches.
func Matches(c appengine.Context, matchIds []int64) []*Tmatch {

	var matches []*Tmatch

	for _, matchId := range matchIds {

		m, err := MatchById(c, matchId)
		if err != nil {
			log.Errorf(c, " Matches, cannot find match with ID=%", matchId)
		} else {
			matches = append(matches, m)
		}
	}
	return matches

}

func GetMatchByIdNumber(c appengine.Context, tournament Tournament, matchInternalId int64) *Tmatch {
	matches1stStage := Matches(c, tournament.Matches1stStage)
	for _, m := range matches1stStage {
		if m.IdNumber == matchInternalId {
			return m
		}
	}
	matches2ndStage := Matches(c, tournament.Matches2ndStage)
	for _, m := range matches2ndStage {
		if m.IdNumber == matchInternalId {
			return m
		}
	}
	return nil
}

// From tournament entity build map of teams.
func MapOfIdTeams(c appengine.Context, tournament Tournament) map[int64]string {

	var mapIdTeams map[int64]string
	mapIdTeams = make(map[int64]string)

	groups := Groups(c, tournament.GroupIds)
	for _, g := range groups {
		for _, t := range g.Teams {
			mapIdTeams[t.Id] = t.Name
		}
	}
	return mapIdTeams
}

// return a pointer to a match key given a match id
func KeyByIdMatch(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Tmatch", "", id, nil)

	return key
}

// return a pointer to a group key given a group id
func KeyByIdGroup(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "Tgroup", "", id, nil)

	return key
}

// Update a match given a match pointer
func UpdateMatch(c appengine.Context, m *Tmatch) error {

	k := KeyByIdMatch(c, m.Id)
	oldMatch := new(Tmatch)
	if err := datastore.Get(c, k, oldMatch); err == nil {
		if _, err = datastore.Put(c, k, m); err != nil {
			return err
		}
	}
	return nil
}

// Update an array of matches
func UpdateMatches(c appengine.Context, matches []*Tmatch) error {
	keys := make([]*datastore.Key, len(matches))
	for i, _ := range keys {
		keys[i] = KeyByIdMatch(c, matches[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, matches); err != nil {
		return err
	}
	return nil
}

// Update an array of groups
func UpdateGroups(c appengine.Context, groups []*Tgroup) error {
	keys := make([]*datastore.Key, len(groups))
	for i, _ := range keys {
		keys[i] = KeyByIdGroup(c, groups[i].Id)
	}
	if _, err := datastore.PutMulti(c, keys, groups); err != nil {
		return err
	}
	return nil
}

// Update a group given an a group pointer
func UpdateGroup(c appengine.Context, g *Tgroup) error {
	k := KeyByIdGroup(c, g.Id)
	oldGroup := new(Tgroup)
	if err := datastore.Get(c, k, oldGroup); err == nil {
		if _, err = datastore.Put(c, k, g); err != nil {
			return err
		}
	}
	return nil
}

// Set results in an array of matches and triggers a match update and group update.
func SetResults(c appengine.Context, matches []*Tmatch, results1 []int64, results2 []int64, t *Tournament) error {
	log.Infof(c, "Set Results: begin")
	if len(matches) != len(results1) || len(matches) != len(results2) {
		log.Errorf(c, "Set Result: unable to set result on matches")
		return errors.New(helpers.ErrorCodeMatchesCannotUpdate)
	}

	for i, m := range matches {
		log.Infof(c, "Tournament Set Results: current match: %v", m.Id)
		if results1[i] < 0 || results2[i] < 0 {
			log.Errorf(c, "Set Result: unable to set result on match with id: %v, %v", m.Id)
			return errors.New(helpers.ErrorCodeMatchCannotUpdate)
		}
		m.Result1 = results1[i]
		m.Result2 = results2[i]
	}

	log.Infof(c, "Set Results: matches ready")
	// batch match update
	if err := UpdateMatches(c, matches); err != nil {
		log.Errorf(c, "Set Results: unable to set results on matches: %v", err)
		return err
	}
	log.Infof(c, "Set Results: matches updated")
	allMatches := GetAllMatchesFromTournament(c, t)
	phases := MatchesGroupByPhase(allMatches)

	for _, m := range matches {
		log.Infof(c, "Tournament Set Results: Trigger current match: %v", m.Id)

		if ismatch, g := IsMatchInGroup(c, t, m); ismatch == true {
			if err := UpdatePointsAndGoals(c, g, m, t); err != nil {
				log.Errorf(c, "Update Points and Goals: unable to update points and goals for group for match with id:%v error: %v", m.IdNumber, err)
				return errors.New(helpers.ErrorCodeMatchCannotUpdate)
			}
			if err := UpdateGroup(c, g); err != nil {
				log.Errorf(c, "Set Results: unable to update group: %v", err)
				return err
			}
		}
		if isLast, phaseId := lastMatchOfPhase(c, m, &phases); isLast == true {
			log.Infof(c, "Tournament Set Results: -------------------------------------------------->")
			log.Infof(c, "Tournament Set Results: Trigger update of next phase here: next phase: %v", phaseId+1)
			log.Infof(c, "Tournament Set Results: Trigger update of next phase here: next phase: %v", m)
			if int(phaseId+1) < len(phases) {
				UpdateNextPhase(c, t, &phases[phaseId], &phases[phaseId+1])
			}
			log.Infof(c, "Tournament Set Results: -------------------------------------------------->")
		}
	}

	log.Infof(c, "Set Results: points and goals updated")

	return nil
}

// Set result in match entity and triggers a match update in datastore.
func SetResult(c appengine.Context, m *Tmatch, result1 int64, result2 int64, t *Tournament) error {

	if result1 < 0 || result2 < 0 {
		log.Errorf(c, "Set Result: unable to set result on match with id: %v, %v", m.Id)
		return errors.New(helpers.ErrorCodeMatchCannotUpdate)
	}
	m.Result1 = result1
	m.Result2 = result2

	if err := UpdateMatch(c, m); err != nil {
		log.Errorf(c, "Set Result: unable to set result on match with id: %v, %v", m.Id, err)
		return err
	}
	if ismatch, g := IsMatchInGroup(c, t, m); ismatch == true {
		if err := UpdatePointsAndGoals(c, g, m, t); err != nil {
			log.Errorf(c, "Update Points and Goals: unable to update points and goals for group for match with id:%v error: %v", m.IdNumber, err)
			return errors.New(helpers.ErrorCodeMatchCannotUpdate)
		}
		UpdateGroup(c, g)
	}
	allMatches := GetAllMatchesFromTournament(c, t)
	phases := MatchesGroupByPhase(allMatches)
	if isLast, phaseId := lastMatchOfPhase(c, m, &phases); isLast == true {
		log.Infof(c, "Tournament Set Result: -------------------------------------------------->")
		log.Infof(c, "Tournament Set Result: Trigger update of next phase here: next phase: %v", phaseId+1)
		log.Infof(c, "Tournament Set Result: Trigger update of next phase here: next phase: %v", m)
		if int(phaseId+1) < len(phases) {
			UpdateNextPhase(c, t, &phases[phaseId], &phases[phaseId+1])
		}
		log.Infof(c, "Tournament Set Results: -------------------------------------------------->")
	}

	return nil
}

// Check if the match is part of a group phase in current tournament.
func IsMatchInGroup(c appengine.Context, t *Tournament, m *Tmatch) (bool, *Tgroup) {
	groups := Groups(c, t.GroupIds)
	for i, g := range groups {
		for _, match := range g.Matches {
			if m.Id == match.Id {
				return true, groups[i]
			}
		}
	}
	return false, nil
}

// Update points in group with result of match
func UpdatePointsAndGoals(c appengine.Context, g *Tgroup, m *Tmatch, tournament *Tournament) error {
	for i, t := range g.Teams {
		if t.Id == m.TeamId1 {
			if m.Result1 > m.Result2 {
				g.Points[i] += 3
			} else if m.Result1 == m.Result2 {
				g.Points[i] += 1
			}
			g.GoalsF[i] += m.Result1
			g.GoalsA[i] += m.Result2
		} else if t.Id == m.TeamId2 {
			if m.Result2 > m.Result1 {
				g.Points[i] += 3
			} else if m.Result2 == m.Result1 {
				g.Points[i] += 1
			}
			g.GoalsF[i] += m.Result2
			g.GoalsA[i] += m.Result1
		}
	}
	return nil
}

// Reset tournament values: Points, GoalsF, GoalsA to zero.
func Reset(c appengine.Context, t *Tournament) error {
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

	// reset all matches rules
	return nil
}

// Check if the match m passed as argument is the last match of a phase in a specific tournament.
// it returns a boolean and the index of the phase the match was found
func lastMatchOfPhase(c appengine.Context, m *Tmatch, phases *[]Tphase) (bool, int64) {

	for i, ph := range *phases {
		if n := len(ph.Days); n >= 1 {
			lastDay := ph.Days[n-1]
			if n = len(lastDay.Matches); n >= 1 {
				lastMatch := lastDay.Matches[n-1]
				if lastMatch.IdNumber == m.IdNumber {
					return true, int64(i)
				}
			}
		}
	}
	return false, int64(-1)
}

// Update next phase in tournament
func UpdateNextPhase(c appengine.Context, t *Tournament, currentphase *Tphase, nextphase *Tphase) error {

	// the array of phases that will be update.
	// it is an array as a phase can trigger an update in multiple phases, like semi-finals
	// trigger update of Third place and Finals
	var phases []*Tphase
	phases = append(phases, nextphase)
	// compute ranking of previous phase
	var mapOfTeams map[string]*Tteam
	mapOfTeams = make(map[string]*Tteam)

	if currentphase.Name == cFirstStage {
		// compute ranking of groups
		// get all groups.
		groups := Groups(c, t.GroupIds)
		for _, g := range groups {
			team1, argTeam1 := getFirstTeamInGroup(c, g)
			team2, _ := getSecondTeamInGroup(c, g, argTeam1)
			mapOfTeams["1"+g.Name] = team1
			mapOfTeams["2"+g.Name] = team2
		}
	} else {
		// compute ranking just by match winners
		if currentphase.Name == cFinals || currentphase.Name == cThirdPlace {
			// nothing to do.
			return nil
		}

		currentmatches := getMatchesByPhase(c, t, currentphase.Name)
		if currentphase.Name != cSemiFinals {
			log.Infof(c, "Update Next phase: current %v", currentphase.Name)
			log.Infof(c, "Update Next phase: next %v", nextphase.Name)

			for _, m := range currentmatches {
				// ToDo: handle penalties
				if m.Result1 >= m.Result2 {
					team1, _ := TeamById(c, m.TeamId1)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team1
					log.Infof(c, "Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team1.Name)

				} else if m.Result1 < m.Result2 {
					team2, _ := TeamById(c, m.TeamId2)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team2
					log.Infof(c, "Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team2.Name)
				}
			}
		} else {
			// append finals phases to array of phases to update.
			var finals Tphase
			finals.Name = cFinals
			phases = append(phases, &finals)

			for _, m := range currentmatches {
				// ToDo: handle penalties
				if m.Result1 >= m.Result2 {
					team1, _ := TeamById(c, m.TeamId1)
					team2, _ := TeamById(c, m.TeamId2)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team1
					mapOfTeams["L"+strconv.Itoa(int(m.IdNumber))] = team2
					log.Infof(c, "Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team1.Name)
					log.Infof(c, "Update Next phase: rule: L%v teams: %v", strconv.Itoa(int(m.IdNumber)), team2.Name)

				} else if m.Result1 < m.Result2 {
					team2, _ := TeamById(c, m.TeamId2)
					team1, _ := TeamById(c, m.TeamId1)
					mapOfTeams["W"+strconv.Itoa(int(m.IdNumber))] = team2
					mapOfTeams["L"+strconv.Itoa(int(m.IdNumber))] = team1
					log.Infof(c, "Update Next phase: rule: W%v teams: %v", strconv.Itoa(int(m.IdNumber)), team2.Name)
					log.Infof(c, "Update Next phase: rule: L%v teams: %v", strconv.Itoa(int(m.IdNumber)), team1.Name)
				}
			}

		}
	}

	// update phase (matches) with new teams
	for _, ph := range phases {
		matches := getMatchesByPhase(c, t, ph.Name)
		for i, m := range matches {
			log.Infof(c, "Update Next phase: current rule: %v", m.Rule)
			rule := strings.Split(m.Rule, " ")
			if len(rule) != 2 {
				continue
			}

			if val, ok := mapOfTeams[rule[0]]; ok {
				log.Infof(c, "Update Next phase: match found: %v", val.Name)
				matches[i].TeamId1 = val.Id
			} else {
				return errors.New(fmt.Sprintf("Cannot parse rule in tournament =%d", t.Id))
			}
			if val, ok := mapOfTeams[rule[1]]; ok {
				log.Infof(c, "Update Next phase: match found: %v", val.Name)
				matches[i].TeamId2 = val.Id
			} else {
				return errors.New(fmt.Sprintf("Cannot parse rule in tournament =%d", t.Id))
			}
			matches[i].Rule = ""
		}
		if err := UpdateMatches(c, matches); err != nil {
			log.Errorf(c, "Set Results: unable to set results on matches: %v", err)
			return err
		}
	}

	return nil
}

// From a tournament entity return an array of Matches.
func GetAllMatchesFromTournament(c appengine.Context, tournament *Tournament) []*Tmatch {

	matches := Matches(c, tournament.Matches1stStage)
	matches2ndPhase := Matches(c, tournament.Matches2ndStage)

	// append 2nd round to first one
	for _, m := range matches2ndPhase {
		matches = append(matches, m)
	}

	return matches
}

func getMatchesByPhase(c appengine.Context, t *Tournament, phaseName string) []*Tmatch {

	limits := MapOfPhaseIntervals()

	low := limits[phaseName][0]
	high := limits[phaseName][1]

	matches := GetAllMatchesFromTournament(c, t)

	var filteredMatches []*Tmatch
	for i, v := range matches {
		if v.IdNumber >= low && v.IdNumber <= high {
			filteredMatches = append(filteredMatches, matches[i])
		}
	}
	return filteredMatches
}

func MatchesGroupByPhase(matches []*Tmatch) []Tphase {
	limits := MapOfPhaseIntervals()
	phaseNames := ArrayOfPhases()

	phases := make([]Tphase, len(limits))
	for i, _ := range phases {
		phases[i].Name = phaseNames[i]
		low := limits[phases[i].Name][0]
		high := limits[phases[i].Name][1]

		var filteredMatches []Tmatch
		for _, v := range matches {
			if v.IdNumber >= low && v.IdNumber <= high {
				filteredMatches = append(filteredMatches, *v)
			}
		}
		phases[i].Days = MatchesGroupByDay(filteredMatches)
	}
	return phases
}

func MatchesGroupByDay(matches []Tmatch) []Tday {

	mapOfDays := make(map[string][]Tmatch)

	const shortForm = "Jan/02/2006"
	for _, m := range matches {
		currentDate := m.Date.Format(shortForm)
		_, ok := mapOfDays[currentDate]
		if ok {
			mapOfDays[currentDate] = append(mapOfDays[currentDate], m)
		} else {
			var arrayMatches []Tmatch
			arrayMatches = append(arrayMatches, m)
			mapOfDays[currentDate] = arrayMatches
		}
	}

	var days []Tday
	days = make([]Tday, len(mapOfDays))
	i := 0
	for key, value := range mapOfDays {
		days[i].Date, _ = time.Parse(shortForm, key)
		days[i].Matches = value
		i++
	}

	sort.Sort(ByDate(days))
	return days
}

// ByDate implements sort.Interface for []Tday based on the date field.
type ByDate []Tday

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func ArgMaxInt64(arr []int64) (index int, max int64) {
	for i, n := range arr {
		if i == 0 || n > max {
			index = i
			max = n
		}
	}
	return
}
func getFirstTeamInGroup(c appengine.Context, g *Tgroup) (*Tteam, int) {
	points := make([]int64, len(g.Points))
	copy(points, g.Points)

	argmax1, max1 := ArgMaxInt64(points)
	points[argmax1] = -1
	argmax2, max2 := ArgMaxInt64(points)
	if max1 == max2 { // points equal try by difference
		diff := make([]int64, len(points))
		for i, _ := range points {
			diff[i] = g.GoalsF[i] - g.GoalsA[i]
		}
		if diff[argmax1] > diff[argmax2] {
			return &g.Teams[argmax1], argmax1
		} else if diff[argmax1] < diff[argmax2] {
			return &g.Teams[argmax2], argmax2
		} else { // diff equal try by greatest number of goals scored
			if g.GoalsF[argmax1] > g.GoalsF[argmax2] {
				return &g.Teams[argmax1], argmax1
			} else if g.GoalsF[argmax1] < g.GoalsF[argmax2] {
				return &g.Teams[argmax2], argmax2
			} else { // still equal try at random for now...
				if rand.Intn(2) == 0 {
					return &g.Teams[argmax1], argmax1
				}
				return &g.Teams[argmax2], argmax2
			}
		}
	} else {
		return &g.Teams[argmax1], argmax1
	}
}

func getSecondTeamInGroup(c appengine.Context, g *Tgroup, indexOfFirst int) (*Tteam, int) {

	points := make([]int64, len(g.Points))
	copy(points, g.Points)

	points[indexOfFirst] = -1

	argmax1, max1 := ArgMaxInt64(points)
	points[argmax1] = -1
	argmax2, max2 := ArgMaxInt64(points)
	if max1 == max2 { // points equal try by difference
		diff := make([]int64, len(points))
		for i, _ := range points {
			diff[i] = g.GoalsF[i] - g.GoalsA[i]
		}
		if diff[argmax1] > diff[argmax2] {
			return &g.Teams[argmax1], argmax1
		} else if diff[argmax1] < diff[argmax2] {
			return &g.Teams[argmax2], argmax2
		} else { // diff equal try by greatest number of goals scored
			if g.GoalsF[argmax1] > g.GoalsF[argmax2] {
				return &g.Teams[argmax1], argmax1
			} else if g.GoalsF[argmax1] < g.GoalsF[argmax2] {
				return &g.Teams[argmax2], argmax2
			} else { // still equal try at random for now...
				if rand.Intn(2) == 0 {
					return &g.Teams[argmax1], argmax1
				}
				return &g.Teams[argmax2], argmax2
			}
		}
	} else {
		return &g.Teams[argmax1], argmax1
	}
}
