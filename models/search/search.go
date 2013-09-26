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
	"strings"
	"strconv"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
)

type TeamInvertedIndex struct {
	Id int64
	KeyName string
	TeamIds []byte
}

type TournamentInvertedIndex struct {
	Id int64
	KeyName string
	TournamentIds []byte
}

func CreateTeamInvertedIndex(r *http.Request, word string, teamIds string) *TeamInvertedIndex {
	c := appengine.NewContext(r)
	c.Infof("pw: CreateTeamInvertedIndex")
	id, _, _ := datastore.AllocateIDs(c, "TeamInvertedIndex", nil, 1)
	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)
	
	byteIds := []byte(teamIds)
	t := &TeamInvertedIndex{ id, word, byteIds }

	_, err := datastore.Put(c, key, t)
	if err != nil {
		c.Errorf("pw: CreateTeamInvertedIndex: %v", err)
	}

	return t
}


func CreateTournamentInvertedIndex(r *http.Request, name string, tournamentIds string) *TournamentInvertedIndex {
	c := appengine.NewContext(r)
	
	id, _, _ := datastore.AllocateIDs(c, "TournamentInvertedIndex", nil, 1)
	key := datastore.NewKey(c, "TournamentInvertedIndex", "", id, nil)
	
	byteIds := []byte(tournamentIds)
	t := &TournamentInvertedIndex{ id, helpers.TrimLower(name), byteIds }

	_, err := datastore.Put(c, key, t)
	if err != nil {
		c.Errorf("Create: %v", err)
	}

	return t
}

// AddToTeamInvertedIndex
// Split name by words.
// For each word check if it exist in the Team Inverted Index. 
// if not create a line with the word as key and team id as value.
func AddToTeamInvertedIndex(r *http.Request, name string, id int64){

	c := appengine.NewContext(r)
	c.Infof("pw: AddToTeamInvertedIndex start")
	words := strings.Split(name, " ")
	for _, w:= range words{
		c.Infof("pw: AddToTeamInvertedIndex: Word: %v", w)

		if inv_id := FindTeamInvertedIndex(r, "KeyName", w); inv_id == nil{
			c.Infof("pw: create inv id as word does not exist in table")
			CreateTeamInvertedIndex(r, w, strconv.FormatInt(id, 10))
		} else{
			// update row with new info
			c.Infof("pw: update row with new info")
			c.Infof("pw: current info: keyname: %v", inv_id.KeyName)
			c.Infof("pw: current info: teamIDs: %v", string(inv_id.TeamIds))
			k := KeyByIdTeamInvertedIndex(r, inv_id.Id)

			if newIds := mergeIds(inv_id.TeamIds, id);len(newIds) > 0{
				c.Infof("pw: current info: new team ids: %v", newIds)
				inv_id.TeamIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, inv_id);err != nil{
					c.Errorf("pw: AddToTeamInvertedIndex error on update")
				}
			}
		}
	}
	c.Infof("pw: AddToTeamInvertedIndex end")
	
}

