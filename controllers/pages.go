package controllers

import (
	"net/http"
	"html/template"
	"bytes"
	"appengine"
	
	"github.com/santiaago/purple-wing/models"
	"github.com/santiaago/purple-wing/helpers"
)

// Data struct holds the data for templates
type data struct{
	User *models.User
	Msg string
}

//main handler: for home page
func Home(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	data := data{
		CurrentUser,		
		"Home handler",
	}
	
	funcs := template.FuncMap{"LoggedIn": LoggedIn}
	
	t := template.Must(template.New("tmpl_main").
		Funcs(funcs).
		ParseFiles("templates/pages/main.html"))
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_main", data)
	main := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error executing template  main: %q", err)
	}

	renderMain(c, w, helpers.Content{template.HTML(main)}, funcs)
}

// renderMain executes the main template.
// c is a Content type
// funcs are the functions needed for the main template
func renderMain(c appengine.Context, 
	w http.ResponseWriter, 
	content helpers.Content, 
	funcs template.FuncMap){

	tmpl := template.Must(template.New("renderMain").
		Funcs(funcs).
		ParseFiles(	"templates/layout/application.html",
					"templates/layout/header.html",
					"templates/layout/container.html",
					"templates/layout/footer.html",
					"templates/layout/scripts.html"))

	err := tmpl.ExecuteTemplate(w,"tmpl_application",content)
	if err != nil{
		c.Errorf("error in execute template: %q", err)
	}
}