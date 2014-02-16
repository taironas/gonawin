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
	Matches1stStage []int64    `json:",omitempty"`
	Matches2ndStage *[]int64   `json:",omitempty"`
}

type TournamentCounter struct {
	Count int64
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

// update a tournament given a tournament id and a tournament pointer.
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

// increment tournament counter
func incrementTournamentCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TournamentCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// decrement tournament counter
func decrementTournamentCounter(c appengine.Context, key *datastore.Key) (int64, error) {
	var x TournamentCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// get the current tournament counter
func GetTournamentCounter(c appengine.Context) (int64, error) {
	key := datastore.NewKey(c, "TournamentCounter", "singleton", 0, nil)
	var x TournamentCounter
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
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
	mapWCGroups = make(map[string][]string)
	mapWCGroups["A"] = []string{"Brazil", "Croatia", "Mexico", "Cameroon"}
	mapWCGroups["B"] = []string{"Spain", "Netherlands", "Chile", "Australia"}
	mapWCGroups["C"] = []string{"Colombia", "Greece", "Côte d'Ivoire", "Japan"}
	mapWCGroups["D"] = []string{"Uruguay", "Costa Rica", "England", "Italy"}
	mapWCGroups["E"] = []string{"Switzerland", "Ecuador", "France", "Honduras"}
	mapWCGroups["F"] = []string{"Argentina", "Bosnia-Herzegovina", "Iran", "Nigeria"}
	mapWCGroups["G"] = []string{"Germany", "Portugal", "Ghana", "United States"}
	mapWCGroups["H"] = []string{"Belgium", "Algeria", "Russia", "South Korea"}

	// build map of matches
	var mapTeamId map[string]int64
	mapTeamId = make(map[string]int64)

	mapGroupMatches := make(map[string][][]string)

	const (
		cMatchId       = 0
		cMatchDate     = 1
		cMatchTeam1    = 2
		cMatchTeam2    = 3
		cMatchLocation = 4
	)

	mA1 := []string{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}
	mA2 := []string{"2", "Jun/13/2014", "Mexico", "Cameroon", "Estádio das Dunas, Natal"}
	mA3 := []string{"17", "Jun/17/2014", "Brazil", "Mexico", "Estádio Castelão, Fortaleza"}
	mA4 := []string{"18", "Jun/18/2014", "Cameroon", "Croatia", "Arena Amazônia, Manaus"}
	mA5 := []string{"33", "Jun/23/2014", "Cameroon", "Brazil", "Brasília"}
	mA6 := []string{"34", "Jun/23/2014", "Croatia", "Mexico", "Recife"}

	mB1 := []string{"3", "Jun/13/2014", "Spain", "Netherlands", "Arena Fonte Nova, Salvador"}
	mB2 := []string{"4", "Jun/13/2014", "Chile", "Australia", "Arena Pantanal, Cuiabá"}
	mB3 := []string{"19", "Jun/18/2014", "Spain", "Chile", "Estádio do Maracanã, Rio de Janeiro"}
	mB4 := []string{"20", "Jun/18/2014", "Australia", "Netherlands", "Estádio BeiraRio, Porto Alegre"}
	mB5 := []string{"35", "Jun/23/2014", "Australia", "Spain", "Curitiba"}
	mB6 := []string{"36", "Jun/23/2014", "Netherlands", "Chile", "São Paulo"}

	mC1 := []string{"5", "Jun/14/2014", "Colombia", "Greece", "Estádio Mineirão, Belo Horizonte"}
	mC2 := []string{"6", "Jun/14/2014", "Côte d'Ivoire", "Japan", "Arena Pernambuco, Recife"}
	mC3 := []string{"21", "Jun/19/2014", "Colombia", "Côte d'Ivoire", "Estádio Nacional Mané Garrincha, Brasília"}
	mC4 := []string{"22", "Jun/19/2014", "Japan", "Greece", "Estádio das Dunas, Natal"}
	mC5 := []string{"37", "Jun/24/2014", "Japan", "Colombia", "Cuiabá"}
	mC6 := []string{"38", "Jun/24/2014", "Côte d'Ivoire", "Greece", "Fortaleza"}

	mD1 := []string{"7", "Jun/14/2014", "Uruguay", "Costa Rica", "Estádio Castelão, Fortaleza"}
	mD2 := []string{"8", "Jun/14/2014", "England", "Italy", "Arena Amazônia, Manaus"}
	mD3 := []string{"23", "Jun/19/2014", "Uruguay", "England", "Arena de São Paulo, São Paulo"}
	mD4 := []string{"24", "Jun/20/2014", "Italy", "Costa Rica", "Arena Pernambuco, Recife"}
	mD5 := []string{"39", "Jun/24/2014", "Italy", "Uruguay", "Natal"}
	mD6 := []string{"40", "Jun/24/2014", "Costa Rica", "England", "Belo Horizonte"}

	mE1 := []string{"9", "Jun/15/2014", "Switzerland", "Ecuador", "Estádio Nacional Mané Garrincha, Brasília"}
	mE2 := []string{"10", "Jun/15/2014", "France", "Honduras", "Estádio BeiraRio, Porto Alegre"}
	mE3 := []string{"25", "Jun/20/2014", "Switzerland", "France", "Arena Fonte Nova, Salvador"}
	mE4 := []string{"26", "Jun/20/2014", "Honduras", "Ecuador", "Arena da Baixada, Curitiba"}
	mE5 := []string{"41", "Jun/25/2014", "Honduras", "Switzerland", "Manaus"}
	mE6 := []string{"42", "Jun/25/2014", "Ecuador", "France", "Rio de Janeiro"}

	mF1 := []string{"11", "Jun/15/2014", "Argentina", "Bosnia-Herzegovina", "Estádio do Maracanã, Rio de Janeiro"}
	mF2 := []string{"12", "Jun/16/2014", "Iran", "Nigeria", "Arena da Baixada, Curitiba"}
	mF3 := []string{"27", "Jun/21/2014", "Argentina", "Iran", "Estádio Mineirão, Belo Horizonte"}
	mF4 := []string{"28", "Jun/21/2014", "Nigeria", "Bosnia-Herzegovina", "Arena Pantanal, Cuiabá"}
	mF5 := []string{"43", "Jun/25/2014", "Nigeria", "Argentina", "Porto Alegre"}
	mF6 := []string{"44", "Jun/25/2014", "Bosnia-Herzegovina", "Iran", "Salvador"}

	mG1 := []string{"13", "Jun/16/2014", "Germany", "Portugal", "Arena Fonte Nova, Salvador"}
	mG2 := []string{"14", "Jun/16/2014", "Ghana", "United States", "Estádio das Dunas, Natal"}
	mG3 := []string{"29", "Jun/21/2014", "Germany", "Ghana", "Fortaleza"}
	mG4 := []string{"30", "Jun/22/2014", "United States", "Portugal", "Manaus"}
	mG5 := []string{"45", "Jun/26/2014", "United States", "Germany", "Recife"}
	mG6 := []string{"46", "Jun/26/2014", "Portugal", "Ghana", "Brasília"}

	mH1 := []string{"15", "Jun/17/2014", "Belgium", "Algeria", "Estádio Mineirão, Belo Horizonte"}
	mH2 := []string{"16", "Jun/17/2014", "Russia", "South Korea", "Arena Pantanal, Cuiabá"}
	mH3 := []string{"31", "Jun/22/2014", "Belgium", "Russia", "Rio de Janeiro"}
	mH4 := []string{"32", "Jun/22/2014", "South Korea", "Algeria", "Porto Alegre"}
	mH5 := []string{"47", "Jun/26/2014", "South Korea", "Belgium", "São Paulo"}
	mH6 := []string{"48", "Jun/26/2014", "Algeria", "Russia", "Curitiba"}

	var matchesA [][]string
	var matchesB [][]string
	var matchesC [][]string
	var matchesD [][]string
	var matchesE [][]string
	var matchesF [][]string
	var matchesG [][]string
	var matchesH [][]string

	matchesA = append(matchesA, mA1, mA2, mA3, mA4, mA5, mA6)
	matchesB = append(matchesB, mB1, mB2, mB3, mB4, mB5, mB6)
	matchesC = append(matchesC, mC1, mC2, mC3, mC4, mC5, mC6)
	matchesD = append(matchesD, mD1, mD2, mD3, mD4, mD5, mD6)
	matchesE = append(matchesE, mE1, mE2, mE3, mE4, mE5, mE6)
	matchesF = append(matchesF, mF1, mF2, mF3, mF4, mF5, mF6)
	matchesG = append(matchesG, mG1, mG2, mG3, mG4, mG5, mG6)
	matchesH = append(matchesH, mH1, mH2, mH3, mH4, mH5, mH6)

	mapGroupMatches["A"] = matchesA
	mapGroupMatches["B"] = matchesB
	mapGroupMatches["C"] = matchesC
	mapGroupMatches["D"] = matchesD
	mapGroupMatches["E"] = matchesE
	mapGroupMatches["F"] = matchesF
	mapGroupMatches["G"] = matchesG
	mapGroupMatches["H"] = matchesH

	// matches1stStageIds is an array of  int64
	// where we allocate IDs of the Tmatches entities
	// we will store them in the tournament entity for easy retreival later on.
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
			log.Infof(c, "World Cup: match: %v allocateIDs ok")

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "World Cup: match: new key ok")

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])
			emtpyrule := ""
			match := &Tmatch{
				matchID,
				int64(matchInternalId),
				matchTime,
				mapTeamId[matchData[cMatchTeam1]],
				mapTeamId[matchData[cMatchTeam2]],
				matchData[cMatchLocation],
				emtpyrule,
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

	// matches 2nd stage ----------------------------------------------------
	matches2ndStageIds := make([]int64, 0)
	// Round of 16
	m2nd1 := []string{"49", "Jun/28/2014", "1A", "2B", "Belo Horizonte"}
	m2nd2 := []string{"50", "Jun/28/2014", "1C", "2D", "Rio de Janeiro"}
	m2nd3 := []string{"51", "Jun/29/2014", "1B", "2A", "Fortaleza"}
	m2nd4 := []string{"52", "Jun/29/2014", "1D", "2C", "Recife"}
	m2nd5 := []string{"53", "Jun/30/2014", "1E", "2F", "Brasília"}
	m2nd6 := []string{"54", "Jun/30/2014", "1G", "2H", "Porto Alegre"}
	m2nd7 := []string{"55", "Jul/01/2014", "1F", "2E", "São Paulo"}
	m2nd8 := []string{"56", "Jul/01/2014", "1H", "2G", "Salvador"}
	// 17 Quarter-finals
	m2nd9 := []string{"57", "Jul/04/2014", "W49", "W50", "Fortaleza"}
	m2nd10 := []string{"58", "Jul/04/2014", "W53", "W54", "Rio de Janeiro"}
	m2nd11 := []string{"59", "Jul/05/2014", "W51", "W52", "Salvador"}
	m2nd12 := []string{"60", "Jul/05/2014", "W55", "W56", "Brasília"}
	// 18 Semi-finals
	m2nd13 := []string{"61", "Jul/08/2014", "W57", "W58", "Belo Horizonte"}
	m2nd14 := []string{"62", "Jul/09/2014", "W59", "W69", "São Paulo"}
	//19 Round 19  -  Match for third place
	m2nd15 := []string{"63", "Jul/12/2014", "L61", "L62", "Brasília"}
	//"20" Final
	m2nd16 := []string{"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}

	var round16 [][]string
	var round17 [][]string
	var round18 [][]string
	var round19 [][]string
	var round20 [][]string

	round16 = append(round16, m2nd1, m2nd2, m2nd3, m2nd4, m2nd5, m2nd6, m2nd7, m2nd8)
	round17 = append(round17, m2nd9, m2nd10, m2nd11, m2nd12)
	round18 = append(round18, m2nd13, m2nd14)
	round19 = append(round19, m2nd15)
	round20 = append(round20, m2nd16)

	mapMatches2ndRound := make(map[string][][]string)
	mapMatches2ndRound["16"] = round16
	mapMatches2ndRound["17"] = round17
	mapMatches2ndRound["18"] = round18
	mapMatches2ndRound["19"] = round19
	mapMatches2ndRound["20"] = round20

	// build matches 2nd phase
	const shortForm = "Jan/02/2006"
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "World Cup: building 2nd round matches: round number %v", roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "World Cup: second phase match data: %v", matchData)

			matchID, _, err := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			log.Infof(c, "World Cup: match: %v allocateIDs ok")

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "World Cup: match: new key ok")

			matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
			matchInternalId, _ := strconv.Atoi(matchData[cMatchId])

			rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])

			match := &Tmatch{
				matchID,
				int64(matchInternalId),
				matchTime,
				0, // second round matches start with ids at 0
				0, // second round matches start with ids at 0
				matchData[cMatchLocation],
				rule,
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

	// ----------------------------------------------------------------------
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

// get a Tgroup entity by id
func GroupById(c appengine.Context, groupId int64) (*Tgroup, error) {
	var g Tgroup
	key := datastore.NewKey(c, "Tgroup", "", groupId, nil)

	if err := datastore.Get(c, key, &g); err != nil {
		log.Errorf(c, "group not found : %v", err)
		return &g, err
	}
	return &g, nil
}

// get a Tmatch entity by id
func MatchById(c appengine.Context, matchId int64) (*Tmatch, error) {
	var m Tmatch
	key := datastore.NewKey(c, "Tmatch", "", matchId, nil)

	if err := datastore.Get(c, key, &m); err != nil {
		log.Errorf(c, "match not found : %v", err)
		return &m, err
	}
	return &m, nil
}

// from a tournament id returns an array of groups the participate in it.
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

// from an array of ids return the corresponding array of matches
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

// from tournament entity build map of teams
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
