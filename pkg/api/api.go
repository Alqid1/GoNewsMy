package api

import (
	"GoNewsMy/pkg/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type API struct {
	r  *mux.Router
	db *db.DatBase
}

func New(db *db.DatBase) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/news", api.NewPostHandler).Methods(http.MethodPost)
	/* api.r.HandleFunc("/news", api.UpdatePostHandler).Methods(http.MethodPatch)
	api.r.HandleFunc("/news", api.DeletePostHandler).Methods(http.MethodDelete) */
	// получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.PostsHandler).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	//api.r.Handle("GET /", http.FileServer(http.Dir("cmd/webapp")))
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Обработчик NewPost
func (api *API) NewPostHandler(w http.ResponseWriter, r *http.Request) {
	var p db.Post
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.db.NewPost(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Обработчик Posts
func (api *API) PostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)
	posts, err := api.db.Posts(n)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(posts)
}

/* func (api *API) UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var p db.Post
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.ID = id
	api.db.UpdatePost(p)
	w.WriteHeader(http.StatusOK)
}

func (api *API) DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api.db.DeletePost(id)
	w.WriteHeader(http.StatusOK)
} */
