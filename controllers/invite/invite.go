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
 
package invite

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	
	"appengine"
	"appengine/mail"
	
	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/log"

	"github.com/santiaago/purple-wing/helpers/auth"
	templateshlp "github.com/santiaago/purple-wing/helpers/templates"
)

type InviteForm struct {
	Name string
	EmailsList string
	Error string
}

const inviteMessage = `
Hi,
Join us on Purple-wing.

You will be able to bet on tournament and win some rewards!

Sign in here: %s

Have fun,
Your friends @ purple-wing


`

func Email(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	var form InviteForm
	form.Name = auth.CurrentUser(r).Name
	
	if r.Method == "GET" {
		form.EmailsList = ""
		form.Error = ""
	} else if r.Method == "POST" {
		log.Infof(c, " Form Value = %v", r.FormValue("emails_area"))
		form.EmailsList = r.FormValue("emails_area")
		
		if len(form.EmailsList) <= 0 {
			form.Error = "No email address has been entered"
		} else {
			emails := strings.Split(form.EmailsList, ",")
			
			// validate emails
			if !helpers.AreEmailsValid(emails) {
				form.Error = "Your list of emails is not properly formatted"
			} else {
				url := fmt.Sprintf("http://%s/m/auth", r.Host) 
				for _, email := range emails {
					msg := &mail.Message{
					Sender:  "No Reply purple-wing <no-reply@purple-wing.com>",
					To:      []string{email},
					Subject: auth.CurrentUser(r).Name + " wants you to join Purple-wing!",
					Body:    fmt.Sprintf(inviteMessage, url),
					}
					
					if err := mail.Send(c, msg); err != nil {
						log.Errorf(c, " couldn't send email: %v", err)
					}
				}				
				http.Redirect(w, r, "/m/", http.StatusFound)
				return
			} 
		}
	}else {
		helpers.Error404(w)
	}
	t := template.Must(template.New("tmpl_invite_email").
		ParseFiles("templates/invite/email.html"))
	
	funcs := template.FuncMap{}
	templateshlp.RenderWithData(w, r, t, form, funcs, "renderEmail")
}
