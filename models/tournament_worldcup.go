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
	"fmt"
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
)

const (
	cFirstStage    = "First Stage"
	cRoundOf16     = "Round of 16"
	cQuarterFinals = "Quarter-finals"
	cSemiFinals    = "Semi-finals"
	cThirdPlace    = "Third Place"
	cFinals        = "Finals"
)

// Map of groups, key: group name, value: string array of teams.
func MapOfGroups() map[string][]string {
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

	return mapWCGroups
}

// Map of country codes, key: team name, value: ISO code
// example: Brazil: BR
func MapOfCountryCodes() map[string]string {

	var codes map[string]string
	codes = make(map[string]string)

	codes["Brazil"] = "br"
	codes["Croatia"] = "hr"
	codes["Mexico"] = "mx"
	codes["Cameroon"] = "cm"
	codes["Spain"] = "es"
	codes["Netherlands"] = "nl"
	codes["Chile"] = "cl"
	codes["Australia"] = "at"
	codes["Colombia"] = "co"
	codes["Greece"] = "gr"
	codes["Côte d'Ivoire"] = "ci"
	codes["Japan"] = "jp"
	codes["Uruguay"] = "uy"
	codes["Costa Rica"] = "cr"
	codes["England"] = "gb"
	codes["Italy"] = "it"
	codes["Switzerland"] = "ch"
	codes["Ecuador"] = "ec"
	codes["France"] = "fr"
	codes["Honduras"] = "hn"
	codes["Argentina"] = "ar"
	codes["Bosnia-Herzegovina"] = "ba"
	codes["Iran"] = "ir"
	codes["Nigeria"] = "ng"
	codes["Germany"] = "de"
	codes["Portugal"] = "pt"
	codes["Ghana"] = "gh"
	codes["United States"] = "us"
	codes["Belgium"] = "be"
	codes["Algeria"] = "dz"
	codes["Russia"] = "ru"
	codes["South Korea"] = "kr"

	return codes
}

// Map of group matches, key: group name, value: array of array of strings with match information ( MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation)
//
// Example:
//
// 	Group A:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
func MapOfGroupMatches() map[string][][]string {

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

	return mapGroupMatches
}

// Returns the Map of 2nd round matches, of the world cup tournament.
// key: round number, value: array of array of strings with match information ( MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation)
//
// Example:
//
//	round 16:[{"1", "Jun/12/2014", "Brazil", "Croatia", "Arena de São Paulo, São Paulo"}, ...]
//
func MapOf2ndRoundMatches() map[string][][]string {

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
	m2nd14 := []string{"62", "Jul/09/2014", "W59", "W60", "São Paulo"}
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

	return mapMatches2ndRound
}

// Return an array of the phases names of world cup tournament: (FirstStage, RoundOf16, QuarterFinals, SemiFinals, ThirdPlace, Finals)
func ArrayOfPhases() []string {
	return []string{cFirstStage, cRoundOf16, cQuarterFinals, cSemiFinals, cThirdPlace, cFinals}
}

// Build a map with key the corresponding phase in the world cup tournament
// at value a tuple that represent the match number interval in which the phase take place:
// first stage: matches 1 to 48
// Round of 16: matches 49 to 56
// Quarte-finals: matches 57 to 60
// Semi-finals: matches 61 to 62
// Thrid Place: match 63
// Finals: match 64
func MapOfPhaseIntervals() map[string][]int64 {

	limits := make(map[string][]int64)
	limits[cFirstStage] = []int64{1, 48}
	limits[cRoundOf16] = []int64{49, 56}
	limits[cQuarterFinals] = []int64{57, 60}
	limits[cSemiFinals] = []int64{61, 62}
	limits[cThirdPlace] = []int64{63, 63}
	limits[cFinals] = []int64{64, 64}
	return limits
}

// Create world cup tournament entity 2014.
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

	var mapCountryCodes map[string]string
	mapCountryCodes = MapOfCountryCodes()
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

	// for date parsing
	const shortForm = "Jan/02/2006"

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

			teamID, _, err1 := datastore.AllocateIDs(c, "Tteam", nil, 1)
			if err1 != nil {
				return nil, err1
			}
			log.Infof(c, "World Cup: team: %v allocateIDs ok", teamName)

			teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
			log.Infof(c, "World Cup: team: %v NewKey ok", teamName)

			team := &Tteam{teamID, teamName, mapCountryCodes[teamName]}
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

		groupMatches := mapGroupMatches[groupName]
		group.Matches = make([]Tmatch, len(groupMatches))

		for matchIndex, matchData := range groupMatches {
			log.Infof(c, "World Cup: match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

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
				false,
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

		groupID, _, err1 := datastore.AllocateIDs(c, "Tgroup", nil, 1)
		if err1 != nil {
			return nil, err1
		}

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
	userIds := make([]int64, 0)
	teamIds := make([]int64, 0)
	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: {"64", "Jul/13/2014", "W61", "W62", "Rio de Janeiro"}
	mapMatches2ndRound := MapOf2ndRoundMatches()

	// build matches 2nd phase
	for roundNumber, roundMatches := range mapMatches2ndRound {
		log.Infof(c, "World Cup: building 2nd round matches: round number %v", roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "World Cup: second phase match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

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
				false,
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

	tstart, _ := time.Parse(shortForm, "Jun/12/2014")
	tend, _ := time.Parse(shortForm, "Jul/13/2014")
	adminIds := make([]int64, 1)
	adminIds[0] = adminId
	tournament := &Tournament{
		tournamentID,
		helpers.TrimLower("world cup"),
		"World Cup",
		"FIFA World Cup",
		tstart,
		tend,
		adminIds,
		time.Now(),
		groupIds,
		matches1stStageIds,
		matches2ndStageIds,
		userIds,
		teamIds,
		false,
	}
	log.Infof(c, "World Cup: instance of tournament ready")

	_, err = datastore.Put(c, key, tournament)
	if err != nil {
		return nil, err
	}
	log.Infof(c, "World Cup:  tournament put in datastore ok")

	AddToTournamentInvertedIndex(c, helpers.TrimLower("world cup"), tournamentID)

	return tournament, nil
}
