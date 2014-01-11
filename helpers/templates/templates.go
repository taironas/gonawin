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

package templates

import (
	"bytes"
	"encoding/json"
	"net/http"
	"html/template"

	"appengine"

	"github.com/santiaago/purple-wing/helpers/log"
	usermdl "github.com/santiaago/purple-wing/models/user"
	"github.com/santiaago/purple-wing/helpers/auth"
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
	c appengine.Context,
	dynamicTemplate []byte,
	pfuncs *template.FuncMap,
	name string) error{
	
	userdata := UserData{auth.CurrentUser(r, c),}

	var funcs template.FuncMap
 
	if pfuncs == nil {
		log.Errorf(c, " Render: pFuncs should not be nil")
	} else {
		funcs = *pfuncs
	}	
	initNavFuncMap(&funcs, r, c)

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
		log.Errorf(c, "error in execute template: %q", err)
		return err
	}
	return nil
}

// Executes and Render template with the data structure and the func map passed as argument
func RenderWithData(w http.ResponseWriter, r *http.Request, c appengine.Context, t *template.Template, data interface{}, funcs template.FuncMap, id string){
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf, t.Name(), data)
	templateBytes := buf.Bytes()
	
	if err != nil{
		log.Errorf(c, " error in parse template %v: %v", t.Name(), err)
	}
	
	err = Render(w, r, c, templateBytes, &funcs, id)
	if err != nil{
		log.Errorf(c, " error when calling Render from helpers: %v", err)
	}
}

// renders data to json and writes it to response writer
func RenderJson(w http.ResponseWriter, c appengine.Context, data interface{}) error{
	return json.NewEncoder(w).Encode(data)
}

// Set all navigation pages to false caller should define only the active one
func initNavFuncMap(pfuncs *template.FuncMap, r *http.Request, c appengine.Context) {
	
	if pfuncs != nil{
		funcs := *pfuncs
		
		if _,ok := funcs[""]; !ok {
			funcs["LoggedIn"] = func() bool { return auth.LoggedIn(r, c) }
		}
		if _,ok := funcs["Teams"]; !ok {
			funcs["Teams"] = func() bool {return false}
		}
		if _,ok := funcs["Tournaments"]; !ok {
			funcs["Tournaments"] = func() bool {return false}
		}
		if _,ok := funcs["Profile"]; !ok {
			funcs["Profile"] = func() bool {return false}
		}
		if _,ok := funcs["Admin"]; !ok {
			funcs["Admin"] = func() bool {return auth.IsAdmin(r, c)}
		}
		
	} else {
		log.Errorf(c, "error in initNavFuncMap, funcs is nil, unable to init funcs map")
	}
}
