package hello

import (
	"bytes"
	"html/template"
	"net/http"
	"appengine"
)

func init(){
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/", mainHandler)
	
}

type Content struct{
	ContainerHTML template.HTML 
}

type Data struct{
	Msg string
}

// writeMain executes the index  template.
func writeMain(w http.ResponseWriter, c Content){
	tmpl, err := template.ParseFiles("templates/index.html","templates/main.html","templates/header.html","templates/container.html","templates/footer.html","templates/scripts.html")

	if err != nil{
		print ("error in parse files")
		print (err.Error())

	}
	print("ok parse files\n")

	err = tmpl.ExecuteTemplate(w,"tmpl_index",c)
	if err != nil{
		print ("error in execute template")
		print(err.Error())
	}
	print("ok execute template\n")

}


func mainHandler(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	c.Infof("pw: mainHandler")
	c.Infof("pw: Requested URL: %v", r.URL)

	data := Data{"hello main handler"}
	print(data.Msg)
	t, err := template.ParseFiles("templates/main.html")
	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_main", data)
	main := buf.Bytes()
	
	if err != nil{
		print ("error in parse template main")
		print(err.Error())
	}
	print("ok parseTemplate main\n")
	c.Infof("pw: profile: %v", main)

	writeMain(w, Content{template.HTML(main)})

}


// writeProfile executes the index  template.
func writeProfile(w http.ResponseWriter, c Content){
	tmpl, err := template.ParseFiles("templates/index.html", "templates/container.html","templates/header.html","templates/footer.html","templates/scripts.html" )
	if err != nil{
		print ("error in parse files")
		print (err.Error())

	}
	print("ok parse files\n")

	err = tmpl.ExecuteTemplate(w,"tmpl_index",c)
	if err != nil{
		print ("error in execute template")
		print(err.Error())
	}
	print("ok execute template\n")
}

func profileHandler(w http.ResponseWriter, r *http.Request){
	c := appengine.NewContext(r)
	c.Infof("pw: profileHandler")
	c.Infof("pw: Requested URL: %v", r.URL)
	data := Data{"hello profile handler"}
	print(data.Msg)
	t, err := template.ParseFiles("templates/profile.html")
	var buf bytes.Buffer
	err = t.ExecuteTemplate(&buf,"tmpl_profile", data)
	profile := buf.Bytes()
	
	if err != nil{
		print ("error in parse template profile")
		print(err.Error())
	}
	print("ok parseTemplate profile\n")
	c.Infof("pw: profile: %v", profile)
	writeProfile(w, Content{template.HTML(profile)})
}

func parseTemplate(file string, data interface{}) (out []byte, error error) {
        var buf bytes.Buffer
        t, err := template.ParseFiles(file)
        if err != nil {
		print(err.Error())
                return nil, err
        }
        err = t.Execute(&buf, data)
        if err != nil {
		print(err.Error())
                return nil, err
        }
        return buf.Bytes(), nil
}





