## API TESTING ROUTES

### Send GET request for users list
```shell
curl -sS -X GET http://localhost:8080/posts | jq .
```

### POST request to create user with firsname and lastname 

```shell
curl -sS -X POST -d '{"firstname": "John", "lastname":"Doe"}' http://localhost:8080/posts | jq .
```
