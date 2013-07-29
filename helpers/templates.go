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
	"net/http"
	"html/template"

	"appengine"

	usermdl "github.com/santiaago/purple-wing/models/user"
)

// Content struct holds the parts to merge multiple templates.
// Contents are of type HTML to prevent escaping HTML.
type Content struct{
	ContainerHTML template.HTML
	UserData *UserData
}

// Data struct holds the data for templates
type UserData struct{
	User *usermdl.User
}

func Render(w http.ResponseWriter, 
	r *http.Request,
	dynamicTemplate []byte,
	funcs template.FuncMap,
	name string) error{
	
	c := appengine.NewContext(r)

	userdata := UserData{CurrentUser(r),}

	if funcs == nil {
		funcs = template.FuncMap{}
	}
	initNavFuncMap(&funcs, r)

	tmpl := template.Must(template.New(name).
		Funcs(funcs).
		ParseFiles("templates/layout/app.html",
		"templates/layout/header.html",
		"templates/layout/container.html",
		"templates/layout/nav.html",
		"templates/layout/footer.html",
		"templates/layout/scripts.html"))
	
	content := Content{
		template.HTML(dynamicTemplate),
		&userdata,
	}
	err := tmpl.ExecuteTemplate(w,"tmpl_app",content)
	if err != nil{
		c.Errorf("error in execute templateO: %q", err)
		return err
	}
	return nil
}

// set all navigation pages to false caller should define only the active one
func initNavFuncMap(pfuncs *template.FuncMap, r *http.Request) {
	
	funcs := *pfuncs
	if funcs != nil{
		if _,ok := funcs[""]; !ok {
			funcs["LoggedIn"] = func() bool { return LoggedIn(r) }
		}
		if _,ok := funcs["Home"]; !ok {
			funcs["Home"] = func() bool {return false}
		}
		if _,ok := funcs["About"]; !ok {
			funcs["About"] = func() bool {return false}
		}
		if _,ok := funcs["Contact"]; !ok {
			funcs["Contact"] = func() bool {return false}
		}
		if _,ok := funcs["Profile"]; !ok {
			funcs["Profile"] = func() bool {return false}
		}
	}else{
		c := appengine.NewContext(r)
		c.Errorf("error in initNavFuncMap, funcs is nil, unable to init funcs map")
	}
}
