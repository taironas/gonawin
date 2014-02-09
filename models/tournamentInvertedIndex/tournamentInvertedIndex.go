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

package tournamentinvid

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
)

type TournamentInvertedIndex struct {
	Id            int64
	KeyName       string
	TournamentIds []byte
}

type TournamentInvertedIndexJson struct {
	Id            *int64  `json:",omitempty"`
	KeyName       *string `json:",omitempty"`
	TournamentIds *[]byte `json:",omitempty"`
}

type WordCountTournament struct {
	Count int64
}

// Create a tournament inverted index
func Create(c appengine.Context, name string, tournamentIds string) (*TournamentInvertedIndex, error) {

	id, _, err := datastore.AllocateIDs(c, "TournamentInvertedIndex", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "TournamentInvertedIndex", "", id, nil)

	byteIds := []byte(tournamentIds)
	t := &TournamentInvertedIndex{id, helpers.TrimLower(name), byteIds}

	_, err = datastore.Put(c, key, t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// AddToTournamentInvertedIndex
// Split name by words.
// For each word check if it exist in the Tournament Inverted Index.
// if not create a line with the word as key and team id as value.
func Add(c appengine.Context, name string, id int64) error {

	words := strings.Split(name, " ")
	for _, w := range words {
		log.Infof(c, " AddToTournamentInvertedIndex: Word: %v", w)

		if invId, err := Find(c, "KeyName", w); err != nil {
			return errors.New(fmt.Sprintf(" tournamentinvid.Add, unable to find KeyName=%s: %v", w, err))
		} else if invId == nil {
			log.Infof(c, " create inv id as word does not exist in table")
			Create(c, w, strconv.FormatInt(id, 10))
			log.Infof(c, " create done tournament inv id")
		} else {
			// update row with new info
			log.Infof(c, " update row with new info")
			log.Infof(c, " current info: keyname: %v", invId.KeyName)
			log.Infof(c, " current info: teamIDs: %v", string(invId.TournamentIds))
			k := KeyById(c, invId.Id)

			if newIds := helpers.MergeIds(invId.TournamentIds, id); len(newIds) > 0 {
				log.Infof(c, " current info: new team ids: %v", newIds)
				invId.TournamentIds = []byte(newIds)
				if _, err := datastore.Put(c, k, invId); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// update a team inverted index given an oldname, a new name and an id
func Update(c appengine.Context, oldname string, newname string, id int64) error {

	var err error

	// if word in old and new do nothing
	old_w := strings.Split(oldname, " ")
	new_w := strings.Split(newname, " ")

	// remove id  from words in old name that are not present in new name
	for _, wo := range old_w {
		innew := false
		for _, wn := range new_w {
			if wo == wn {
				innew = true
			}
		}
		if !innew {
			log.Infof(c, " remove: %v", wo)
			err = removeWord(c, wo, id)
			//remove it
		}
	}

	// add all id words in new name
	for _, wn := range new_w {
		inold := false
		for _, wo := range old_w {
			if wo == wn {
				inold = true
			}
		}
		if !inold {
			// add it
			log.Infof(c, " add: %v", wn)
			err = addWord(c, wn, id)
		}
	}

	return err
}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well
func removeWord(c appengine.Context, w string, id int64) error {

	invId, err := Find(c, "KeyName", w)
	if err != nil {
		return errors.New(fmt.Sprintf(" tournamentinvid.removeWord, unable to find KeyName=%s: %v", w, err))
	} else if invId == nil {
		log.Infof(c, " word %v does not exist in Tournament InvertedIndex so nothing to remove", w)
	} else {
		// update row with new info
		k := KeyById(c, invId.Id)

		if newIds, err := helpers.RemovefromIds(invId.TournamentIds, id); err == nil {
			log.Infof(c, " new tournament ids after removal: %v", newIds)
			if len(newIds) == 0 {
				// this entity does not have ids so remove it from the datastore.
				log.Infof(c, " removing key %v from datastore as it is no longer used", k)
				datastore.Delete(c, k)
			} else {
				invId.TournamentIds = []byte(newIds)
				if _, err := datastore.Put(c, k, invId); err != nil {
					log.Errorf(c, " RemoveWordFromTournamentInvertedIndex error on update")
				}
			}
		} else {
			return errors.New(fmt.Sprintf(" unable to remove id from ids: %v", err))
		}
	}

	return nil
}

// add word to tournament inverted index is handled the same way as AddToTournamentInvertedIndex.
// we add this function for clarity
func addWord(c appengine.Context, word string, id int64) error {
	return Add(c, word, id)
}

func Find(c appengine.Context, filter string, value interface{}) (*TournamentInvertedIndex, error) {
	q := datastore.NewQuery("TournamentInvertedIndex").Filter(filter+" =", value).Limit(1)

	var t []*TournamentInvertedIndex

	if _, err := q.GetAll(c, &t); err == nil && len(t) > 0 {
		return t[0], nil
	} else {
		return nil, err
	}
}

// get key pointer to tournament inverted index given an id
func KeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "TournamentInvertedIndex", "", id, nil)

	return key
}

// get an array of ids given a set of words pass as paramenters.
func GetIndexes(c appengine.Context, words []string) ([]int64, error) {
	var err1 error = nil
	strMerge := ""
	for _, w := range words {
		l := ""
		res, err := Find(c, "KeyName", w)
		if err != nil {
			log.Infof(c, "tournamentinvid.GetIndexes, unable to find KeyName=%s: %v", w, err)
			err1 = errors.New(fmt.Sprintf(" tournamentinvid.GetIndexes, unable to find KeyName=%s: %v", w, err))
		} else if res != nil {
			strTournamentIds := string(res.TournamentIds)
			if len(l) == 0 {
				l = strTournamentIds
			} else {
				l = l + " " + strTournamentIds
			}
		}
		if len(strMerge) == 0 {
			strMerge = l
		} else {
			// build intersection between merge and l
			strMerge = helpers.Intersect(strMerge, l)
		}
	}
	// no need to continue if no results were found, just return empty array
	if len(strMerge) == 0 {
		intIds := make([]int64, 0)
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
				log.Infof(c, "tournamentinvid.GetIndexes, unable to parse %v, error:%v", w, err)
			}
		}
	}
	return intIds, err1
}

// increment word count for tournaments
func incrementWordCountTournament(c appengine.Context, key *datastore.Key) (int64, error) {
	var x WordCountTournament
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// decrement word count for tournaments
func decrementWordCountTournament(c appengine.Context, key *datastore.Key) (int64, error) {
	var x WordCountTournament
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// get current word count in tournaments
func GetWordCount(c appengine.Context) (int64, error) {
	key := datastore.NewKey(c, "WordCountTournament", "singleton", 0, nil)
	var x WordCountTournament
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

// get frequency of a given word base on tournamentinverted index data
func GetTournamentFrequencyForWord(c appengine.Context, word string) (int64, error) {

	if invId, err := Find(c, "KeyName", word); err != nil {
		return 0, errors.New(fmt.Sprintf(" tournamentinvid.GetTournamentFrequencyForWord, unable to find KeyName=%s: %v", word, err))
	} else if invId == nil {
		return 0, nil
	} else {
		return int64(len(strings.Split(string(invId.TournamentIds), " "))), nil
	}
}
