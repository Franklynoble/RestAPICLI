package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/Franklynoble/todocli"
)

var (
	ErrNoteFound   = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		replyError(w, r, http.StatusNotFound, "")
		return

	}
	content := "There's  An API here"
	replyTextContent(w, r, http.StatusOK, content)

}

func replyTextContent(w http.ResponseWriter, r *http.Request,
	status int, content string) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content))

}

// For this example, you’re using the method l.Lock() from the sync.Locker interface
// to lock the entire request handling
// This prevents concurrent access to the
//file represented by the variable todoFile which could lead to data loss. This is
func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		list := &todocli.List{}

		l.Lock()
		defer l.Unlock()
		if err := list.Get(todoFile); err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		if r.URL.Path == "" {
			switch r.Method {
			case http.MethodGet:
				getAllHandler(w, r, list)
			case http.MethodPost:
				addHandler(w, r, list, todoFile)

			default:
				message := "Method not supported"
				replyError(w, r, http.StatusMethodNotAllowed, message)
			}
			//The return statement at the end of switch ensures that we finish processing any
			// requests to the /todo root.
			return
		}

		id, err := validateID(r.URL.Path, list)

		if err != nil {
			if errors.Is(err, ErrNoteFound) {
				replyError(w, r, http.StatusNotFound, err.Error())
				return
			}
			replyError(w, r, http.StatusNotFound, err.Error())
		}

		switch r.Method {
		case http.MethodGet:
			getOneHandler(w, r, list, id)

		case http.MethodDelete:
			deleteHandler(w, r, list, id, todoFile)
		case http.MethodPatch:
			patchHandler(w, r, list, id, todoFile)
		default:
			message := "Method not supported"
			replyError(w, r, http.StatusMethodNotAllowed, message)

		}
	}

}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todocli.List) {
	resp := &todoResponse{
		Results: *list,
	}
	replyJSONContent(w, r, http.StatusOK, resp)
}

func getOneHandler(w http.ResponseWriter, r *http.Request, list *todocli.List, id int) {

	resp := &todoResponse{
		Results: (*list)[id-1 : id],
	}
	replyJSONContent(w, r, http.StatusOK, resp)

}

func deleteHandler(w http.ResponseWriter, r *http.Request, list *todocli.List, id int, todoFile string) {
	list.Delete(id)

	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())

		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}

func patchHandler(w http.ResponseWriter, r *http.Request, list *todocli.List, id int, todoFile string) {
	q := r.URL.Query()

	if _, ok := q["complete"]; !ok {
		message := "Missing query param 'complete'"
		replyError(w, r, http.StatusBadRequest, message)
		return
	}
	list.Complete(id)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusNoContent, "")
}

func addHandler(w http.ResponseWriter, r *http.Request, list *todocli.List, todoFile string) {
	item := struct {
		Task string `json:"task"`
	}{}
	//  representing the task name to include in the list. It decodes the JSON into
	// an anonymous struct and uses it to add the new item. It replies with the HTTP
	// Status Created if successful or an appropriate error status if errors occur
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		message := fmt.Sprintf("Invalid JSON: %s", err)
		replyError(w, r, http.StatusBadRequest, message)
		return
	}

	list.Add(item.Task)
	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	replyTextContent(w, r, http.StatusCreated, "")
}

func validateID(path string, list *todocli.List) (int, error) {
	//convert string id to int using strconv.Atoi
	id, err := strconv.Atoi(path)

	if err != nil {
		return 0, fmt.Errorf("%w: Invalid ID: %s", ErrInvalidData, err)
	}
	if id < 1 {
		return 0, fmt.Errorf("%W, invalid ID: Less than one", ErrInvalidData)
	}
	if id > len(*list) {
		return id, fmt.Errorf("%w: ID %d not found", ErrNoteFound, id)
	}

	return id, nil
}
