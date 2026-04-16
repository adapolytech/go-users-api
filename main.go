package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/adapolytech/go_rest_api/packages/server"
)

type User struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var InMemoryDatabase map[string]User = map[string]User{"1": User{Id: "1", Firstname: "TheNuulest", Lastname: "Developer"}}

func main() {
	mux := server.CreateServer()
	mux.HandleFunc("GET /posts", server.LoggerMiddleware(getUsers))
	mux.HandleFunc("POST /posts", createUser)
	log.Println("Server starting on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	data := make([]User, 0, len(InMemoryDatabase))
	for key, _ := range InMemoryDatabase {
		data = append(data, InMemoryDatabase[key])
	}
	w.Header().Add("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
	}
	firstname := body["firstname"].(string)
	lastname := body["lastname"].(string)
	nextId := strconv.Itoa(len(InMemoryDatabase) + 1)
	user := User{Id: nextId, Firstname: firstname, Lastname: lastname}
	InMemoryDatabase[nextId] = user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
