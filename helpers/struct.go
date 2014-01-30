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

// from two structures source and destination.
// structures defined as follows:
// both structures have the same fields
// destination structure field types are pointers of the source types
// example:
// type TA struct {
// 	Field1      string
// 	Field2      string
// 	IsSomething bool
// }
// type TB struct {
// 	Field1      *string `json:",omitempty"`
// 	Field2      *string `json:",omitempty"`
// 	IsSomething *bool   `json:",omitempty"`
// }
// use source structure to build destination structure
func CopyToPointerStructure(tSrc interface{}, tDest interface{}) {
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

// Set to nil all fields of t structure not present in array
// we supose that t is a structure of pointer types and fieldsToKeep is an array of the fields you wish to keep.
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

// from two structures, source and destination.
// structures defined as follows:
// both structures have the same fields
// destination structure field types are pointers of the source types
// example:
// type TA struct {
// 	Field1      string
// 	Field2      string
// 	IsSomething bool
// }
// type TB struct {
// 	Field1      *string `json:",omitempty"`
// 	Field2      *string `json:",omitempty"`
// 	IsSomething *bool   `json:",omitempty"`
// }
// use source structure to build destination structure
// and at the same time set to nil all fields of Destination structure not present in array fieldsToKeep
// we supose that pDest is a structure of pointer types and fieldsToKeep is an array of the fields you wish to keep.
func InitPointerStructure(pSrc interface{}, pDest interface{}, fieldsToKeep arrayOfStrings) {
	s1 := reflect.ValueOf(pSrc).Elem()
	s2 := reflect.ValueOf(pDest).Elem()
	for i := 0; i < s1.NumField(); i++ {
		f1 := s1.Field(i)
		f2 := s2.Field(i)
		if f2.CanSet() {
			if fieldsToKeep.Contains(s2.Type().Field(i).Name) {
				s2.Field(i).Set(f1.Addr())
			} else {
				s2.Field(i).Set(reflect.Zero(s2.Type().Field(i).Type))
			}
		}
	}
}

// works like InitPointerStructure but for arrays
// source array has values
func TransformFromArrayOfPointers(pArraySrc interface{}, pArrayDest interface{}, fieldsToKeep arrayOfStrings) {
	arraySrc := reflect.ValueOf(pArraySrc).Elem()
	arrayDest := reflect.ValueOf(pArrayDest).Elem()
	for i := 0; i < arraySrc.Len(); i++ {
		src := arraySrc.Index(i).Elem() // as we are working with array of pointers get true value
		dest := arrayDest.Index(i)
		for j := 0; j < src.NumField(); j++ {
			srcField := src.Field(j)
			destField := dest.Field(j)
			if destField.CanSet() {
				if fieldsToKeep.Contains(dest.Type().Field(j).Name) {
					destField.Set(srcField.Addr())
				} else {
					destField.Set(reflect.Zero(dest.Type().Field(j).Type))
				}
			}
		}
	}
}

// works like InitPointerStructure but for arrays
// source array has pointers
func TransformFromArrayOfValues(pArraySrc interface{}, pArrayDest interface{}, fieldsToKeep arrayOfStrings) {
	arraySrc := reflect.ValueOf(pArraySrc).Elem()
	arrayDest := reflect.ValueOf(pArrayDest).Elem()
	for i := 0; i < arraySrc.Len(); i++ {
		src := arraySrc.Index(i)
		dest := arrayDest.Index(i)
		for j := 0; j < src.NumField(); j++ {
			srcField := src.Field(j)
			destField := dest.Field(j)
			if destField.CanSet() {
				if fieldsToKeep.Contains(dest.Type().Field(j).Name) {
					destField.Set(srcField.Addr())
				} else {
					destField.Set(reflect.Zero(dest.Type().Field(j).Type))
				}
			}
		}
	}
}
