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

// Package invite provides the JSON handlers to send invitations to gonawin app.
package invite

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"appengine"
	"appengine/taskqueue"

	"github.com/taironas/gonawin/helpers"
	"github.com/taironas/gonawin/helpers/log"
	templateshlp "github.com/taironas/gonawin/helpers/templates"

	mdl "github.com/taironas/gonawin/models"
)

// Invite handler, use it to invite users to use gonawin.
// It expects a list of emails using the url param 'emails'
//
func Invite(w http.ResponseWriter, r *http.Request, u *mdl.User) error {
	if r.Method != "POST" {
		return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
	}

	desc := "invite handler:"
	c := appengine.NewContext(r)

	var emailsList string
	if emailsList = r.FormValue("emails"); len(emailsList) <= 0 {
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInviteNoEmailAddr)}
	}

	emails := parseEmails(emailsList)

	if !helpers.AreEmailsValid(emails) {
		log.Errorf(c, "%s emails are not valid", desc, emails)
		return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInviteEmailsInvalid)}
	}

	body := buildEmailBody(r)

	// send tasks to send emails
	processEmails(c, desc, emails, body, u)

	vm := buildInviteViewModel()

	return templateshlp.RenderJson(w, c, vm)
}

func processEmails(c appengine.Context, desc string, emails []string, body string, u *mdl.User) error {

	bname := encode(c, desc, u.Name)
	bbody := encode(c, desc, body)

	for _, email := range emails {
		bemail := encode(c, desc, email)

		task := taskqueue.NewPOSTTask("/a/invite/", url.Values{
			"email": []string{string(bemail)},
			"name":  []string{string(bname)},
			"body":  []string{string(bbody)},
		})
		if _, err := taskqueue.Add(c, task, ""); err != nil {
			log.Errorf(c, "%s unable to add task to taskqueue %v", desc, err)
			return err
		}
	}
	return nil
}

type inviteViewModel struct {
	MessageInfo string `json:",omitempty"`
}

func buildInviteViewModel() inviteViewModel {
	msg := fmt.Sprintf("Email invitations have been successfully sent.")
	return inviteViewModel{msg}
}

// parseEmails split string by ',' and
// removes the leading and trailing spaces from each email.
//
func parseEmails(emailsList string) []string {
	rawEmails := strings.Split(emailsList, ",")
	emails := make([]string, 0)
	for _, e := range rawEmails {
		emails = append(emails, strings.Trim(e, " "))
	}
	return emails
}

func encode(c appengine.Context, desc, s string) []byte {
	var b []byte
	var err error
	if b, err = json.Marshal(s); err != nil {
		log.Errorf(c, "%s Error marshaling %v", desc, err)
	}
	return b
}

const inviteMessage = `
Hi there,
Join us at gonawin.

You will be able to join your friends and compete with them by predicting the results of your favorite sports events!

Sign in here: %s

Have fun,
Your friends @ Gonawin


`

func buildEmailBody(r *http.Request) string {
	currenturl := fmt.Sprintf("http://%s/#", r.Host)
	return fmt.Sprintf(inviteMessage, currenturl)
}
