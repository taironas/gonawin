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

package teaminvid

import (
	"errors"
	"fmt"
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

type WordCountTeam struct{
	Count int64
}

func Create(r *http.Request, word string, teamIds string) (*TeamInvertedIndex, error) {
	c := appengine.NewContext(r)
	
	id, _, err := datastore.AllocateIDs(c, "TeamInvertedIndex", nil, 1)
	if err != nil {
		return nil, err
	}
	
	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)
	
	byteIds := []byte(teamIds)
	t := &TeamInvertedIndex{ id, word, byteIds }

	_, err = datastore.Put(c, key, t)
	if err != nil {
		return nil, err
	}
	// increment counter
	errIncrement := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var err1 error
		_, err1 = incrementWordCountTeam(c, datastore.NewKey(c, "WordCountTeam", "singleton", 0, nil))
		return err1
	}, nil)
	if errIncrement != nil {
		c.Errorf("pw: Error incrementing WordCountTeam")
	}
	
	return t, err
}

// AddToTeamInvertedIndex
// Split name by words.
// For each word check if it exist in the Team Inverted Index. 
// if not create a line with the word as key and team id as value.
func Add(r *http.Request, name string, id int64) error {
	c := appengine.NewContext(r)
	
	words := strings.Split(name, " ")
	for _, w:= range words{
		c.Infof("pw: AddToTeamInvertedIndex: Word: %v", w)

		if invId, err := Find(r, "KeyName", w); err != nil {
			return errors.New(fmt.Sprintf("pw: teaminvid.Add, unable to find KeyName=%s: %v", w, err))
		} else if	invId == nil {
			c.Infof("pw: create inv id as word does not exist in table")
			Create(r, w, strconv.FormatInt(id, 10))
		} else {
			// update row with new info
			c.Infof("pw: update row with new info")
			c.Infof("pw: current info: keyname: %v", invId.KeyName)
			c.Infof("pw: current info: teamIDs: %v", string(invId.TeamIds))
			k := KeyById(r, invId.Id)

			if newIds := helpers.MergeIds(invId.TeamIds, id);len(newIds) > 0 {
				c.Infof("pw: current info: new team ids: %v", newIds)
				invId.TeamIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, invId);err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// from the oldname and the new name we handle removal of words no longer present.
// the addition of new words.
func Update(r *http.Request, oldname string, newname string, id int64) error {
	c := appengine.NewContext(r)
	
	var err error

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
			err = removeWord(r, wo, id)
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
			c.Infof("pw: add: %v", wn)
			err = addWord(r, wn, id)
		}
	}
	
	return err
}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well
func removeWord(r *http.Request, w string, id int64) error {
	c := appengine.NewContext(r)

	invId, err := Find(r, "KeyName", w)
	if err != nil {
		return errors.New(fmt.Sprintf("pw: teaminvid.removeWord, unable to find KeyName=%s: %v", w, err))
	} else if invId == nil {
		c.Infof("pw: word %v does not exist in Team InvertedIndex so nothing to remove", w)
	} else {
		// update row with new info
		k := KeyById(r, invId.Id)

		if newIds, err := helpers.RemovefromIds(invId.TeamIds, id); err == nil{
			c.Infof("pw: new team ids after removal: %v", newIds)
			if len(newIds) == 0 {
				// this entity does not have ids so remove it from the datastore.
				c.Infof("pw: removing key %v from datastore as it is no longer used", k)
				datastore.Delete(c, k)
				// decrement word counter
				errDec := datastore.RunInTransaction(c, func(c appengine.Context) error {
					var err1 error
					_, err1 = decrementWordCountTeam(c, datastore.NewKey(c, "WordCountTeam", "singleton", 0, nil))
					return err1
				}, nil)
				if errDec != nil {
					return errors.New(fmt.Sprintf("pw: Error decrementing WordCountTeam: %v", errDec))
				}
			} else {
				invId.TeamIds  = []byte(newIds)
				if _, err := datastore.Put(c, k, invId);err != nil{
					return errors.New(fmt.Sprintf("pw: RemoveWordFromTeamInvertedIndex error on update: %v", err))
				}
			}
		} else {
			return errors.New(fmt.Sprintf("pw: unable to remove id from ids: %v", err))
		}
	}
	
	return nil
}


// add word to team inverted index is handled the same way as AddToTeamInvertedIndex.
// we add this function for clarity
func addWord(r *http.Request, word string, id int64) error {
	return Add(r, word, id)
}

func Find(r *http.Request, filter string, value interface{}) (*TeamInvertedIndex, error) {
	q := datastore.NewQuery("TeamInvertedIndex").Filter(filter + " =", value).Limit(1)
	
	var t[]*TeamInvertedIndex
	
	if _, err := q.GetAll(appengine.NewContext(r), &t); err == nil && len(t) > 0 {
		return t[0], nil
	} else {
		return nil, err
	}
}

func KeyById(r *http.Request, id int64) (*datastore.Key) {
	c := appengine.NewContext(r)

	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)

	return key
}

func GetIndexes(r *http.Request, words []string) ([]int64, error) {
	var err1 error = nil
	strMerge := ""
	
	for _, w := range words {
		l := ""
		
		res, err := Find(r, "KeyName", w)
		if err != nil {
			err1 = errors.New(fmt.Sprintf("pw: teaminvid.GetIndexes, unable to find KeyName=%s: %v", w, err))
		} else if res !=nil {
			strTeamIds := string(res.TeamIds)
			if len(l) == 0 {
				l = strTeamIds
			} else {
				l = l + " " + strTeamIds
			}
		}
		if len(strMerge) == 0 {
			strMerge = l
		} else {
			// build intersection between merge and l
			strMerge = helpers.Intersect(strMerge,l)
		}
	}
	strIds := strings.Split(strMerge, " ")
	intIds := make([]int64, len(strIds)) 

	i := 0
	for _, w := range strIds {
		if n, err := strconv.ParseInt(w, 10, 64); err == nil {
			intIds[i] = n
			i = i + 1
		} else {
			err1 = errors.New(fmt.Sprintf("pw: teaminvid.GetIndexes, unable to parse %v, error:%v", w, err))
		}
	}
	
	return intIds, err1	
}

func incrementWordCountTeam(c appengine.Context, key *datastore.Key) (int64, error) {
	var x WordCountTeam
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

func decrementWordCountTeam(c appengine.Context, key *datastore.Key) (int64, error) {
	var x WordCountTeam
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

func GetWordCount(c appengine.Context)(int64, error) {
	key := datastore.NewKey(c, "WordCountTeam", "singleton", 0, nil)
	var x WordCountTeam
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

func GetTeamFrequencyForWord(r *http.Request, word string) (int64, error) {
	
	if invId, err := Find(r, "KeyName", word); err != nil {
		return 0, errors.New(fmt.Sprintf("pw: teaminvid.GetTeamFrequencyForWord, unable to find KeyName=%s: %v", word, err))
	} else if invId == nil {
		return 0, nil
	} else {
		return int64(len(strings.Split(string(invId.TeamIds), " "))), nil
	}
}
