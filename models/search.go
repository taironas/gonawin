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
	"math"
	"sort"
	"strings"

	"appengine"

	helpers "github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
)

// Given a query string and an array of ids, computes a score vector that has the doc ids and the score of each id with respect to the query.
func TournamentScore(c appengine.Context, query string, ids []int64) []int64 {

	words := strings.Split(query, " ")
	setOfWords := helpers.SetOfStrings(query)
	nbTournamentWords, _ := TournamentInvertedIndexGetWordCount(c)

	// query vector
	q := make([]float64, len(setOfWords))
	for i, w := range setOfWords {
		dft := 0
		if invId, err := FindTournamentInvertedIndex(c, "KeyName", w); err != nil {
			log.Errorf(c, " search.TournamentScore, unable to find KeyName=%s: %v", w, err)
		} else if invId != nil {
			dft = len(strings.Split(string(invId.TournamentIds), " "))
		}
		q[i] = math.Log10(1+float64(helpers.CountTerm(words, w))) * math.Log10(float64(nbTournamentWords+1)/float64(dft+1))
	}
	log.Infof(c, "query vector: %v", q)

	// d vector
	vec_d := make([][]float64, len(ids))
	for i, id := range ids {
		d := make([]float64, len(setOfWords))
		for j, wi := range setOfWords {
			// get word frequency by tournament (id, wi)
			wordFreqByTournament := GetWordFrequencyForTournament(c, id, wi)
			// get number of tournaments with word (wi)
			tournamentFreqForWord, err := GetTournamentFrequencyForWord(c, wi)
			if err != nil {
				log.Errorf(c, " search.TournamentScore, error occurred when getting tournament frequency for word=%s: %v", wi, err)
			}

			d[j] = math.Log10(float64(1+wordFreqByTournament)) * math.Log10(float64(nbTournamentWords+1)/float64(tournamentFreqForWord+1))
		}
		vec_d[i] = make([]float64, len(setOfWords))
		vec_d[i] = d
	}
	log.Infof(c, "d vector: %v", vec_d)

	// compute score vector
	var score map[int64]float64
	score = make(map[int64]float64)
	for i, vec_di := range vec_d {
		score[ids[i]] = dotProduct(vec_di, q)
	}
	log.Infof(c, "score vector :%v", score)
	sortedScore := sortMapByValueDesc(score)
	log.Infof(c, "sorted score: %v", sortedScore)

	return getKeysFrompairList(sortedScore)
}

// Given a query string and an array of ids, computes a score vector that has the doc ids and the score of each id with respect to the query.
func TeamScore(c appengine.Context, query string, ids []int64) []int64 {

	words := strings.Split(query, " ")
	setOfWords := helpers.SetOfStrings(query)
	nbTeamWords, _ := TeamInvertedIndexGetWordCount(c)

	// query vector
	q := make([]float64, len(setOfWords))
	for i, w := range setOfWords {
		dft := 0
		if invId, err := FindTeamInvertedIndex(c, "KeyName", w); err != nil {
			log.Errorf(c, " search.TeamScore, unable to find KeyName=%s: %v", w, err)
		} else if invId != nil {
			dft = len(strings.Split(string(invId.TeamIds), " "))
		}
		q[i] = math.Log10(1+float64(helpers.CountTerm(words, w))) * math.Log10(float64(nbTeamWords+1)/float64(dft+1))
	}
	log.Infof(c, "query vector: %v", q)

	// d vector
	vec_d := make([][]float64, len(ids))
	for i, id := range ids {
		d := make([]float64, len(setOfWords))
		for j, wi := range setOfWords {
			// get word frequency by team (id, wi)
			wordFreqByTeam := GetWordFrequencyForTeam(c, id, wi)
			// get number of teams with word (wi)
			teamFreqForWord, err := GetTeamFrequencyForWord(c, wi)
			if err != nil {
				log.Errorf(c, " search.TeamScore, error occurred when getting team frequency for word=%s: %v", wi, err)
			}

			d[j] = math.Log10(float64(1+wordFreqByTeam)) * math.Log10(float64(nbTeamWords+1)/float64(teamFreqForWord+1))
		}
		vec_d[i] = make([]float64, len(setOfWords))
		vec_d[i] = d
	}
	log.Infof(c, "d vector: %v", vec_d)

	// compute score vector
	var score map[int64]float64
	score = make(map[int64]float64)
	for i, vec_di := range vec_d {
		score[ids[i]] = dotProduct(vec_di, q)
	}
	log.Infof(c, "score vector :%v", score)
	sortedScore := sortMapByValueDesc(score)
	log.Infof(c, "sorted score: %v", sortedScore)

	return getKeysFrompairList(sortedScore)
}

