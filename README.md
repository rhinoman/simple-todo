Simple TODO
===========

This is a simple TODO server side application.

## API:

**GET /todo** - List of all todo items

**POST /todo** - Create a new todo item


**GET /todo/{item-id}** - Fetch a todo item

**DELETE /todo/{item-id}** - Delete a todo item

Todo JSON example:
```
{
	"_id": "6bde7e44f6974338b920046d4e933c01",
	"_rev": "1-66753c14df1b19023762ffda70355b77",
	"type": "todo_item",
	"task_name": "Buy Milk",
	"completed": false,
	"due": "0001-01-01T00:00:00Z",
	"created": "2015-11-18T02:26:09.348165369Z"
}
```

I picked CouchDB for my database - using my couchdb driver here: https://github.com/rhinoman/couchdb-go

NOTE: The unit tests require a running couchdb instance (the docker image exposes the couchdb port)
