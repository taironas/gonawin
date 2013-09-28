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

package search

import (
	"net/http"
	"math"
	"sort"
	"strings"

	"appengine"

	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
	teammdl "github.com/santiaago/purple-wing/models/team"
	helpers "github.com/santiaago/purple-wing/helpers"
)

func Score(r *http.Request, query string, ids []int64){
	c := appengine.NewContext(r)

	words := strings.Split(query, " ")
	setOfWords := helpers.SetOfStrings(query)
	nbTeamWords, _ := teaminvidmdl.GetWordCount(c) 

	// query vector
	q := make([]float64, len(setOfWords))
	for i,w := range setOfWords{
		dft := 0
		if inv_id := teaminvidmdl.Find(r, "KeyName", w); inv_id != nil{
			dft = len(strings.Split(string(inv_id.TeamIds), " "))
		}
		q[i] = math.Log10(1+float64(helpers.CountTerm(words, w))) * math.Log10(float64(nbTeamWords+1)/float64(dft+1))
	}
	c.Infof("query vector: %v",q)

	// d vector
	//nbTeams, _ := teammdl.GetTeamCounter(c)
	vec_d := make([][]float64, len(ids))
	for i, id := range ids{
		d := make([]float64, len(setOfWords))
		for j, wi := range setOfWords{
			// get word frequency by team (id, wi)
			wordFreqByTeam := teammdl.GetWordFrequencyForTeam(r, id, wi)
			// get number of teams with word (wi)
			teamFreqForWord := teaminvidmdl.GetTeamFrequencyForWord(r, wi)
			d[j] = math.Log10(float64(1+wordFreqByTeam)) * math.Log10(float64(nbTeamWords+1)/float64(teamFreqForWord+1))
		}
		vec_d[i] = make([]float64, len(setOfWords))
		vec_d[i] = d		
	}
	c.Infof("d vector: %v",vec_d)

	// compute score vector
	var score map[int64]float64
	score = make(map[int64]float64)
	for i, vec_di := range vec_d{
		score[ids[i]] = dotProduct(vec_di, q)
	}
	c.Infof("score vector :%v", score)

	sortedScore := make([]int, len(score))
	i := 0
	for k, _ := range score{
		sortedScore[i] = int(k)
		i++
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedScore)))
	c.Infof("sorted score vector %v",sortedScore)
}

func dotProduct(vec1 []float64,vec2 []float64)float64{
	if len(vec1) != len(vec2){
		return 0
	}else{
		sum := float64(0)
		for i, _ :=  range vec1 {
			sum = sum + vec1[i]*vec2[i]
		}
		return sum
	}
}















