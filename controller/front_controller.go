package controller

import (
	"html/template"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	. "tabula-rasa/util"
	"tabula-rasa/service"
	"io/ioutil"
)

type Param struct {
	Key   string
	Value string
}

type Params []Param

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

func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	t, _ := template.ParseFiles("src/tabula-rasa/web/web/index.html")
	t.Execute(w, nil)
}

func processRest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	var out string
	name := ps.ByName("name")
	spec := ps.ByName("spec")
	id := r.URL.Query().Get("id")
	var args = []string{name, id, spec}

	switch r.Method {
	case "GET":
		out, err = service.Get(args...)
	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		CheckError(err)
		defer r.Body.Close()
		out, err = service.Create(name, b)
	case "DELETE":
		id := r.URL.Query().Get("id")
		out, err = service.Delete(name, id)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write([]byte(out))
}
