package hello

import (
	"html/template"
	"net/http"
	"controllers"
	"models"
	"helpers"
)

func init(){
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/auth", controllers.Auth)
	http.HandleFunc("/oauth2callback", controllers.AuthCallback)
	http.HandleFunc("/logout", controllers.Logout)
}

func mainHandler(w http.ResponseWriter, r *http.Request){
	renderHomePage(w)
}

func renderHomePage(w http.ResponseWriter) {
	funcs := template.FuncMap{"LoggedIn": helpers.LoggedIn}
	mainTemplate := template.Must(template.New("tmpl_index").Funcs(funcs).ParseFiles("templates/index.html"))
	if err := mainTemplate.Execute(w, models.CurrentUser); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}




















