package main

import (
	"encoding/json"
	"log"
	"net/http"

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
	mux.HandleFunc("/", server.LoggerMiddleware(getUsers))
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
	jsondata, _ := json.Marshal(data)
	w.Write(jsondata)
}
