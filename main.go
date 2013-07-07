package hello

import (
	"bytes"
	"html/template"
	"net/http"
	"controllers"
	"models"
	"helpers"
	"appengine"
)

func init(){
	http.HandleFunc("/profile", controllers.ProfileHandler)
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/auth", controllers.Auth)
	http.HandleFunc("/oauth2callback", controllers.AuthCallback)
	http.HandleFunc("/logout", controllers.Logout)	
}

// Data struct holds the data for templates
type data struct{
	User *models.GoogleUser
	Msg string
}

// renderMain executes the main template.
// c is a Content type
// funcs are the functions needed for the main template
func renderMain(c appengine.Context, 
	w http.ResponseWriter, 
	content helpers.Content, 
	funcs template.FuncMap){

	c.Infof("pw: render main\n")

	tmpl := template.Must(template.New("renderMain").
		Funcs(funcs).
		ParseFiles("templates/index.html",
		"templates/header.html",
		"templates/container.html",
		"templates/footer.html",
		"templates/scripts.html"))

	c.Infof("template ready!\n")
	c.Infof("parse files done!\n")
	c.Infof("funcs done!\n")
	c.Infof("executing index template\n")
	err := tmpl.ExecuteTemplate(w,"tmpl_index",content)
	if err != nil{
		c.Errorf("error in execute template")
		c.Errorf(err.Error())
	}
	c.Infof("execute template done!\n")
}

//main handler: for home page
func mainHandler(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	c.Infof("pw: mainHandler")
	c.Infof("pw: Requested URL: %v", r.URL)

	c.Infof("pw: setting data")
	data := data{
		models.CurrentUser,		
		"hello main handler\n",
	}
	c.Infof("pw: data ready")
	c.Infof("pw: setting functions for template")
	funcs := template.FuncMap{"LoggedIn": helpers.LoggedIn}
	c.Infof("pw: functions ready")
	
	c.Infof("pw: preparing template main")
	t := template.Must(template.New("tmpl_main").
		Funcs(funcs).
		ParseFiles("templates/main.html"))
	c.Infof("pw: template main ready")
	
	c.Infof("pw: executing main template in standalone")
	var buf bytes.Buffer
	err := t.ExecuteTemplate(&buf,"tmpl_main", data)
	main := buf.Bytes()
	
	if err != nil{
		c.Errorf("pw: error executing template  main")
		c.Errorf("pw: %v",err.Error())
	}
	c.Infof("pw: calling renderMain()")	
	renderMain(c, w, helpers.Content{template.HTML(main)}, funcs)
	c.Infof("pw: main handler done!")
}
