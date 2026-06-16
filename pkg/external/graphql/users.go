package graphql

import (
	"errors"
)

const BASE_URL string = "https://graphqlplaceholder.vercel.app/graphql"

const USER_BY_ID_QUERY string = `
query usersById($ID: Int!) {
  user: userById(id: $ID){
    id
    name
    username
    email
  }
}`

type UserData struct {
	User User `json:"user"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func UserByIDQuery(id int) (user User, err error) {
	operation := GraphQLRequest{Query: USER_BY_ID_QUERY, Variables: map[string]any{"ID": id}}
	response, err := SendGqlRequest[UserData](BASE_URL, operation)
	if err != nil {
		return User{}, err
	}
	if len(response.Errors) > 0 {
		return User{}, errors.New(response.Errors[0].Message)
	}
	//netip.AddrFrom4([4]byte{127, 0, 0, 1})
	return response.Data.User, nil
}
