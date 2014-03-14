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

// Package invite provides the JSON handlers to send invitations to gonawin app.
package invite

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"appengine"
	"appengine/mail"

	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
)

const inviteMessage = `
Hi,
Join us on Gonawin.

You will be able to bet on tournament and win some rewards!

Sign in here: %s

Have fun,
Your friends @ Gonawin


`

// invite json handler
func InviteJson(w http.ResponseWriter, r *http.Request) error {
	c := appengine.NewContext(r)

	if r.Method == "POST" {
		emailsList := r.FormValue("emails")
		name := r.FormValue("name")

		if len(emailsList) <= 0 {
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInviteNoEmailAddr)}
		}
		emails := strings.Split(emailsList, ",")
		// validate emails
		if !helpers.AreEmailsValid(emails) {
			return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInviteEmailsInvalid)}
		}

		url := fmt.Sprintf("http://%s/ng#", r.Host)
		for _, email := range emails {
			msg := &mail.Message{
				Sender:  "No Reply gonawin <no-reply@gonawin.com>",
				To:      []string{email},
				Subject: name + " wants you to join Gonawin!",
				Body:    fmt.Sprintf(inviteMessage, url),
			}

			if err := mail.Send(c, msg); err != nil {
				log.Errorf(c, "Invite Handler: couldn't send email: %v", err)
				return &helpers.InternalServerError{Err: errors.New(helpers.ErrorCodeInviteEmailCannotSend)}
			}
		}
		return templateshlp.RenderJson(w, c, "Email has been sent successfully")
	}
	return &helpers.BadRequest{Err: errors.New(helpers.ErrorCodeNotSupported)}
}
