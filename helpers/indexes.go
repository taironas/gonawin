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
	"strings"
	"strconv"
)

// merge ids in slice of byte with id if it is not already there
// if id is already in the slice return empty string
func MergeIds(teamIds []byte, id int64) string{
	
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
func RemovefromIds(teamIds []byte, id int64)(string,error){
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
