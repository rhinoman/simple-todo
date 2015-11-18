/**
 * controller.go
 *
 * REST interface
 */

package todo_service

import (
	"encoding/json"
	"fmt"
	"github.com/rhinoman/couchdb-go"
	"log"
	"net/http"
)

type Controller struct{}

type ApiHandler func(resp http.ResponseWriter, req *http.Request) error

func (ah ApiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := ah(w, req); err != nil {
		log.Printf("[ERROR]: %v", err)
	}
}

// Routes

// GET /todo - List of all todo items
// POST /todo - Save a todo item
func (c Controller) HandleTodo(w http.ResponseWriter, req *http.Request) error {
	//We can have a GET or a POST here
	switch req.Method {
	case "GET":
		log.Printf("GET /todo")
		li, err := todoDbm.GetTodoList()
		if err != nil {
			cerr, ok := err.(*couchdb.Error)
			if ok {
				code := cerr.StatusCode
				w.WriteHeader(code)
			}
			return err
		}
		if err := writeBody(w, li); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		w.Header().Add("Content-Type", "application/json")
	case "POST":
		log.Printf("POST /todo")
		item := TodoItem{}
		if err := parseBody(req, &item); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return err
		}
		id, err := todoDbm.CreateTodoItem(&item)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return err
		} else {
			//Just need a simple reply struct
			type CreatedResponse struct {
				Id string `json:"new_id"`
			}
			w.WriteHeader(http.StatusCreated)
			writeBody(w, &CreatedResponse{Id: id})
		}
	default:
		//Return a 405 I suppose
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return nil
}

// GET /todo/{item-id} - Fetch a single item
// DELETE /todo/{item-id} - Delete a single item
func (c Controller) HandleTodoItem(w http.ResponseWriter, req *http.Request) error {
	//We can have a GET or DELETE here
	itemId := req.URL.Path[len("/todo/"):]
	if itemId == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	switch req.Method {
	case "GET":
		log.Printf("GET /todo/%v", itemId)
		item, err := todoDbm.GetTodoItem(itemId)
		if err != nil {
			cerr, ok := err.(*couchdb.Error)
			if ok {
				code := cerr.StatusCode
				w.WriteHeader(code)
			}
			return err
		}
		w.Header().Add("Content-Type", "application/json")
		writeBody(w, item)
	case "DELETE":
		log.Printf("DELETE /todo/%v", itemId)
		err := todoDbm.DeleteTodoItem(itemId)
		if err != nil {
			cerr, ok := err.(*couchdb.Error)
			if ok {
				code := cerr.StatusCode
				w.WriteHeader(code)
			}
			return err
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		//Method not allowed
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	return nil
}

// Root path
func (c Controller) ServeHome(w http.ResponseWriter, req *http.Request) error {
	log.Printf("Serving the home page...")
	html := `
		<!doctype html>
		<html>
		<head>
			<title>TODO</title>
		</head>
		<body>
			<h1>TODO API</h1>
			<p>Try these routes:</p>
			<ul>
				<li>GET /todo - List of Todo items</li>
				<li>POST /todo - Create a new Todo item</li>
				<li>GET /todo/{item-id} - Fetch a Todo item</li>
				<li>DELETE /todo/{item-id} - Delete a Todo item</li>
			</ul>
		</body>
		</html>`
	fmt.Fprintf(w, html)
	return nil
}

// unmarshalls a JSON Request body to a TodoItem
func parseBody(req *http.Request, item *TodoItem) error {
	return json.NewDecoder(req.Body).Decode(item)
}

// marshalls data to JSON and writes to response body
func writeBody(w http.ResponseWriter, o interface{}) error {
	if o == nil {
		return nil
	}
	if buf, err := json.Marshal(o); err != nil {
		return err
	} else {
		w.Write(buf)
	}
	return nil
}
