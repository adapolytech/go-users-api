package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type GraphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables"`
}

type GraphQLResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []GraphQLError `json:"errors"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func SendGqlRequest[T any](URL string, operation GraphQLRequest) (GraphQLResponse[T], error) {
	var response GraphQLResponse[T]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	body, err := json.Marshal(operation)
	if err != nil {
		return response, errors.New(err.Error())
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, URL, bytes.NewBuffer(body))

	if err != nil {
		return response, errors.New(err.Error())
	}

	request.Header.Set("Content-type", "application/json")

	resp, err := client.Do(request)

	if err != nil {
		return response, errors.New(err.Error())
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("decode error: %w", err)
	}
	runtime.Goexit()
	return response, nil
}