// Given a query string and an array of ids, computes a score vector that has the doc ids and the score of each id with respect to the query.
func UserScore(c appengine.Context, query string, ids []int64) []int64 {

	words := strings.Split(query, " ")
	setOfWords := helpers.SetOfStrings(query)
	nbUserWords, _ := UserInvertedIndexGetWordCount(c)

	// query vector
	q := make([]float64, len(setOfWords))
	for i, w := range setOfWords {
		dft := 0
		if invId, err := FindUserInvertedIndex(c, "KeyName", w); err != nil {
			log.Errorf(c, "search.UserScore, unable to find KeyName=%s: %v", w, err)
		} else if invId != nil {
			dft = len(strings.Split(string(invId.UserIds), " "))
		}
		q[i] = math.Log10(1+float64(helpers.CountTerm(words, w))) * math.Log10(float64(nbUserWords+1)/float64(dft+1))
	}
	log.Infof(c, "query vector: %v", q)

	// d vector
	vec_d := make([][]float64, len(ids))
	for i, id := range ids {
		d := make([]float64, len(setOfWords))
		for j, wi := range setOfWords {
			// get word frequency by user (id, wi)
			wordFreqByUser := GetWordFrequencyForUser(c, id, wi)
			// get number of users with word (wi)
			userFreqForWord, err := GetUserFrequencyForWord(c, wi)
			if err != nil {
				log.Errorf(c, " search.UserScore, error occurred when getting user frequency for word=%s: %v", wi, err)
			}

			d[j] = math.Log10(float64(1+wordFreqByUser)) * math.Log10(float64(nbUserWords+1)/float64(userFreqForWord+1))
		}
		vec_d[i] = make([]float64, len(setOfWords))
		vec_d[i] = d
	}
	log.Infof(c, "d vector: %v", vec_d)

	// compute score vector
	var score map[int64]float64
	score = make(map[int64]float64)
	for i, vec_di := range vec_d {
		score[ids[i]] = dotProduct(vec_di, q)
	}
	log.Infof(c, "score vector :%v", score)
	sortedScore := sortMapByValueDesc(score)
	log.Infof(c, "sorted score: %v", sortedScore)

	return getKeysFrompairList(sortedScore)
}

// A data structure to hold a key/value pair.
type pair struct {
	Key   int64
	Value float64
}

// A slice of pairs that implements sort.Interface to sort by Value.
type pairList []pair

func (p pairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p pairList) Len() int           { return len(p) }
func (p pairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a pairList, then sort and return it.
func sortMapByValue(m map[int64]float64) pairList {
	p := make(pairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

// A function to turn a map into a pairList, then sort in descending order and return it.
func sortMapByValueDesc(m map[int64]float64) pairList {
	p := make(pairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}

// From a pair list returns an array of keys present in the pair list.
func getKeysFrompairList(p pairList) []int64 {
	keys := make([]int64, len(p))
	i := 0
	for _, pair := range p {
		keys[i] = pair.Key
		i++
	}
	return keys
}

// Compute the dot product of two float vectors.
func dotProduct(vec1 []float64, vec2 []float64) float64 {
	if len(vec1) != len(vec2) {
		return 0
	} else {
		sum := float64(0)
		for i, _ := range vec1 {
			sum = sum + vec1[i]*vec2[i]
		}
		return sum
	}
}
