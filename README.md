# graphql-handler

HTTP Handler for creating GraphQL servers in Golang

This implementation follow the guidelines for adapting http requests to GraphQL requests from https://graphql.org/learn/serving-over-http/

## Usage

* Create your Golang code for handling web requests

```go
package main

import ghandler "github.com/flaviostutz/graphql-handler"

func main() {
    http.HandleFunc("/graphql", ghandler.NewGraphQLHandler(sampleSchema(), true))
    log.Println("Listening on :1234...")
    http.ListenAndServe(":1234", nil)
}
```

* Run `go run example.go`

* Create a todo item with `curl -g 'http://localhost:1234/graphql?query=mutation+_{createTodo(title:"My+new+todo"){title}}'`

* List all todo items with `curl -g 'http://localhost:1234/graphql?query={todoList{title}}'`

