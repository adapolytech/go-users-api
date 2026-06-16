package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/adapolytech/go_rest_api/pkg/database"
	"github.com/adapolytech/go_rest_api/pkg/handlers"
	"github.com/adapolytech/go_rest_api/pkg/middlewares"
	"github.com/adapolytech/go_rest_api/pkg/server"
)

type User struct {
	Id        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var InMemoryDatabase map[string]User = map[string]User{"1": User{"1", "TheNuulest", "Developer"}}

func main() {
	_, database, err := database.CreateClient()
	if err != nil {
		log.Fatal("Cannot create database connection")
	}
	mux := server.CreateServer()
	UserService := handlers.NewUserHandler(database)
	mux.Handle("GET /users", middlewares.AuthMiddleware(UserService.GetUsers))
	mux.HandleFunc("POST /users", UserService.Register)
	mux.HandleFunc("GET /users/export", UserService.ExportUserToCSVFile)
	mux.HandleFunc("GET /users/external/{id}", UserService.GetExternalUserById)
	mux.HandleFunc("GET /users/{id}", UserService.GetUserById)
	log.Println("Server starting on http://localhost:8080")
	httpServer := http.Server{Addr: ":8080", Handler: mux, ReadTimeout: time.Second * 30}
	httpServer.ListenAndServe()
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	data := make([]User, 0, len(InMemoryDatabase))
	for _, user := range InMemoryDatabase {
		data = append(data, user)
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

func getUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, exist := InMemoryDatabase[id]
	if !exist {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)

}

func PrintFileContent(path string) {

	file, err := os.OpenFile(path, os.O_RDONLY, os.FileMode(os.O_RDONLY))

	if err != nil {
		fmt.Printf("Cannot open file %s", err.Error())
	}

	reader := bufio.NewReader(file)

	for {
		inputString, readerError := reader.ReadString('\n')
		if errors.Is(readerError, io.EOF) {
			return
		}
		fmt.Printf("The input was: %s", inputString)
	}
}
