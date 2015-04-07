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
	"strconv"
	"time"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers/log"
)

const (
	cMatchId       = 0
	cMatchDate     = 1
	cMatchTeam1    = 2
	cMatchTeam2    = 3
	cMatchLocation = 4
)

type ChampionsLeagueTournament struct{}

// Map of groups, key: group name, value: string array of teams.
func (clt ChampionsLeagueTournament) MapOfGroups() map[string][]string {
	mapWCGroups := make(map[string][]string)
	return mapWCGroups
}

// Map of team codes, key: team name, value: code
// example: Paris Saint-Germain: PSG
func (clt ChampionsLeagueTournament) MapOfTeamCodes() map[string]string {

	var codes map[string]string
	codes = make(map[string]string)

	codes["Club Athletico de Madrid"] = "ATM"
	codes["FC Barcelona"] = "FCB"
	codes["FC Bayern Munchen"] = "BYM"
	codes["Juventus"] = "JUV"
	codes["AS Monaco FC"] = "ASM"
	codes["Paris Saint-Germain"] = "PSG"
	codes["FC Porto"] = "POR"
	codes["Real Madrid CF"] = "RLM"

	return codes
}

func (clt ChampionsLeagueTournament) MapOfGroupMatches() map[string][][]string {
	mapGroupMatches := make(map[string][][]string)
	return mapGroupMatches
}

// Returns the Map of 2nd round matches, of the world cup tournament.
// key: round number, value: array of array of strings with match information ( MatchId, MatchDate, MatchTeam1, MatchTeam2, MatchLocation)
//
// Example:
//
//	round 4:[{"1", "Apr/14/2014", "Paris Saint-Germain", "AS Monaco FC", "Parc des Princes, Paris"}, ...]
//
func (clt ChampionsLeagueTournament) MapOf2ndRoundMatches() map[string][][]string {

	// Quarter-finals
	m2nd1 := []string{"1", "Apr/14/2014", "Juventus", "AS Monaco FC", "Juventus Stadium, Turin"}
	m2nd2 := []string{"2", "Apr/14/2014", "FC Porto", "FC Bayern Munchen", "Estádio do Dragão, Porto"}
	m2nd3 := []string{"3", "Apr/15/2014", "Paris Saint-Germain", "FC Barcelona", "Parc des Princes, Paris"}
	m2nd4 := []string{"4", "Apr/15/2014", "Club Athletico de Madrid", "Real Madrid CF", "Stade Vicente-Calderón, Madrid"}
	m2nd5 := []string{"5", "Apr/21/2014", "AS Monaco FC", "Juventus", "Stade Louis-II, Monaco"}
	m2nd6 := []string{"6", "Apr/21/2014", "FC Bayern Munchen", "FC Porto", "Allianz Arena, Munchen"}
	m2nd7 := []string{"7", "Apr/22/2014", "FC Barcelona", "Paris Saint-Germain", "Camp Nou, Barcelona"}
	m2nd8 := []string{"8", "Apr/22/2014", "Real Madrid CF", "Club Athletico de Madrid", "Stade Santiago Bernabéu, Madrid"}
	// Semi-finals
	m2nd9 := []string{"9", "May/05/2014", "TBD1", "TBD2", "TBD"}
	m2nd10 := []string{"10", "May/06/2014", "TBD3", "TBD4", "TBD"}
	m2nd11 := []string{"11", "May/12/2014", "TBD5", "TBD6", "TBD"}
	m2nd12 := []string{"12", "May/13/2014", "TBD7", "TBD8", "TBD"}
	// Final
	m2nd13 := []string{"13", "Jun/06/2015", "TBD9", "TBD10", "Olympiastadion, Berlin"}

	var quarterFinals [][]string
	var semiFinals [][]string
	var final [][]string

	quarterFinals = append(quarterFinals, m2nd1, m2nd2, m2nd3, m2nd4, m2nd5, m2nd6, m2nd7, m2nd8)
	semiFinals = append(semiFinals, m2nd9, m2nd10, m2nd11, m2nd12)
	final = append(final, m2nd13)

	mapMatches2ndStage := make(map[string][][]string)
	mapMatches2ndStage[cQuarterFinals] = quarterFinals
	mapMatches2ndStage[cSemiFinals] = semiFinals
	mapMatches2ndStage[cFinals] = final

	return mapMatches2ndStage
}

