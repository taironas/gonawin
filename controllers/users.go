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

package controllers

import (
	"bytes"
	"appengine"
	"net/http"
	"html/template"
	"github.com/santiaago/purple-wing/helpers"
)

// Data struct holds the data for the template
type data struct{
	Msg string
}

// writeProfile executes the index  template.
func renderProfile(c appengine.Context, w http.ResponseWriter, content helpers.Content){
	tmpl, err := template.ParseFiles("templates/index.html", 
		"templates/container.html",
		"templates/header.html",
		"templates/footer.html",
		"templates/scripts.html" )
	if err != nil{
		print ("error in parse files")
		print (err.Error())
	}
	c.Infof("ok parse files\n")

	err = tmpl.ExecuteTemplate(w,"tmpl_index",content)
	if err != nil{
		c.Errorf("error in execute template")
		c.Errorf(err.Error())
	}
	c.Infof("ok execute template\n")
}

func ProfileHandler(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	c.Infof("pw: profileHandler")
	c.Infof("pw: Requested URL: %v", r.URL)
	
	c.Infof("pw: preparing data")
	data := data{Msg: "hello profile handler"}
	
	c.Infof("pw: data ready")
	
	c.Infof("pw: preparing profile template")
	t, err := template.ParseFiles("templates/profile.html")

	c.Infof("pw: executing profile template in standalone")
	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_profile", data)
	profile := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error in parse template profile")
		c.Errorf("pw: %v",err.Error())
	}
	c.Infof("pw: calling renderProfile()")
	renderProfile(c, w, helpers.Content{template.HTML(profile)})
	c.Infof("pw: profile handler done!")
}
