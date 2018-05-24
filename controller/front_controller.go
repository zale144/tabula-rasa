package controller

import (
	"html/template"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	. "tabula-rasa/util"
	"io/ioutil"
	"tabula-rasa/dao"
)

type Param struct {
	Key   string
	Value string
}
// TODO refactor everything
type Params []Param
// setting up the router with endpoints
func SetupRouter()  {
	router := httprouter.New()
	GetDBConnection()
	router.GET("/", homePage)
	router.NotFound = http.StripPrefix("/static/", http.FileServer(http.Dir("./src/tabula-rasa/web")))
	router.GET("/rest/:name", processRest)
	router.GET("/rest/:name/:spec", processRest)
	router.POST("/rest/:name", processRest)
	router.DELETE("/rest/:name", processRest)
	log.Fatal(http.ListenAndServe(":8080", router))
}
// serve the home page from this handler
func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	t, _ := template.ParseFiles("src/tabula-rasa/web/public/index.html")
	t.Execute(w, nil)
}
// generic processing of RESTful http requests
func processRest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	var out string
	name := ps.ByName("name")
	spec := ps.ByName("spec")
	id := r.URL.Query().Get("id")
	var args = []string{name, id, spec}

	switch r.Method {
	case "GET":
		out, err = dao.Get(args...)
	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		CheckError(err)
		defer r.Body.Close()
		out, err = dao.Create(name, b)
	case "DELETE":
		id := r.URL.Query().Get("id")
		out, err = dao.Delete(name, id)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}
