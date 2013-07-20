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
)

// Content struct holds the parts to merge multiple templates.
// Contents are of type HTML to prevent escaping HTML.
type Content struct{
	ContainerHTML template.HTML 
}

func Render(c appengine.Context, 
	w http.ResponseWriter, 
	dynamicTemplate []byte,
	funcs template.FuncMap,
	name string) error{

	tmpl := template.Must(template.New(name).
		Funcs(funcs).
		ParseFiles("templates/layout/application.html",
		"templates/layout/header.html",
		"templates/layout/container.html",
		"templates/layout/footer.html",
		"templates/layout/scripts.html"))
	
	err := tmpl.ExecuteTemplate(w,"tmpl_application",Content{template.HTML(dynamicTemplate)})
	if err != nil{
		c.Errorf("error in execute template: %q", err)
		return err
	}
	return nil
}














