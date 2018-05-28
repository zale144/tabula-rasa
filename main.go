package main

import (
	"html/template"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
	"tabula-rasa/resource"
	. "tabula-rasa/db"
	"tabula-rasa/libs/memcache"
)

// setting up the router with endpoints
func main()  {
	go memcache.MemCache()
	router := httprouter.New()
	GetDBConnection()
	// REST handlers
	router.GET("/", homePage)
	router.NotFound = http.StripPrefix("/static/", http.FileServer(http.Dir("./src/tabula-rasa/web")))
	router.GET("/rest/:name/:spec", resource.TableResource{}.Get)
	router.POST("/rest/:name/:typ", resource.TableResource{}.Save)
	router.DELETE("/rest/:name/:typ", resource.TableResource{}.Delete)
	log.Fatal(http.ListenAndServe(":8080", router))
}
// serve the home page from this handler
func homePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	t, _ := template.ParseFiles("src/tabula-rasa/web/public/index.html")
	t.Execute(w, nil)
}