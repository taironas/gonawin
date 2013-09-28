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
	"strings"

	"appengine"

	teaminvidmdl "github.com/santiaago/purple-wing/models/teamInvertedIndex"
)

func Score(r *http.Request, words []string, ids []int64){
	c := appengine.NewContext(r)
		
	nbTeamWords, _ := teaminvidmdl.GetWordCount(c) 
	i := 0
	q := make([]float64, len(words))
	for _,w := range words{
		dft := 0
		if inv_id := teaminvidmdl.Find(r, "KeyName", w); inv_id != nil{
			dft = len(strings.Split(string(inv_id.TeamIds), " "))
		}
		q[i] = math.Log10(1+float64(countTerm(words, w))) * math.Log10(float64(nbTeamWords+1)/float64(dft+1))
		i = i + 1
	}
	// d
	//nbTeam := 
}

func countTerm(words []string, w string)int64{
	var c int64 = 0
	for _,wi := range words{
		if wi == w{
			c = c + 1
		}
	}
	return c
}
