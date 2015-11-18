/*
 * todo_test.go
 * Unit Tests
 */

package todo_service

import (
	"testing"
)

// NOTE: These tests *REQUIRE* a running CouchDB instance
const testDbName = "todo_test_db"
const couchServer = "localhost"
const couchPort = 5984

func TestInitDb(t *testing.T) {
	defer cleanup()
	InitDb(couchServer, couchPort, testDbName)

	//Make sure the database was created
	dbs, err := connection.GetDBList()
	if err != nil {
		t.Errorf("Error getting DB List: %v", err)
	}
	foundTestDb := false
	for _, db := range dbs {
		if db == testDbName {
			foundTestDb = true
		}
	}
	if !foundTestDb {
		t.Errorf("Database not created")
	}
	//Make sure the design doc was written to the database
	ddoc := DesignDocument{}
	_, err = todoDbm.db.Read("_design/todo", &ddoc, nil)
	if err != nil {
		t.Errorf("Error reading design doc: %v", err)
	}
	if ddoc.Language != "javascript" {
		t.Errorf("Design Doc language not set!")
	}
	if len(ddoc.Views) != 1 {
		t.Errorf("Theres should be exactly ONE view!")
	}

}

func TestTodoDatabaseManager(t *testing.T) {
	defer cleanup()
	InitDb(couchServer, couchPort, testDbName)
	//Create some todos
	ti1 := TodoItem{
		TaskName: "Buy Milk",
	}
	ti2 := TodoItem{
		TaskName: "Buy Eggs",
	}

	ti1Id, err := todoDbm.CreateTodoItem(&ti1)
	if err != nil {
		t.Errorf("Error saving todo 1: %v", err)
	}
	t.Logf("Todo Item 1 ID: %v", ti1Id)
	ti2Id, err := todoDbm.CreateTodoItem(&ti2)
	if err != nil {
		t.Errorf("Error saving todo 2: %v", err)
	}
	t.Logf("Todo Item 2 ID: %v", ti2Id)

	//Now, read a todo back out
	readTi1, err := todoDbm.GetTodoItem(ti1Id)
	if err != nil {
		t.Errorf("Error fetching todo 1: %v", err)
	}
	if readTi1.TaskName != "Buy Milk" {
		t.Errorf("The Task Name is not set!")
	}

	//Now, read a list of todos
	tiList, err := todoDbm.GetTodoList()
	if err != nil {
		t.Errorf("Error fetching todo list: %v", err)
	}
	numItems := len(tiList.TodoItems)
	if numItems != 2 {
		t.Errorf("There should be 2 Todo items but there were %v", numItems)
	}

	//Now, delete some todos
	err = todoDbm.DeleteTodoItem(ti2Id)
	if err != nil {
		t.Errorf("Error Deleting todo item: %v", err)
	}
	//Fetch the list again to verify deletion
	tiList, err = todoDbm.GetTodoList()
	numItems = len(tiList.TodoItems)
	if err != nil {
		t.Errorf("Error fetching todo list: %v", err)
	}
	if len(tiList.TodoItems) != 1 {
		t.Errorf("There should be only 1 Todo item but there was %v", numItems)
	}

}

func cleanup() {
	connection.DeleteDB(testDbName, nil)
}