// AddToTournamentInvertedIndex
// Split name by words.
// For each word check if it exist in the Tournament Inverted Index. 
// if not create a line with the word as key and team id as value.
func AddToTournamentInvertedIndex(r *http.Request, name string, id int64){

	c := appengine.NewContext(r)
	c.Infof("pw: AddToTournamentInvertedIndex start")
	words := strings.Split(name, " ")
	for _, w:= range words{
		c.Infof("pw: AddToTournamentInvertedIndex: Word: %v", w)

		if inv_id := FindTournamentInvertedIndex(r, "KeyName", w); inv_id == nil{
			c.Infof("pw: create inv id as word does not exist in table")
			CreateTournamentInvertedIndex(r, w, strconv.FormatInt(id, 10))
		} else{
			// update row with new info
			c.Infof("pw: update row with new info")
			c.Infof("pw: current info: keyname: %v", inv_id.KeyName)
			c.Infof("pw: current info: teamIDs: %v", string(inv_id.TournamentIds))
			k := KeyByIdTournamentInvertedIndex(r, inv_id.Id)

			if newIds := mergeIds(inv_id.TournamentIds, id);len(newIds) > 0{
				c.Infof("pw: current info: new team ids: %v", newIds)
				inv_id.TournamentIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, inv_id);err != nil{
					c.Errorf("pw: AddToTournamentInvertedIndex error on update")
				}
			}
		}
	}
	c.Infof("pw: AddToTournamentInvertedIndex end")	
}
// from the oldname and the new name we handle removal of words no longer present.
// the addition of new words.
func UpdateToTeamInvertedIndex(r *http.Request, oldname string, newname string, id int64){
	c := appengine.NewContext(r)
	c.Infof("pw: UpdateToTeamInvertedIndex start")

	// if word in old and new do nothing
	old_w := strings.Split(oldname," ")
	new_w := strings.Split(newname, " ")

	// remove id  from words in old name that are not present in new name
	for _, wo := range old_w{
		innew := false
		for _, wn := range new_w{
			if wo == wn{
				innew = true
			}
		}
		if !innew{
			c.Infof("pw: remove: %v",wo)
			RemoveWordFromTeamInvertedIndex(r, wo, id)
			//remove it
		}
	}

	// add all id words in new name
	for _, wn := range new_w{
		inold := false
		for _, wo := range old_w{
			if wo == wn{
				inold = true
			}
		}
		if !inold{
			// add it
			c.Infof("pw: add: %v", wn)
			AddWordToTeamInvertedIndex(r, wn, id)
		}
	}

	c.Infof("pw: AddToTeamInvertedIndex end")	

}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well
func RemoveWordFromTeamInvertedIndex(r *http.Request, w string, id int64){

	c := appengine.NewContext(r)
	c.Infof("pw: RemoveWordFromTeamInvertedIndex start")

	if inv_id := FindTeamInvertedIndex(r, "KeyName", w); inv_id == nil{
		c.Infof("pw: word %v does not exist in Team InvertedIndex so nothing to remove",w)
	} else{
		// update row with new info
		k := KeyByIdTeamInvertedIndex(r, inv_id.Id)

		if newIds, err := removeFromIds(inv_id.TeamIds, id); err == nil{
			c.Infof("pw: new team ids after removal: %v", newIds)
			if len(newIds) == 0{
				// this entity does not have ids so remove it from the datastore.
				c.Infof("pw: removing key %v from datastore as it is no longer used", k)
				datastore.Delete(c, k)
			}else{
				inv_id.TeamIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, inv_id);err != nil{
					c.Errorf("pw: RemoveWordFromTeamInvertedIndex error on update")
				}
			}
		}else{
			c.Errorf("pw: unable to remove id from ids: %v",err)
		}
	}
	c.Infof("pw: RemoveWordFromTeamInvertedIndex end")
	
}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well
func RemoveWordFromTournamentInvertedIndex(r *http.Request, w string, id int64){

	c := appengine.NewContext(r)
	c.Infof("pw: RemoveWordFromTournamentInvertedIndex start")

	if inv_id := FindTournamentInvertedIndex(r, "KeyName", w); inv_id == nil{
		c.Infof("pw: word %v does not exist in Tournament InvertedIndex so nothing to remove",w)
	} else{
		// update row with new info
		k := KeyByIdTournamentInvertedIndex(r, inv_id.Id)

		if newIds, err := removeFromIds(inv_id.TournamentIds, id); err == nil{
			c.Infof("pw: new tournament ids after removal: %v", newIds)
			if len(newIds) == 0{
				// this entity does not have ids so remove it from the datastore.
				c.Infof("pw: removing key %v from datastore as it is no longer used", k)
				datastore.Delete(c, k)
			}else{
				inv_id.TournamentIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, inv_id);err != nil{
					c.Errorf("pw: RemoveWordFromTournamentInvertedIndex error on update")
				}
			}
		}else{
			c.Errorf("pw: unable to remove id from ids: %v",err)
		}
	}
	c.Infof("pw: RemoveWordFromTournamentInvertedIndex end")
	
}


// add word to team inverted index is handled the same way as AddToTeamInvertedIndex.
// we add this function for clarity
func AddWordToTeamInvertedIndex(r *http.Request, word string, id int64){
	AddToTeamInvertedIndex(r, word, id)
}


// add word to tournament inverted index is handled the same way as AddToTournamentInvertedIndex.
// we add this function for clarity
func AddWordToTournamentInvertedIndex(r *http.Request, word string, id int64){
	AddToTournamentInvertedIndex(r, word, id)
}

