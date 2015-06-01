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

package tasks

import (
	"encoding/json"
	"errors"
	"net/http"

	"appengine"
	"appengine/datastore"
	"appengine/mail"

	"github.com/santiaago/gonawin/helpers"
	"github.com/santiaago/gonawin/helpers/log"
)

// Invite task handler, use it to send an invitation via email.
//
func Invite(w http.ResponseWriter, r *http.Request) error {

	c := appengine.NewContext(r)
	desc := "Task queue - Invite Handler:"

	if r.Method != "POST" {
		log.Infof(c, "%s something went wrong...")
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	log.Infof(c, "%s processing...", desc)

	err := datastore.RunInTransaction(c, func(c appengine.Context) error {

		log.Infof(c, "%s reading data...", desc)

		msg := buildMessage(c, desc, r)

		if err := mail.Send(c, msg); err != nil {
			log.Errorf(c, "%s: couldn't send email: %v", desc, err)
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInviteEmailCannotSend)}
		}

		log.Infof(c, "%s task done!", desc)
		return nil
	}, nil)

	if err != nil {
		c.Errorf("%s error: %v", desc, err)
		log.Infof(c, "%s something went wrong...", desc)
		return err
	}

	return nil
}

func buildMessage(c appengine.Context, desc string, r *http.Request) *mail.Message {

	emailBlob := []byte(r.FormValue("email"))
	nameBlob := []byte(r.FormValue("name"))
	bodyBlob := []byte(r.FormValue("body"))

	email := decode(c, desc, emailBlob)
	name := decode(c, desc, nameBlob)
	body := decode(c, desc, bodyBlob)

	log.Infof(c, "%s value of email: %v", desc, email)
	log.Infof(c, "%s value of name: %v", desc, name)
	log.Infof(c, "%s value of body: %v", desc, body)

	log.Infof(c, "%s crunching data...", desc)

	return &mail.Message{
		Sender:  "No Reply gonawin <no-reply@gonawin.com>",
		To:      []string{email},
		Subject: name + " wants you to join Gonawin!",
		Body:    body,
	}
}

func decode(c appengine.Context, desc string, blob []byte) (v string) {
	if err := json.Unmarshal(blob, &v); err != nil {
		log.Errorf(c, "%s unable to extract object from data, %v", desc, err)
	}
	return
}
