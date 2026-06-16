package handlers

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adapolytech/go_rest_api/pkg/external/graphql"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UserModel struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	Firstname string        `bson:"firstname" json:"firstname"`
	Lastname  string        `bson:"lastname" json:"lastname"`
	UserId    string        `bson:"userid" json:"user_id"`
	Uuid      string        `bson:"uuid" json:"uuid"`
}

type UserHandler struct {
	UserCollection *mongo.Collection
}

func NewUserHandler(database *mongo.Database) *UserHandler {
	return &UserHandler{UserCollection: database.Collection("users")}
}

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
	}
	firstname := body["firstname"].(string)
	lastname := body["lastname"].(string)
	nextId := bson.NewObjectID()
	uuid := uuid.NewString()
	user := UserModel{nextId, firstname, lastname, nextId.Hex(), uuid}
	log.Default().Printf("User to create: %+v", user)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := handler.UserCollection.InsertOne(ctx, user)
	if err != nil {
		log.Default().Printf("resp: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{"error": err.Error()})
	} else {
		log.Default().Printf("resp: %v", resp)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}

}

func (handler *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	user_info := r.Context().Value("Token1")
	fmt.Printf("%v", user_info)
	defer cancel()
	cursor, _ := handler.UserCollection.Find(ctx, bson.M{})
	defer cursor.Close(ctx)
	var users []UserModel
	var user UserModel
	for cursor.Next(ctx) {
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Printf("%v", err)
		}
		users = append(users, user)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)

}

func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user_id := r.PathValue("id")
	if user_id == "" {
		fmt.Print("UserID is required")
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	objID, err := bson.ObjectIDFromHex(user_id)
	if err != nil {
		log.Default().Printf("Not a valid object id hex %s", r.PathValue("id"))
	}
	var user UserModel
	errF := handler.UserCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if errF == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	}
}

func (handler *UserHandler) ExportUserToCSVFile(w http.ResponseWriter, r *http.Request) {
	fileName := strings.Join([]string{"file-", time.Now().String(), ".csv"}, "")
	//fileHanlder, err := os.Create(path.Join("exports", fileName))
	fileHanlder, err := os.Create("exports/" + fileName)
	if err != nil {
		log.Fatalf("Eroor occurent %s", err.Error())
	}
	defer fileHanlder.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, _ := handler.UserCollection.Find(ctx, bson.M{}, options.Find().SetLimit(5), options.Find().SetSkip(4))
	defer cursor.Close(ctx)
	// Write header
	// fmt.Fprintf(fileHanlder, "%s;%s;%s\n", "ID", "FIRSTNAME", "LASTNAME")
	var users []UserModel
	err = cursor.All(ctx, &users)
	if err != nil {
		fmt.Printf("%v", err)
		http.Error(w, "Cannot fetch users data", http.StatusInternalServerError)
	}
	csvWriter := csv.NewWriter(fileHanlder)
	data := make([][]string, len(users)+1)
	data[0] = []string{"ID", "FIRSTNAME", "LASTNAME"}
	for i, user := range users {
		data[i+1] = []string{user.UserId, user.Firstname, user.Lastname}
	}
	csvWriter.WriteAll(data)
	csvWriter.Flush()

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"filename": fileHanlder.Name()})
}

func (*UserHandler) GetExternalUserById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "UserID is required", http.StatusBadRequest)
		return
	}
	ID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := graphql.UserByIDQuery(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

type Book struct {
	title    string
	price    float64
	quantity int
}

var books []Book = make([]Book, 0)

func ReadCVSFile(path string) {

	file, err := os.OpenFile(path, os.O_RDONLY, 0666)

	if err != nil {
		log.Default().Print("No such filename")
		return
	}

	reader := csv.NewReader(file)
	for {
		str, err := reader.Read()
		if errors.Is(err, io.EOF) {
			return
		}
		book := new(Book)
		book.title = str[0]
		book.price, err = strconv.ParseFloat(str[1], 32)
		book.quantity, _ = strconv.Atoi(str[2])
		books = append(books, *book)
	}
}
