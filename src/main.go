package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/florent-formation/MPS/src/lib"
	"github.com/florent-formation/MPS/src/router"
)

var road *router.Tree = router.Make()
var cacheHTTP = map[string]string{}
var cacheMutex sync.Mutex // Mutex pour synchroniser l'accès à cacheHTTP

const MAX_FIB = 92

func init() {
	// Préremplir le cache pour les valeurs de 0 à 5
	for i := 0; i <= MAX_FIB; i++ {
		fmt.Println("Cache: fib", i)
		lib.FibonacciCache(i)
	}
}

func main() {

	road.Add("/fibonacci/:number", "GET", func(w http.ResponseWriter, r *http.Request) {
		cacheMutex.Lock()
		if response, ok := cacheHTTP[r.URL.Path]; ok {
			fmt.Fprint(w, response)
			cacheMutex.Unlock()
			return
		}
		cacheMutex.Unlock()

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "Mauvais format d'URL", http.StatusBadRequest)
			return
		}

		number, err := strconv.Atoi(parts[2])
		if err != nil {
			http.Error(w, "Veuillez fournir un nombre valide", http.StatusBadRequest)
			return
		}

		if number < 0 {
			http.Error(w, "Veuillez fournir un nombre positif", http.StatusBadRequest)
			return
		}

		if number > MAX_FIB {
			http.Error(w, "La valeur est trop grande pour une exécution efficace", http.StatusRequestEntityTooLarge)
			return
		}

		result := lib.FibonacciCache(number)
		cacheMutex.Lock()
		cacheHTTP[r.URL.Path] = fmt.Sprintf("Fibonacci(%d) = %d", number, result)
		cacheMutex.Unlock()

		fmt.Fprintf(w, "Fibonacci(%d) = %d", number, result)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler, exists := road.GetHandler(r.URL.Path, r.Method)
		if !exists {
			http.Error(w, "Page non trouvée ou méthode non autorisée", http.StatusNotFound)
			return
		}
		handler(w, r)
	})

	fmt.Println("Serveur démarré sur le port 8080")
	http.ListenAndServe(":8080", nil)
}
