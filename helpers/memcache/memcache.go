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

// Package memcache provides a set of functions to use memcache in gonawin.
//
package memcache

import (
	"fmt"

	"appengine"
	"appengine/memcache"

	"github.com/taironas/gonawin/helpers/log"
)

// Set sets a key value pair to memcache.
//
func Set(c appengine.Context, key string, value interface{}) error {

	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case int64:
		bytes = []byte(fmt.Sprintf("%d", v))
	}

	item := &memcache.Item{
		Key:   key,
		Value: bytes,
	}
	// Set the item, unconditionally
	if err := memcache.Set(c, item); err != nil {
		log.Errorf(c, " error setting item: %v", err)
		return err
	}

	return nil
}

// Get gets a value from memcache with respect to a key string.
//
func Get(c appengine.Context, key string) (interface{}, error) {
	// Get the item from the memcache
	item, err := memcache.Get(c, key)

	if err != nil {
		log.Errorf(c, " error getting item: %v", err)
		return nil, err
	}

	return item.Value, err
}

// Delete deletes a key from memcache.
//
func Delete(c appengine.Context, key string) error {
	return memcache.Delete(c, key)
}
