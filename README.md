# graphql-handler

HTTP Handler for creating GraphQL servers in Golang using https://github.com/graphql-go/graphql.

This implementation follows the guidelines for adapting http requests to GraphQL requests from https://graphql.org/learn/serving-over-http/

For a complete example, see [/example](/example)

Attention: after creating this lib I found a more "official" one that does almost the same thing: https://github.com/graphql-go/handler

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

* Function "sampleSchema()" is where you build the GraphQL schema using https://github.com/graphql-go/graphql

* Run `go run example.go`

* Create a todo item with `curl -g 'http://localhost:1234/graphql?query=mutation+_{createTodo(title:"My+new+todo"){title}}'`

* List all todo items with `curl -g 'http://localhost:1234/graphql?query={todos{title}}'`

