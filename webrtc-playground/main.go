package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	homeTempl = template.Must(template.ParseFiles("index.html"))
)

type IndexData struct {
	Host string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		fmt.Println(*r)
		if rec := recover(); rec != nil {
			// log.Print("panic in updating counter\n", r, string(debug.Stack()))
			// TODO: change fmt's to log.* and figure out levels
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}()

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println(err)
	}
	for k, v := range body {
		fmt.Printf("%v -> %v\n", k, v)
	}

	homeTempl.ExecuteTemplate(w, "index.html", IndexData{"localhost:8080"})
}

func main() {
	fmt.Println(*addr)
	go h.run()
	//go msgrouter.route()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ws", wsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