func UpdateTournamentInvertedIndex(r *http.Request, oldname string, newname string, id int64){
	c := appengine.NewContext(r)
	c.Infof("pw: UpdateToTournamentInvertedIndex start")

	// if word in old and new do nothing
	old_w := strings.Split(oldname," ")
	new_w := strings.Split(newname, " ")

	// remove id  from words in old name that are not present in new name
	for _, wo := range old_w{
		innew := false
		for _, wn := range new_w{
			if wo == wn{
				innew = true
			}
		}
		if !innew{
			c.Infof("pw: remove: %v",wo)
			RemoveWordFromTournamentInvertedIndex(r, wo, id)
			//remove it
		}
	}

	// add all id words in new name
	for _, wn := range new_w{
		inold := false
		for _, wo := range old_w{
			if wo == wn{
				inold = true
			}
		}
		if !inold{
			// add it
			c.Infof("pw: add: %v", wn)
			AddWordToTournamentInvertedIndex(r, wn, id)
		}
	}

	c.Infof("pw: AddToTournamentInvertedIndex end")	
	
}

func FindTeamInvertedIndex(r *http.Request, filter string, value interface{}) *TeamInvertedIndex {
	q := datastore.NewQuery("TeamInvertedIndex").Filter(filter + " =", value).Limit(1)
	
	var t[]*TeamInvertedIndex
	
	if _, err := q.GetAll(appengine.NewContext(r), &t); err == nil && len(t) > 0 {
		return t[0]
	}
	
	return nil
}


func FindTournamentInvertedIndex(r *http.Request, filter string, value interface{}) *TournamentInvertedIndex {
	q := datastore.NewQuery("TournamentInvertedIndex").Filter(filter + " =", value).Limit(1)
	
	var t[]*TournamentInvertedIndex
	
	if _, err := q.GetAll(appengine.NewContext(r), &t); err == nil && len(t) > 0 {
		return t[0]
	}
	
	return nil
}

func KeyByIdTeamInvertedIndex(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)

	return key
}

func KeyByIdTournamentInvertedIndex(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "TournamentInvertedIndex", "", id, nil)

	return key
}


// merge ids in slice of byte with id if it is not already there
// if id is already in the slice return empty string
func mergeIds(teamIds []byte, id int64) string{
	
	strTeamIds := string(teamIds)
	strIds := strings.Split(strTeamIds, " ")
	strId := strconv.FormatInt(id, 10)
	for _, i := range strIds{
		if i == strId{
			return ""
		}
	}
	return strTeamIds + " " + strId
}

// remove id from slice of byte with ids.
func removeFromIds(teamIds []byte, id int64)(string,error){
	strTeamIds := string(teamIds)
	strIds := strings.Split(strTeamIds, " ")
	strId := strconv.FormatInt(id, 10)
	strRet := ""
	for _,val := range strIds{
		if val != strId{
			if len(strRet)==0{
				strRet = val
			}else{
				strRet = strRet + " " + val
			}
		}
	}
	return strRet, nil
}

func TeamInvertedIndexes(r *http.Request, words []string)[]int64{
	c := appengine.NewContext(r)

	strMerge := ""
	for _, w := range words{
		l := ""
		if res := FindTeamInvertedIndex(r, "KeyName", w);res !=nil{
			strTeamIds := string(res.TeamIds)
			if len(l) == 0{
				l = strTeamIds
			}else{
				l = l + " " + strTeamIds
			}
		}
		if len(strMerge) == 0{
			strMerge = l
		}else{
			// build intersection between merge and l
			strMerge = intersect(strMerge,l)
		}
	}
	strIds := strings.Split(strMerge, " ")
	intIds := make([]int64, len(strIds)) 

	i := 0
	for _, w := range strIds{
		if n, err := strconv.ParseInt(w,10,64); err == nil{
			intIds[i] = n
			i = i + 1
		}else{
			c.Errorf("pw: unable to parse %v, error:%v", w, err)
		}
	}
	return intIds	
}

func intersect(a string, b string) string{
	sa := helpers.SetOfStrings(a)
	sb := helpers.SetOfStrings(b)
	intersect := ""
	for _, val := range sa{
		if helpers.SliceContains(sb,val){
			if len(intersect)==0{
				intersect = val
			}else{
				intersect = intersect + " " + val
			}
		}
	}
	return intersect
}




















