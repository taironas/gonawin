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
	err = helpers.Render(c, w, main, funcs, "renderMain")
	
	if err != nil{
		c.Errorf("pw: error when calling Render from helpers: %q", err)
	}

}
