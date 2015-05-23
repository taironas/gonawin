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

// Package log provide a set of function to log in gonawin app.
//
package log

import (
	"appengine"
)

// information overrides the infof of appengine to be able to centralize the prefix of all logs
func Infof(c appengine.Context, format string, args ...interface{}) {
	c.Infof("gonawin: "+format, args...)
}

// error overrides the infof of appengine to be able to centralize the prefix of all logs
func Errorf(c appengine.Context, format string, args ...interface{}) {
	c.Errorf("gonawin: "+format, args...)
}

// warning overrides the infof of appengine to be able to centralize the prefix of all logs
func Warningf(c appengine.Context, format string, args ...interface{}) {
	c.Warningf("gonawin: "+format, args...)
}
