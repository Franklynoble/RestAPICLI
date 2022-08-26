package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/Franklynoble/todocli"
	"golang.org/x/text/message"
)

var (
	ErrNoteFound   = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return

	}
	content := "There is  An API here"
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
			deletHandler(w, r, list, id, todoFile)
		case http.MethodPatch:
               pathHandler(w,r list, id, todoFile)
		default:
			message := "Method not supported"
			replyError(w, r, http.StatusMethodNotAllowed, message)

		}
	}

}

func getAllHandler(w http.ResponseWriter, r *http.Request, list *todocli.list) {
resp := &todoResponse {
	Results: *list,
	}
	replyJSONContent(w,r, http.StatusOk, resp)	
}


func  getOneHandler(w http.ResponseWriter, r *http.Request, list *todocli.list, id int) {

  resp := &todoResponse {
	Results: (*list)[id-1: id],

  }
  replyJSONContent(w, r, http.StatusOk, resp)

}

func deleteHandler(w http.ResponsWriter, r *http.Request, list *todocli.list, id int, todtodoFile) {
	list.Delete(id)

	if err := list.Save(todoFile); err != nil {
		replyError(w, r, http.StatusInternalServerError, err.Error())
	
	 return
	}
	replyTextContent(w,r , http.StatusaNoContent, "")
}


func patchHandler(w http.ResponseWriter, r *http.Request, list *todocli.Llist, id int, todoFile string) {
	q := r.URL.Query()

	if _, ok := q["complete"]; !ok {
		message := "Missing query param 'complete'"
		replyError(w,r , http.StatusBadRequest, message)
	   return
	}
}