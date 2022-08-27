package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

/*
a new variable mu as a pointer to a sync.Mutex type. The pointer to Mutex implements
the interface sync.Locker, so you can use it as an input to the todoRouter() function.
Then use the variables mu and todoFile to run the function todoRouter(), assigning
its output to a variable t. Finally, use the variable t in the function http.StripPrefix()
to strip the /todo prefix from the URL path, passing its output to the method
*/
func newMux(todoFile string) http.Handler {

	m := http.NewServeMux()

	mu := &sync.Mutex{}

	m.HandleFunc("/", rootHandler)

	t := todoRouter(todoFile, mu)

	m.Handle("/todo", http.StripPrefix("/todo", t))
	m.Handle("/todo/", http.StripPrefix("/todo/", t))

	return m
}

func replyJSONContent(w http.ResponseWriter, r *http.Request, status int, resp *todoResponse) {

	body, err := json.Marshal(resp)

	if err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)

}

func replyError(w http.ResponseWriter, r *http.Request, status int, message string) {
	log.Printf("%s %s: Error: %d %s", r.URL, r.Method, status, message)
	http.Error(w, http.StatusText(status), status)
}
