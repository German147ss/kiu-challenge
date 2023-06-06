package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	// Test invalid content type
	req, _ := http.NewRequest(http.MethodPost, "/hello", nil)
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()
	HelloHandler(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
	//Pass msg
	fmt.Println("Test: invalid content type - passed")

	// Test POST method
	postBody, _ := json.Marshal(map[string]string{
		"name": "test",
	})
	req, _ = http.NewRequest(http.MethodPost, "/hello", bytes.NewBuffer(postBody))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	HelloHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	//Pass msg
	fmt.Println("Test: POST method - passed")

	// Test GET method
	req, _ = http.NewRequest(http.MethodGet, "/hello", nil)
	rr = httptest.NewRecorder()
	HelloHandler(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	//Pass msg
	fmt.Println("Test: GET method - passed")

	// Test unsupported method
	req, _ = http.NewRequest(http.MethodPut, "/hello", nil)
	rr = httptest.NewRecorder()
	HelloHandler(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
	//Pass msg
	fmt.Println("Test: unsupported method - passed")
}

func TestHelloHandlerConcurrency(t *testing.T) {
	const concurrentRequests = 1000 // Aquí se define la cantidad de solicitudes simultáneas que se realizarán.
	const name = "ConcurrentUser"   // Esto es lo que se usará como "name" en las solicitudes POST.

	var wg sync.WaitGroup
	wg.Add(concurrentRequests) // Se está utilizando un WaitGroup para asegurar que todas las goroutines finalicen antes de que el test termine.

	// Aquí se limpia el mapa "names" antes del test.
	names = make(map[string]bool)

	for i := 0; i < concurrentRequests; i++ {
		go func() { // Cada goroutine ejecutará este código.
			defer wg.Done() // Una vez que esta goroutine termine su trabajo, le informará al WaitGroup que ha terminado.

			// Aquí se crea una nueva solicitud POST con "name" en el cuerpo.
			postBody, _ := json.Marshal(map[string]string{
				"name": name,
			})
			req, _ := http.NewRequest(http.MethodPost, "/hello", bytes.NewBuffer(postBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			HelloHandler(rr, req)

			var resp HelloResponse
			if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
				t.Errorf("Error decoding response body: %v", err)
				return
			}

			// Como todas las solicitudes están utilizando el mismo nombre y se están realizando de manera concurrente,
			// sólo la primera debería obtener un "exists" como false y todas las demás deberían obtenerlo como true.
			if resp.Exists && resp.Message != "Hello, "+name+"! Welcome back!" {
				t.Errorf("Expected message is 'Hello, %s! Welcome back!', but got '%s'", name, resp.Message)
			} else if !resp.Exists && resp.Message != "Hello, "+name+"!" {
				t.Errorf("Expected message is 'Hello, %s!', but got '%s'", name, resp.Message)
			}
		}()
	}

	wg.Wait() // Esperamos..

	// Al final del test, se verifica si el nombre se almacenó correctamente.
	mu.Lock()
	defer mu.Unlock()
	if !names[name] {
		t.Errorf("Expected '%s' to be in names map", name)
	}
}
