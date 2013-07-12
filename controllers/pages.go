package controllers

import (
	"net/http"
	"html/template"
	"bytes"
	"appengine"
	
	"models"
	"helpers"
)

// Data struct holds the data for templates
type data struct{
	User *models.User
	Msg string
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
		ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/container.html",
		"templates/footer.html",
		"templates/scripts.html"))

	err := tmpl.ExecuteTemplate(w,"tmpl_index",content)
	if err != nil{
		c.Errorf("error in execute template: %q", err)
	}
}

//main handler: for home page
func Home(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)

	data := data{
		CurrentUser,		
		"Home handler\n",
	}
	
	funcs := template.FuncMap{"LoggedIn": LoggedIn}
	
	t := template.Must(template.New("tmpl_main").
		Funcs(funcs).
		ParseFiles("templates/main.html"))
	
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_main", data)
	main := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error executing template  main: %q", err)
	}

	renderMain(c, w, helpers.Content{template.HTML(main)}, funcs)
}