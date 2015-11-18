/*
 * database.go
 * Database Access functions
 */

package todo_service

import (
	"github.com/rhinoman/couchdb-go"
	"github.com/twinj/uuid"
	"log"
	"net/url"
	"strings"
	"time"
)

type TodoItem struct {
	Id        string    `json:"_id,omitempty"`
	Rev       string    `json:"_rev,omitempty"`
	Type      string    `json:"type"`
	TaskName  string    `json:"task_name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due,omitempty"`
	Created   time.Time `json:"created"`
}

func (ti TodoItem) Validate() bool {
	if ti.Type != "todo_item" {
		return false
	}
	if ti.TaskName == "" {
		return false
	}
	return true
}

type TodoList struct {
	TodoItems []TodoItem `json:"items"`
}

//CouchDB design doc for queries
type DesignDocument struct {
	Language string          `json:"language"`
	Views    map[string]View `json:"views"`
}

type View struct {
	Map    string `json:"map"`
	Reduce string `json:"reduce,omitempty"`
}

//Responses to couchdb view
type ListResponse struct {
	TotalRows int `json:"total_rows"`
	Offset    int `json:"offset"`
	Rows      []struct {
		Id    string   `json:"id"`
		Key   string   `json:"key"`
		Value TodoItem `json:"value"`
	}
}

//CouchDB view to query for all 'todo' items
var todoViews = map[string]View{
	"getAllTodos": {
		Map: `
			function(doc){ 
				if(doc.type==="todo_item"){ 
					emit(doc.created, doc); 
				} 
			}`,
		Reduce: "_count",
	},
}

// Holds the connection to the database
var connection *couchdb.Connection

// Our database 'object'
type DatabaseManager struct {
	db *couchdb.Database
}

var todoDbm *DatabaseManager

// Name of the couch database
var dbName string

// Sets up the database connection
func InitDb(host string, port int, db string) {
	var err error
	dbName = db
	log.Println("Initializing Database Connection")
	timeoutMs := time.Duration(100) * time.Millisecond
	connection, err = couchdb.NewConnection(host, port, timeoutMs)
	if err != nil {
		log.Fatalf("Could not connect to CouchDB: %v", err)
	}
	setupDb()
	writeDDoc()
}

// Checks if database exists, creating it if necessary
func setupDb() {
	//See if our database exists
	dbList, err := connection.GetDBList()
	if err != nil {
		log.Fatalf("Could not query CouchDB for database list: %v", err)
	}
	dbFound := false
	for _, db := range dbList {
		if db == dbName {
			dbFound = true
		}
	}
	//Database does not exists, looks like we have to create it
	if !dbFound {
		if err = connection.CreateDB(dbName, nil); err != nil {
			log.Fatalf("Could not create database: %v", err)
		}
	}
	todoDbm = &DatabaseManager{
		db: connection.SelectDB(dbName, nil),
	}
}

// Checks if design document is in the database, creating it if necessary
func writeDDoc() {
	var o interface{} //we don't care about the returned data
	if _, err := todoDbm.db.Read("_design/todo", o, nil); err != nil {
		if strings.Contains(err.Error(), "404") {
			//Document does not exist, create it
			ddoc := DesignDocument{
				Language: "javascript",
				Views:    todoViews,
			}
			_, err = todoDbm.db.SaveDesignDoc("todo", &ddoc, "")
			if err != nil {
				log.Fatalf("Couldn't save design document: %v", err)
			}
		}
	}
}

// Reads a single Todo Item
// Returns an error on failure
func (dbm *DatabaseManager) GetTodoItem(id string) (*TodoItem, error) {
	ti := TodoItem{}
	if _, err := dbm.db.Read(id, &ti, nil); err != nil {
		return nil, err
	} else {
		return &ti, nil
	}
}

// Create a Todo item
// Returns the id of the new item, or an error on failure
func (dbm *DatabaseManager) CreateTodoItem(item *TodoItem) (string, error) {
	//Generate a UUID for this new Todo item
	id := uuid.Formatter(uuid.NewV4(), uuid.Clean)
	//Fill in some values
	nowTime := time.Now().UTC()
	item.Created = nowTime
	item.Type = "todo_item"
	//validate
	if !item.Validate() {
		return "", &couchdb.Error{
			StatusCode: 400,
			Reason:     "TodoItem is invalid",
		}
	}
	//Save it to the database
	if _, err := dbm.db.Save(item, id, ""); err != nil {
		return "", err
	} else {
		return id, nil
	}
}

// Deletes a Todo item
// Returns an error on failure
func (dbm *DatabaseManager) DeleteTodoItem(id string) error {
	//We're going to need the rev, so read it first
	ti := TodoItem{}
	if rev, err := dbm.db.Read(id, &ti, nil); err != nil {
		return err
	} else {
		_, err = dbm.db.Delete(id, rev)
		return err
	}
}

// Gets all the Todo Items
// Returns an error on failure
func (dbm *DatabaseManager) GetTodoList() (*TodoList, error) {
	til := TodoList{}
	lr := ListResponse{}
	params := url.Values{}
	params.Add("reduce", "false")
	if err := dbm.db.GetView("todo", "getAllTodos", &lr, &params); err != nil {
		return nil, err
	} else {
		//copy the Todo items out of the couch list response
		for _, row := range lr.Rows {
			til.TodoItems = append(til.TodoItems, row.Value)
		}
		return &til, nil
	}
}