// Return an array of the phases names of champions league tournament: (QuarterFinals, SemiFinals, Finals)
func (clt ChampionsLeagueTournament) ArrayOfPhases() []string {
	return []string{cQuarterFinals, cSemiFinals, cFinals}
}

// Build a map with key the corresponding phase in the champions league tournament
// at value a tuple that represent the match number interval in which the phase take place:
// Quarter-finals: matches 1 to 8
// Semi-finals: matches 9 to 12
// Finals: match 13
func (clt ChampionsLeagueTournament) MapOfPhaseIntervals() map[string][]int64 {

	limits := make(map[string][]int64)
	limits[cQuarterFinals] = []int64{1, 8}
	limits[cSemiFinals] = []int64{9, 12}
	limits[cFinals] = []int64{13, 13}
	return limits
}

// Create champions league tournament entity 2015.
func CreateChampionsLeague(c appengine.Context, adminId int64) (*Tournament, error) {
	// create new tournament
	log.Infof(c, "Champions League: start")

	mapTeamId := make(map[string]int64)

	// for date parsing
	const shortForm = "Jan/02/2006"

	// mapMatches2ndRound  is a map where the key is a string which represent the rounds
	// the key is a two dimensional string array. each element in the array represent a specific field in the match
	// mapMatches2ndRound is a map[string][][]string
	// example: "1", "Apr/14/2014", "Paris Saint-Germain", "AS Monaco FC", "Parc des Princes, Paris"}
	clt := ChampionsLeagueTournament{}
	clMatches2ndStage := clt.MapOf2ndRoundMatches()
	clMapTeamCodes := clt.MapOfTeamCodes()

	// matches 2nd stage
	matches2ndStageIds := make([]int64, 0)
	userIds := make([]int64, 0)
	teamIds := make([]int64, 0)

	log.Infof(c, "Champions League: maps ready")

	// build teams
	for teamName, teamCode := range clMapTeamCodes {

		teamID, _, err1 := datastore.AllocateIDs(c, "Tteam", nil, 1)
		if err1 != nil {
			return nil, err1
		}
		log.Infof(c, "Champions League: team: %v allocateIDs ok", teamName)

		teamkey := datastore.NewKey(c, "Tteam", "", teamID, nil)
		log.Infof(c, "Champions League: team: %v NewKey ok", teamName)

		team := &Tteam{teamID, teamName, teamCode}
		log.Infof(c, "Champions League: team: %v instance of team ok", teamName)

		_, err := datastore.Put(c, teamkey, team)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "Champions League: team: %v put in datastore ok", teamName)
		mapTeamId[teamName] = teamID
	}

	// build tmatches quarter finals
	log.Infof(c, "Champions League: building quarter finals matches")
	for _, matchData := range clMatches2ndStage[cQuarterFinals] {
		log.Infof(c, "Champions League: quarter finals match data: %v", matchData)

		matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
		if err1 != nil {
			return nil, err1
		}

		log.Infof(c, "Champions League: match: %v allocateIDs ok", matchID)

		matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
		log.Infof(c, "Champions League: match: new key ok")

		matchTime, _ := time.Parse(shortForm, matchData[cMatchDate])
		matchInternalId, _ := strconv.Atoi(matchData[cMatchId])

		rule := fmt.Sprintf("%s %s", matchData[cMatchTeam1], matchData[cMatchTeam2])
		emptyresult := int64(0)
		match := &Tmatch{
			matchID,
			int64(matchInternalId),
			matchTime,
			mapTeamId[matchData[cMatchTeam1]],
			mapTeamId[matchData[cMatchTeam2]],
			matchData[cMatchLocation],
			rule,
			emptyresult,
			emptyresult,
			false,
			true,
			true,
		}
		log.Infof(c, "Champions League: match 2nd round: build match ok")

		_, err := datastore.Put(c, matchkey, match)
		if err != nil {
			return nil, err
		}
		log.Infof(c, "Champions League: 2nd round match: %v put in datastore ok", matchData)
		// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
		matches2ndStageIds = append(matches2ndStageIds, matchID)
	}

	// build tmatches 2nd phase from semi finals
	restofMatches2ndStage := clMatches2ndStage
	delete(restofMatches2ndStage, cQuarterFinals)
	for roundNumber, roundMatches := range restofMatches2ndStage {
		log.Infof(c, "Champions League: building 2nd stage matches: round number %v", roundNumber)
		for _, matchData := range roundMatches {
			log.Infof(c, "Champions League: second stage match data: %v", matchData)

			matchID, _, err1 := datastore.AllocateIDs(c, "Tmatch", nil, 1)
			if err1 != nil {
				return nil, err1
			}

			log.Infof(c, "Champions League: match: %v allocateIDs ok", matchID)

			matchkey := datastore.NewKey(c, "Tmatch", "", matchID, nil)
			log.Infof(c, "Champions League: match: new key ok")

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
				false,
				true,
			}
			log.Infof(c, "Champions League: match 2nd round: build match ok")

			_, err := datastore.Put(c, matchkey, match)
			if err != nil {
				return nil, err
			}
			log.Infof(c, "Champions League: 2nd round match: %v put in datastore ok", matchData)
			// save in an array of int64 all the allocate IDs to store them in the tournament for easy retreival later on.
			matches2ndStageIds = append(matches2ndStageIds, matchID)
		}
	}

	tstart, _ := time.Parse(shortForm, "Apr/14/2015")
	tend, _ := time.Parse(shortForm, "Jun/06/2015")
	adminIds := make([]int64, 1)
	adminIds[0] = adminId
	name := "2014-2015 UEFA Champions League"
	description := ""
	var tournament *Tournament
	var err error
	if tournament, err = CreateTournament(c, name, description, tstart, tend, adminId); err != nil {
		log.Infof(c, "Champions League: something went wrong when creating tournament.")
		return nil, err
	} else {
		tournament.GroupIds = make([]int64, 0)
		tournament.Matches1stStage = make([]int64, 0)
		tournament.Matches2ndStage = matches2ndStageIds
		tournament.UserIds = userIds
		tournament.TeamIds = teamIds
		tournament.TwoLegged = false
		tournament.IsFirstStageComplete = false
		if err1 := tournament.Update(c); err1 != nil {
			log.Infof(c, "Champions League: unable to udpate tournament.")
		}
	}

	log.Infof(c, "Champions League: instance of tournament ready")

	return tournament, nil
}

// From tournament entity build map of teams.
func (clt ChampionsLeagueTournament) MapOfIdTeams(c appengine.Context, tournament *Tournament) map[int64]string {

	mapIdTeams := make(map[int64]string)

	matches2ndStage := tournament.Matches2ndStage

	for i := 0; i < 4; i++ {
		matchID := matches2ndStage[i]

		m, err := MatchById(c, matchID)
		if err != nil {
			log.Errorf(c, " MapOfIdTeams, cannot find match with ID=%", matchID)
		} else {
			t1, err1 := TTeamById(c, m.TeamId1)
			if err1 != nil {
				log.Errorf(c, " MapOfIdTeams, cannot find tteam with ID=%", m.TeamId1)
			} else {
				mapIdTeams[t1.Id] = t1.Name
			}

			t2, err2 := TTeamById(c, m.TeamId2)
			if err2 != nil {
				log.Errorf(c, " MapOfIdTeams, cannot find tteam with ID=%", m.TeamId2)
			} else {
				mapIdTeams[t2.Id] = t2.Name
			}
		}
	}

	return mapIdTeams
}
