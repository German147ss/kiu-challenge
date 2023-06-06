package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type RequestBody struct {
	Name string `json:"name"`
}

type HelloResponse struct {
	Message string `json:"message"`
	Exists  bool   `json:"exists"`
}

type NamesResponse struct {
	Names []string `json:"names"`
}

var (
	names = make(map[string]bool)
	mu    = &sync.Mutex{}
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}

		var body RequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		message := "Hello, " + body.Name + "!"
		exists := names[body.Name]
		if exists {
			message = "Hello, " + body.Name + "! Welcome back!"
		}
		names[body.Name] = true

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(HelloResponse{Message: message, Exists: exists})

	case http.MethodGet:
		mu.Lock()
		defer mu.Unlock()

		allNames := make([]string, 0, len(names))
		for name := range names {
			allNames = append(allNames, name)
		}

		var resp NamesResponse
		if len(allNames) > 0 {
			resp.Names = allNames
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

	default:
		if r.Method != http.MethodPut && r.Method != http.MethodDelete {
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/hello", HelloHandler)
	http.ListenAndServe(":8080", nil)
}
