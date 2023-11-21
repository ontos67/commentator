// API приложения Commentator.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"commentator/pkg/profanity"
	storage "commentator/pkg/storage/pstg"

	"github.com/gorilla/mux"
)

type API struct {
	db *storage.DB
	r  *mux.Router
}

func New(db *storage.DB) *API {
	a := API{db: db, r: mux.NewRouter()}
	a.endpoints()
	return &a
}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.HandleFunc("/comment/save", api.saveCom).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comment/del", api.deleteCom).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comment/comListP", api.comListP).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comment/comListCont", api.comListCont).Methods(http.MethodGet, http.MethodOptions)
}

// saveCom обрабатывает запрос на сохранение нового комментария
// в браузере http://localhost:9999/comment/save?userid=64&text=заманали%20комары&pubtime=12344134&ptype=A&pid=2345
// storage.Comment
func (api *API) saveCom(w http.ResponseWriter, r *http.Request) {
	var c storage.Comment
	var err error
	form := r.URL.Query()
	userids := form.Get("userid")
	c.User_id, err = strconv.Atoi(userids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Text = form.Get("text")
	pubtimes := form.Get("pubtime")
	c.PubTime, err = strconv.ParseInt(pubtimes, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.ParentType = form.Get("ptype")
	pids := form.Get("pid")
	c.ParentID, err = strconv.Atoi(pids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l, err := profanity.ProfanityCheck([]storage.Comment{c})
	if len(l) > 0 {
		http.Error(w, fmt.Errorf("не ругайтесь. в сохранении отказано").Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.ID, err = api.db.SaveComment(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.db.CChan <- c
	json.NewEncoder(w).Encode(c.ID)
}

// deleteCom обрабатывает запрос на удаление комментария
// в браузере http://localhost:9999/comment/del?id=64
// id
func (api *API) deleteCom(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query()
	ids := form.Get("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.DeleteComment(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(id)
}

// comListP обрабатывает запрос на коммертарии на уровень ниже от родительского.
// в браузере http://localhost:9999/comment/comListP?pT=C&pId=67
// pT, pId
func (api *API) comListP(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query()
	pT := form.Get("pT")
	if pT != "A" && pT != "C" {
		http.Error(w, fmt.Errorf("wrong comment type value").Error(), http.StatusInternalServerError)
		return
	}
	pIds := form.Get("pId")
	pId, err := strconv.Atoi(pIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	l, err := api.db.CommentList(pT, pId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(l)
}

// comListCont обрабатывает запрос на коммертарии, начиная с конкретного ID и последовательно n шт.
// в браузере http://localhost:9999/comment/comListCont?sID=44&n=5
// sID, n
func (api *API) comListCont(w http.ResponseWriter, r *http.Request) {
	form := r.URL.Query()
	sIDs := form.Get("sID")
	sID, err := strconv.Atoi(sIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ns := form.Get("n")
	n, err := strconv.Atoi(ns)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	l, err := api.db.CommentsListCont(sID, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(l)
}
