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

// UserInvertedIndex holds informations needed for User indexing.
//
type UserInvertedIndex struct {
	Id      int64
	KeyName string
	UserIds []byte
}

// UserInvertedIndexJSON is the JSON representation of UserInvertedIndex.
//
type UserInvertedIndexJSON struct {
	Id      *int64  `json:",omitempty"`
	KeyName *string `json:",omitempty"`
	UserIds *[]byte `json:",omitempty"`
}

// WordCountUser holds a word counter for User entities.
//
type WordCountUser struct {
	Count int64
}

// CreateUserInvertedIndex creates a userinvertedindex entity given a word and a list of ids as a string.
//
func CreateUserInvertedIndex(c appengine.Context, word string, ids string) (*UserInvertedIndex, error) {

	id, _, err := datastore.AllocateIDs(c, "UserInvertedIndex", nil, 1)
	if err != nil {
		return nil, err
	}

	key := datastore.NewKey(c, "UserInvertedIndex", "", id, nil)

	byteIds := []byte(ids)
	u := &UserInvertedIndex{id, word, byteIds}

	_, err = datastore.Put(c, key, u)
	if err != nil {
		return nil, err
	}
	// increment counter
	errIncrement := datastore.RunInTransaction(c, func(c appengine.Context) error {
		var err1 error
		_, err1 = incrementWordCountUser(c, datastore.NewKey(c, "WordCountUser", "singleton", 0, nil))
		return err1
	}, nil)
	if errIncrement != nil {
		log.Errorf(c, " Error incrementing WordCountUser")
	}

	return u, err
}

// AddToUserInvertedIndex add a name the user inverted index entity.
//
// We do this by spliting the name in words (split by spaces),
// for each word we check if it already exists a user inverted index entity.
// If it does not yet exist, we create an entity with the word as key and user id as value.
//
func AddToUserInvertedIndex(c appengine.Context, name string, id int64) error {
	words := strings.Split(name, " ")
	for _, w := range words {

		if invID, err := FindUserInvertedIndex(c, "KeyName", w); err != nil {
			return fmt.Errorf(" userinvid.Add, unable to find KeyName=%s: %v", w, err)
		} else if invID == nil {
			CreateUserInvertedIndex(c, w, strconv.FormatInt(id, 10))
		} else {
			// update row with new info
			k := UserInvertedIndexKeyByID(c, invID.Id)

			if newIds := helpers.MergeIds(invID.UserIds, id); len(newIds) > 0 {
				invID.UserIds = []byte(newIds)
				if _, err := datastore.Put(c, k, invID); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// UpdateUserInvertedIndex updates the User inverted index.
// From the old user name and the new user name we handle the removal of the words that are no longer present and the addition of new words.
//
func UpdateUserInvertedIndex(c appengine.Context, oldname string, newname string, id int64) error {

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
			err = userInvertedIndexRemoveWord(c, wo, id)
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
			err = userInvertedIndexAddWord(c, wn, id)
		}
	}
	return err
}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well.
//
func userInvertedIndexRemoveWord(c appengine.Context, w string, id int64) error {

	invID, err := FindUserInvertedIndex(c, "KeyName", w)
	if err != nil {
		return fmt.Errorf(" userinvid.removeWord, unable to find KeyName=%s: %v", w, err)
	} else if invID != nil {
		// update row with new info
		k := UserInvertedIndexKeyByID(c, invID.Id)

		if newIds, err := helpers.RemovefromIds(invID.UserIds, id); err == nil {
			if len(newIds) == 0 {
				// this entity does not have ids so remove it from the datastore.
				datastore.Delete(c, k)
				// decrement word counter
				errDec := datastore.RunInTransaction(c, func(c appengine.Context) error {
					var err1 error
					_, err1 = decrementWordCountUser(c, datastore.NewKey(c, "WordCountUser", "singleton", 0, nil))
					return err1
				}, nil)
				if errDec != nil {
					return fmt.Errorf(" Error decrementing WordCountUser: %v", errDec)
				}
			} else {
				invID.UserIds = []byte(newIds)
				if _, err1 := datastore.Put(c, k, invID); err1 != nil {
					return fmt.Errorf(" RemoveWordFromUserInvertedIndex error on update: %v", err)
				}
			}
		} else {
			return fmt.Errorf(" unable to remove id from ids: %v", err)
		}
	}
	return nil
}

// add word to user inverted index is handled the same way as AddToUserInvertedIndex.
// we add this function for clarity.
//
func userInvertedIndexAddWord(c appengine.Context, word string, id int64) error {
	return AddToUserInvertedIndex(c, word, id)
}

// FindUserInvertedIndex looks for an entity in the datastore given a filter and a value .
//
func FindUserInvertedIndex(c appengine.Context, filter string, value interface{}) (*UserInvertedIndex, error) {
	q := datastore.NewQuery("UserInvertedIndex").Filter(filter+" =", value).Limit(1)

	var u []*UserInvertedIndex

	if _, err := q.GetAll(c, &u); err != nil || len(u) <= 0 {
		return nil, err
	}

	return u[0], nil
}

// UserInvertedIndexKeyByID returns, given an id, a pointer to the corresponding key of a user inverted index entity if found.
//
func UserInvertedIndexKeyByID(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "UserInvertedIndex", "", id, nil)
	return key
}

// GetUserInvertedIndexes returns, Given an array of words, an array of indexes that correspond to the user ids of the users that use these words.
//
func GetUserInvertedIndexes(c appengine.Context, words []string) ([]int64, error) {
	var err1 error
	strMerge := ""
	for _, w := range words {
		l := ""
		res, err := FindUserInvertedIndex(c, "KeyName", w)
		if err != nil {
			log.Errorf(c, "userinvid.GetIndexes, unable to find KeyName=%s: %v", w, err)
			err1 = fmt.Errorf(" userinvid.GetIndexes, unable to find KeyName=%s: %v", w, err)
		} else if res != nil {
			strUserIds := string(res.UserIds)
			if len(l) == 0 {
				l = strUserIds
			} else {
				l = l + " " + strUserIds
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
			}
		}
	}
	return intIds, err1
}

// increment the word count user counter.
//
func incrementWordCountUser(c appengine.Context, key *datastore.Key) (int64, error) {
	var x WordCountUser
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count++
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// decrement the word count user counter.
//
func decrementWordCountUser(c appengine.Context, key *datastore.Key) (int64, error) {
	var x WordCountUser
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	x.Count--
	if _, err := datastore.Put(c, key, &x); err != nil {
		return 0, err
	}
	return x.Count, nil
}

// UserInvertedIndexGetWordCount returns the current number of words on user names.
//
func UserInvertedIndexGetWordCount(c appengine.Context) (int64, error) {
	key := datastore.NewKey(c, "WordCountUser", "singleton", 0, nil)
	var x WordCountUser
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

// GetUserFrequencyForWord gets the number of users that have 'word' in their name.
//
func GetUserFrequencyForWord(c appengine.Context, word string) (int64, error) {

	if invID, err := FindUserInvertedIndex(c, "KeyName", word); err != nil {
		return 0, fmt.Errorf(" userinvid.GetUserFrequencyForWord, unable to find KeyName=%s: %v", word, err)
	} else if invID == nil {
		return 0, nil
	} else {
		return int64(len(strings.Split(string(invID.UserIds), " "))), nil
	}
}
