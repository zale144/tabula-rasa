package resource

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"tabula-rasa/storage"
	"encoding/json"
	"io/ioutil"
)

// table resource
type TableResource struct{}

// method for retrieving resources
func (tr TableResource) Get(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	name := ps.ByName("name")
	spec := ps.ByName("spec")
	id := r.URL.Query().Get("id")
	entity, err := storage.TableStorage{}.Get(name, id, spec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	str, err := json.Marshal(entity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(str)
}

// method for saving resources
func (tr TableResource) Save(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	typ := ps.ByName("typ")
	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "invalid parameter 'name'", http.StatusInternalServerError)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()
	out, err := storage.TableStorage{}.Save(name, typ, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

	}
	str, err := json.Marshal(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(str)
}

// method for deleting resources
func (tr TableResource) Delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params)  {
	typ := ps.ByName("typ")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "invalid parameter 'id'", http.StatusInternalServerError)
		return
	}
	name := ps.ByName("name")
	if name == "" {
		http.Error(w, "invalid parameter 'name'", http.StatusInternalServerError)
		return
	}
	err := storage.TableStorage{}.Delete(name, id, typ)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}