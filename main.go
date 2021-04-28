package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Teacher struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Course struct {
	Id        int      `json:"id"`
	User_id   string   `json:"user_id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Img       string   `json:"img"`
	Desc      string   `json:"desc"`
	Date      string   `json:"date"`
	Timestamp string   `json:"timestamp"`
	Teacher   Teacher  `json:"teacher"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func courses() []Course {
	t := Teacher{"Max", "https://images.unsplash.com/photo-1558531304-a4773b7e3a9c?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=634&q=80"}
	cTags := []string{"colaborate", "git", "cli", "commit", "versionning"}
	c := Course{1, "eba25511-afce-4c8e-8cab-f82822434648", "learn git", cTags, "https://carlchenet.com/wp-content/uploads/2019/04/git-logo.png", "Learn how to create, manage, fork, and collaborate on a project. Git stays a major part of all companies projects. Learning git is learning how to make your project better everyday", "5 nov", "1604577600000", t}
	return []Course{c}
}

func main() {
	r := mux.NewRouter()
	r.Handle("/api/v1/courses", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonResponse, err := json.Marshal(courses())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)

	}))

	http.ListenAndServe(":8000", r)
}
