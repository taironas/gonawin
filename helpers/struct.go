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

package helpers

import (
	"reflect"
)

type arrayOfStrings []string

func (a arrayOfStrings) Contains(s string) bool {
	for _, e := range a {
		if e == s {
			return true
		}
	}
	return false
}

func CopyToPtrBasedStructGeneric(tSrc interface{}, tDest interface{}) {
	s1 := reflect.ValueOf(tSrc).Elem()
	s2 := reflect.ValueOf(tDest).Elem()
	for i := 0; i < s1.NumField(); i++ {
		f1 := s1.Field(i)
		f2 := s2.Field(i)
		if f2.CanSet() {
			s2.Field(i).Set(f1.Addr())
		}
	}
}

func KeepFields(t interface{}, fieldsToKeep arrayOfStrings) {
	s := reflect.ValueOf(t).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !fieldsToKeep.Contains(typeOfT.Field(i).Name) && f.CanSet() {
			s.Field(i).Set(reflect.Zero(typeOfT.Field(i).Type))
		}
	}
}
