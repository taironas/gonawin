package hello

import (
	"html/template"
	"net/http"
	"appengine"
)

func init(){
	http.HandleFunc("/", mainHandler)
}

var mainTemplate = template.Must(template.New("Main").Parse(mainHTML))

const mainHTML = `
<html>
<body>
<h2>hello purple-wing in Go!</h2>
<b>Contributors:</b>
<ul>
<li>Santiago (sar)</li>
<li>Carlos (cab)</li>
</body>
</html>
`


func mainHandler(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	c.Infof("pw: mainHandler")
	c.Infof("pw: Requested URL: %v", r.URL)

	if err := mainTemplate.Execute(w, nil); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}




















