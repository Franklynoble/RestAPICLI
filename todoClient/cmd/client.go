package cmd

import (
	"errors"
	"net/http"
	"time"
)

var (
	ErrConnection        = errors.New("Connection error")
	ErroNotFOund         = errors.New("Not found")
	ErrorInvalidResponse = errors.New("Invalid Server Response")
	ErroInvalid          = errors.New("Invalid data")
	ErrorNotNumber       = errors.New("Not a number")
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

//What  the API response should be
type response struct {
	Results      []item `json:"results"`
	Date         int    `json:"date"`
	TotalResults int    `json:"total_results"`
}

/*
the connection timeout. The default client has no timeout set, which means
your application could take a long time to return or hang forever if the server
has an issue. Let’s define a function to instantiate a new client with a timeout
of 10 seconds. If you want, you could make this value customizable. For now,
we’ll keep it hardcoded

*/
func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	return c
}
