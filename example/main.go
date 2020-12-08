package main

import (
	"fmt"
	"log"
	"net/http"

	ghandler "github.com/flaviostutz/graphql-handler"
	"github.com/graphql-go/graphql"
)

//Todo todo type
type Todo struct {
	Title string
}

var todos []Todo

func init() {
	todos = make([]Todo, 0)
}

func main() {
	schema, err := sampleSchema()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/graphql", ghandler.NewGraphQLHandler(schema, true))
	log.Println("Listening on :1234...")
	err = http.ListenAndServe(":1234", nil)
	log.Fatal(err)
}

var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"title": &graphql.Field{
			Type: graphql.String,
		},
	},
})

func sampleSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queries(),
		Mutation: mutations(),
	})
}

func queries() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			/*
				curl -g 'http://localhost:1234/graphql?query={todoList{title}}'
			*/
			"todoList": &graphql.Field{
				Type:        graphql.NewList(todoType),
				Description: "List of todos",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return todos, nil
				},
			},
		},
	})
}

func mutations() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			/*
				curl -g 'http://localhost:1234/graphql?query=mutation+_{createTodo(title:"My+new+todo"){title}}'
			*/
			"createTodo": &graphql.Field{
				Type:        todoType, // the return type for this field
				Description: "Create new todo",
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					title, ok := params.Args["title"]
					if !ok {
						return nil, fmt.Errorf("title is required")
					}
					nt := Todo{
						Title: title.(string),
					}
					todos = append(todos, nt)
					return nt, nil
				},
			},
		},
	})
}
