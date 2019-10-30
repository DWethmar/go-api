package main

import (
	"encoding/json"
	"html"
	"fmt"
    "log"
	"net/http"
	guuid "github.com/google/uuid"
)

// https://flaviocopes.com/golang-tutorial-rest-api/
type Node struct {
	Id 		string 	`json:"id"`
	Name 	string 	`json:"name"`
	Text 	string 	`json:"text"`
}

func main() {
    http.HandleFunc("/api/index", indexHandler)
    http.HandleFunc("/api/repo/", repoHandler)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	id := guuid.New();
	node := &Node {
		Id: id.String(),
		Name: "asd",
		Text: "lorum",
	}
	out, err := json.Marshal(node)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    fmt.Fprintf(w, string(out))
}

func repoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, %q", html.EscapeString(r.URL.Path))
}
