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
	"strings"

	"appengine"
	"appengine/datastore"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
)

// TeamInvertedIndex holds informations needed for Team indexing.
//
type TeamInvertedIndex struct {
	Id      int64
	KeyName string
	TeamIds []byte
}

// TeamInvertedIndexJSON is the JSON representation of TeamInvertedIndex.
//
type TeamInvertedIndexJSON struct {
	Id      *int64  `json:",omitempty"`
	KeyName *string `json:",omitempty"`
	TeamIds *[]byte `json:",omitempty"`
}

// WordCountTeam holds a word counter for Team entities.
//
type WordCountTeam struct {
	Count int64
}

// CreateTeamInvertedIndex creates a teaminvertedindex entity given a word and a list of ids as a string.
//
func CreateTeamInvertedIndex(c appengine.Context, word string, teamIds string) (*TeamInvertedIndex, error) {

	id, _, err := datastore.AllocateIDs(c, "TeamInvertedIndex", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)

	byteIds := []byte(teamIds)
	t := &TeamInvertedIndex{id, word, byteIds}

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
		log.Errorf(c, " Error incrementing WordCountTeam")
	}

	return t, err
}

// AddToTeamInvertedIndex adds name to team inverted index entity.
//
// We do this by spliting the name in words (split by spaces),
// for each word we check if it already exists a team inverted index entity.
// If it does not yet exist, we create an entity with the word as key and team id as value.
//
func AddToTeamInvertedIndex(c appengine.Context, name string, id int64) error {

	words := strings.Split(name, " ")
	for _, w := range words {

		if invID, err := FindTeamInvertedIndex(c, "KeyName", w); err != nil {
			return fmt.Errorf(" teaminvid.Add, unable to find KeyName=%s: %v", w, err)
		} else if invID == nil {
			CreateTeamInvertedIndex(c, w, strconv.FormatInt(id, 10))
		} else {
			// update row with new info
			k := TeamInvertedIndexKeyByID(c, invID.Id)

			if newIds := helpers.MergeIds(invID.TeamIds, id); len(newIds) > 0 {
				invID.TeamIds = []byte(newIds)
				if _, err := datastore.Put(c, k, invID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateTeamInvertedIndex updates the Team inverted index.
// From the old team name and the new team name we handle the removal of the words that are no longer present and the addition of new words.
//
func UpdateTeamInvertedIndex(c appengine.Context, oldname string, newname string, id int64) error {

	var err error

	// if word in old and new do nothing
	oldW := strings.Split(oldname, " ")
	newW := strings.Split(newname, " ")

	// remove id  from words in old name that are not present in new name
	for _, wo := range oldW {
		innew := false
		for _, wn := range newW {
			if wo == wn {
				innew = true
			}
		}
		if !innew {
			err = teamInvertedIndexRemoveWord(c, wo, id)
		}
	}

	// add all id words in new name
	for _, wn := range newW {
		inold := false
		for _, wo := range oldW {
			if wo == wn {
				inold = true
			}
		}
		if !inold && (len(wn) > 0) {
			err = teamInvertedIndexAddWord(c, wn, id)
		}
	}

	return err
}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well.
//
func teamInvertedIndexRemoveWord(c appengine.Context, w string, id int64) error {

	invID, err := FindTeamInvertedIndex(c, "KeyName", w)
	if err != nil {
		return fmt.Errorf(" teaminvid.removeWord, unable to find KeyName=%s: %v", w, err)
	} else if invID != nil {
		// update row with new info
		k := TeamInvertedIndexKeyByID(c, invID.Id)

		if newIds, err := helpers.RemovefromIds(invID.TeamIds, id); err == nil {
			if len(newIds) == 0 {
				// this entity does not have ids so remove it from the datastore.
				datastore.Delete(c, k)
				// decrement word counter
				errDec := datastore.RunInTransaction(c, func(c appengine.Context) error {
					var err1 error
					_, err1 = decrementWordCountTeam(c, datastore.NewKey(c, "WordCountTeam", "singleton", 0, nil))
					return err1
				}, nil)
				if errDec != nil {
					return fmt.Errorf(" Error decrementing WordCountTeam: %v", errDec)
				}
			} else {
				invID.TeamIds = []byte(newIds)
				if _, err1 := datastore.Put(c, k, invID); err1 != nil {
					return fmt.Errorf(" RemoveWordFromTeamInvertedIndex error on update: %v", err)
				}
			}
		} else {
			return fmt.Errorf(" unable to remove id from ids: %v", err)
		}
	}

	return nil
}

// add word to team inverted index is handled the same way as AddToTeamInvertedIndex.
// we add this function for clarity
func teamInvertedIndexAddWord(c appengine.Context, word string, id int64) error {
	return AddToTeamInvertedIndex(c, word, id)
}

// FindTeamInvertedIndex looks for an entity in the datastore given a filter and a value .
//
func FindTeamInvertedIndex(c appengine.Context, filter string, value interface{}) (*TeamInvertedIndex, error) {
	q := datastore.NewQuery("TeamInvertedIndex").Filter(filter+" =", value).Limit(1)

	var t []*TeamInvertedIndex

	if _, err := q.GetAll(c, &t); err != nil || len(t) <= 0 {
		return nil, err
	}

	return t[0], nil
}

// TeamInvertedIndexKeyByID returns, given an id, a pointer to the corresponding key of a team inverted index entity if found.
//
func TeamInvertedIndexKeyByID(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "TeamInvertedIndex", "", id, nil)
	return key
}

// GetTeamInvertedIndexes returns, Given an array of words, an array of indexes that correspond to the team ids of the teams that use these words.
//
func GetTeamInvertedIndexes(c appengine.Context, words []string) ([]int64, error) {
	var err1 error
	strMerge := ""
	for _, w := range words {
		l := ""
		res, err := FindTeamInvertedIndex(c, "KeyName", w)
		if err != nil {
			log.Errorf(c, "teaminvid.GetIndexes, unable to find KeyName=%s: %v", w, err)
			err1 = fmt.Errorf(" teaminvid.GetIndexes, unable to find KeyName=%s: %v", w, err)
		} else if res != nil {
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
			strMerge = helpers.Intersect(strMerge, l)
		}
	}
	// no need to continue if no results were found, just return emtpy array
	if len(strMerge) == 0 {
		var intIds []int64
		return intIds, err1
	}
	strIds := strings.Split(strMerge, " ")
	intIds := make([]int64, len(strIds))
	i := 0
	for _, w := range strIds {
		if len(w) > 0 {
			if n, err := strconv.ParseInt(w, 10, 64); err == nil {
				intIds[i] = n
				i = i + 1
			} else {
				log.Errorf(c, "teaminvid.GetIndexes, unable to parse %v, error:%v", w, err)
			}
		}
	}
	return intIds, err1
}

// increment the word count team counter
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

// decrement the word count team counter
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

// TeamInvertedIndexGetWordCount returns the current number of words on team names.
//
func TeamInvertedIndexGetWordCount(c appengine.Context) (int64, error) {
	key := datastore.NewKey(c, "WordCountTeam", "singleton", 0, nil)
	var x WordCountTeam
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

// GetTeamFrequencyForWord gets the number of teams that have 'word' in their name.
//
func GetTeamFrequencyForWord(c appengine.Context, word string) (int64, error) {

	if invID, err := FindTeamInvertedIndex(c, "KeyName", word); err != nil {
		return 0, fmt.Errorf(" teaminvid.GetTeamFrequencyForWord, unable to find KeyName=%s: %v", word, err)
	} else if invID == nil {
		return 0, nil
	} else {
		return int64(len(strings.Split(string(invID.TeamIds), " "))), nil
	}
}
