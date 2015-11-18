/**
 * main.go
 *
 * Starts the server
 */

package main

import (
	"flag"
	"github.com/rhinoman/simple-todo/todo_service"
	"log"
	"net/http"
	"time"
)

func main() {
	// Get command line arguments
	var dbHost = flag.String("dbHost", "localhost", "location(URL) of CouchDB server")
	var dbPort = flag.Int("dbPort", 5984, "Port of CouchDB server")
	var port = flag.String("p", "8085", "port for HTTP server")
	flag.Parse()
	log.Printf("CouchDB host: %v on port %v", dbHost, dbPort)
	//Initialize the database
	todo_service.InitDb(*dbHost, *dbPort, "todo_db")
	controller := todo_service.Controller{}
	mux := http.NewServeMux()
	mux.Handle("/", todo_service.ApiHandler(controller.ServeHome))
	mux.Handle("/todo/", todo_service.ApiHandler(controller.HandleTodoItem))
	mux.Handle("/todo", todo_service.ApiHandler(controller.HandleTodo))
	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Printf("Listening for HTTP clients on port %v", *port)
	//Start serving
	log.Fatal(server.ListenAndServe())
}
