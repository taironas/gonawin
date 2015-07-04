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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"appengine"
	"appengine/datastore"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
)

type UserInvertedIndex struct {
	Id      int64
	KeyName string
	UserIds []byte
}

type UserInvertedIndexJson struct {
	Id      *int64  `json:",omitempty"`
	KeyName *string `json:",omitempty"`
	UserIds *[]byte `json:",omitempty"`
}

type WordCountUser struct {
	Count int64
}

// Create a userinvertedindex entity given a word and a list of ids as a string.
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
	desc := "AddToUserInvertedIndex: "
	words := strings.Split(name, " ")
	for _, w := range words {
		log.Infof(c, "%s Word: %v", desc, w)

		if invId, err := FindUserInvertedIndex(c, "KeyName", w); err != nil {
			return errors.New(fmt.Sprintf(" userinvid.Add, unable to find KeyName=%s: %v", w, err))
		} else if invId == nil {
			log.Infof(c, " create inv id as word does not exist in table")
			CreateUserInvertedIndex(c, w, strconv.FormatInt(id, 10))
			log.Infof(c, " create done user inv id")
		} else {
			// update row with new info
			log.Infof(c, " update row with new info")
			log.Infof(c, " current info: keyname: %v", invId.KeyName)
			log.Infof(c, " current info: userIds: %v", string(invId.UserIds))
			k := UserInvertedIndexKeyById(c, invId.Id)

			if newIds := helpers.MergeIds(invId.UserIds, id); len(newIds) > 0 {
				log.Infof(c, " current info: new user ids: %v", newIds)
				invId.UserIds = []byte(newIds)
				if _, err := datastore.Put(c, k, invId); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// From the old user name and the new user name we handle the removal of the words that are no longer present and the addition of new words.
func UpdateUserInvertedIndex(c appengine.Context, oldname string, newname string, id int64) error {

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
			err = userInvertedIndexRemoveWord(c, wo, id)
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
		if !inold && (len(wn) > 0) {
			log.Infof(c, " add: %v", wn)
			err = userInvertedIndexAddWord(c, wn, id)
		}
	}
	return err
}

// if the removal of the id makes the entity useless (no more ids in it)
// we will remove the entity as well
func userInvertedIndexRemoveWord(c appengine.Context, w string, id int64) error {

	invId, err := FindUserInvertedIndex(c, "KeyName", w)
	if err != nil {
		return errors.New(fmt.Sprintf(" userinvid.removeWord, unable to find KeyName=%s: %v", w, err))
	} else if invId == nil {
		log.Infof(c, " word %v does not exist in User InvertedIndex so nothing to remove", w)
	} else {
		// update row with new info
		k := UserInvertedIndexKeyById(c, invId.Id)

		if newIds, err := helpers.RemovefromIds(invId.UserIds, id); err == nil {
			log.Infof(c, " new user ids after removal: %v", newIds)
			if len(newIds) == 0 {
				// this entity does not have ids so remove it from the datastore.
				log.Infof(c, " removing key %v from datastore as it is no longer used", k)
				datastore.Delete(c, k)
				// decrement word counter
				errDec := datastore.RunInTransaction(c, func(c appengine.Context) error {
					var err1 error
					_, err1 = decrementWordCountUser(c, datastore.NewKey(c, "WordCountUser", "singleton", 0, nil))
					return err1
				}, nil)
				if errDec != nil {
					return errors.New(fmt.Sprintf(" Error decrementing WordCountUser: %v", errDec))
				}
			} else {
				invId.UserIds = []byte(newIds)
				if _, err1 := datastore.Put(c, k, invId); err1 != nil {
					return errors.New(fmt.Sprintf(" RemoveWordFromUserInvertedIndex error on update: %v", err))
				}
			}
		} else {
			return errors.New(fmt.Sprintf(" unable to remove id from ids: %v", err))
		}
	}
	return nil
}

// add word to user inverted index is handled the same way as AddToUserInvertedIndex.
// we add this function for clarity
func userInvertedIndexAddWord(c appengine.Context, word string, id int64) error {
	return AddToUserInvertedIndex(c, word, id)
}

// given a filter and a value look for an entity in the datastore
func FindUserInvertedIndex(c appengine.Context, filter string, value interface{}) (*UserInvertedIndex, error) {
	q := datastore.NewQuery("UserInvertedIndex").Filter(filter+" =", value).Limit(1)

	var u []*UserInvertedIndex

	if _, err := q.GetAll(c, &u); err == nil && len(u) > 0 {
		return u[0], nil
	} else {
		return nil, err
	}
}

// Given an id returns a pointer to the corresponding key of a user inverted index entity if found.
func UserInvertedIndexKeyById(c appengine.Context, id int64) *datastore.Key {

	key := datastore.NewKey(c, "UserInvertedIndex", "", id, nil)
	return key
}

// Given an array of words, return an array of indexes that correspond to the user ids of the users that use these words.
func GetUserInvertedIndexes(c appengine.Context, words []string) ([]int64, error) {
	var err1 error = nil
	strMerge := ""
	for _, w := range words {
		l := ""
		res, err := FindUserInvertedIndex(c, "KeyName", w)
		if err != nil {
			log.Infof(c, "userinvid.GetIndexes, unable to find KeyName=%s: %v", w, err)
			err1 = errors.New(fmt.Sprintf(" userinvid.GetIndexes, unable to find KeyName=%s: %v", w, err))
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
				log.Infof(c, "userinvid.GetIndexes, unable to parse %v, error:%v", w, err)
			}
		}
	}
	return intIds, err1
}

// increment the word count user counter.
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

// Returns the current number of words on user names.
func UserInvertedIndexGetWordCount(c appengine.Context) (int64, error) {
	key := datastore.NewKey(c, "WordCountUser", "singleton", 0, nil)
	var x WordCountUser
	if err := datastore.Get(c, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
		return 0, err
	}
	return x.Count, nil
}

// Get the number of users that have 'word' in their name.
func GetUserFrequencyForWord(c appengine.Context, word string) (int64, error) {

	if invId, err := FindUserInvertedIndex(c, "KeyName", word); err != nil {
		return 0, errors.New(fmt.Sprintf(" userinvid.GetUserFrequencyForWord, unable to find KeyName=%s: %v", word, err))
	} else if invId == nil {
		return 0, nil
	} else {
		return int64(len(strings.Split(string(invId.UserIds), " "))), nil
	}
}
