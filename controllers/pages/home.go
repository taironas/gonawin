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
 
package pages

import (
	"net/http"
	"html/template"
	"bytes"
	"fmt"
	
	"appengine"
	
	"github.com/santiaago/purple-wing/helpers"
	"github.com/santiaago/purple-wing/helpers/auth"
)

//temporary main handler: for landing page
func TempHome(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Hello, purple wing!")
}

//main handler: for home page
func Home(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	data := data{
		auth.CurrentUser(r),		
		"Home handler",
	}
	
	funcs := template.FuncMap{
		"LoggedIn": func() bool { return auth.LoggedIn(r) },
		"Home": func() bool {return true},
	}
	
	t := template.Must(template.New("tmpl_main").
		Funcs(funcs).
		ParseFiles("templates/pages/main.html"))
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_main", data)
	main := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error executing template  main: %v", err)
	}
	err = helpers.Render(w, r, main, &funcs, "renderMain")
	
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers in Home Handler: %v", err)
	}

}
