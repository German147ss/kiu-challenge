package main

import (
	"encoding/json"
	"fmt"
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

const (
	contentTypeJSON        = "application/json"
	statusUnprocessable    = http.StatusUnprocessableEntity
	statusBadRequest       = http.StatusBadRequest
	statusMethodNotAllowed = http.StatusMethodNotAllowed
)

const (
	messageHello       = "Hello, %s!"
	messageWelcomeBack = "Hello, %s! Welcome back!"
	errorEmptyName     = "Name field is required"
)

// HelloHandler maneja las solicitudes a la ruta "/hello"
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.Header.Get("Content-Type") != contentTypeJSON {
			http.Error(w, http.StatusText(statusUnprocessable), statusUnprocessable)
			return
		}

		var body RequestBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), statusBadRequest)
			return
		}

		if body.Name == "" {
			http.Error(w, errorEmptyName, statusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		message := ""
		exists := names[body.Name]
		if exists {
			message = messageWelcomeBack
		} else {
			message = messageHello
		}
		names[body.Name] = true

		w.Header().Set("Content-Type", contentTypeJSON)
		response := HelloResponse{
			Message: fmt.Sprintf(message, body.Name),
			Exists:  exists,
		}
		json.NewEncoder(w).Encode(response)

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

		w.Header().Set("Content-Type", contentTypeJSON)
		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, http.StatusText(statusMethodNotAllowed), statusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/hello", HelloHandler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
