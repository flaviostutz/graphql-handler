package graphqlhandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

//GraphQLHandler http handler for handling GraphQL requests
//according to https://graphql.org/learn/serving-over-http/
type GraphQLHandler struct {
	//Schema graphql schema used for handling requests
	Schema graphql.Schema
	//Debug show info and error messages on output log
	Debug bool
}

//NewGraphQLHandler creates a new handler
func NewGraphQLHandler(schema graphql.Schema, debug bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("query")
		operationName := ""
		variables := make(map[string]interface{})

		if r.Method == "POST" {
			defer r.Body.Close()

			switch r.Header.Get("Content-Type") {

			case "application/graphql":
				qb, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(500)
					if debug {
						fmt.Printf("Error reading request body. err=%s\n", err)
					}
					json.NewEncoder(w).Encode("Error processing request")
					return
				}
				query = string(qb)

			case "application/json":
				var rdata map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&rdata)
				if err != nil {
					w.WriteHeader(500)
					if debug {
						fmt.Printf("Error reading request body. err=%s\n", err)
					}
					json.NewEncoder(w).Encode("Error processing request")
					return
				}

				q, ok := rdata["query"]
				if ok {
					query = q.(string)
				}

				v, ok := rdata["variables"]
				if ok {
					variables = v.(map[string]interface{})
				}

				o, ok := rdata["operationName"]
				if ok {
					operationName = o.(string)
				}
			}

		}

		start := time.Now()
		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  query,
			VariableValues: variables,
			OperationName:  operationName,
		})

		w.Header().Add("Content-Type", "application/json")
		if result.HasErrors() {
			fmt.Printf("Error executing query '%s'. err=%v\n", query, result.Errors)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(result)
			return
		}

		if debug {
			fmt.Printf("Query '%s'; Elapsed=%s\n", query, time.Since(start))
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(result)
	}
}
