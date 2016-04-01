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

package helpers

import (
	"io"
	"net/http"
)

// Error404 writes a 404 not found in ResponseWriter.
//
func Error404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "404: Not Found")
}

// inspired by https://github.com/campoy/todo

// BadRequest is handled by setting the status code in the reply to StatusBadRequest.
type BadRequest struct {
	Err error
}

// Implementation of error, returns string error on Err structure.
func (e *BadRequest) Error() string {
	return e.Err.Error()
}

// NotFound is handled by setting the status code in the reply to StatusNotFound.
type NotFound struct {
	Err error
}

// Implementation of error, returns string error on Err structure.
func (e *NotFound) Error() string {
	return e.Err.Error()
}

// Forbidden is handled by setting the status code in the reply to StatusForbidden.
type Forbidden struct {
	Err error
}

// Implementation of error, returns string error on Err structure.
func (e *Forbidden) Error() string {
	return e.Err.Error()
}

// Unauthorized is handled by setting the status code in the reply to StatusUnauthorized.
type Unauthorized struct {
	Err error
}

// Implementation of error, returns string error on Err structure.
func (e *Unauthorized) Error() string {
	return e.Err.Error()
}

// InternalServerError is handled by setting the status code in the reply to StatusInternalServerError.
type InternalServerError struct {
	Err error
}

// Implementation of error, returns string error on Err structure.
func (e *InternalServerError) Error() string {
	return e.Err.Error()
}
