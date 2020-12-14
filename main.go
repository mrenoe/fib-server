package main

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
)

var mu sync.Mutex
var past = make(map[int]*big.Int, 0)
var count = 0

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/current", current)
	http.HandleFunc("/next", next)
	http.HandleFunc("/previous", previous)
	port, ok := os.LookupEnv("PORT")
	if ok == false || port == "" {
		port = "8080"
	}

	log.Println("listen on localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "This is not a valid endpoint.\nValid endpoints are /current, /next, /previous.")
}

func current(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	current := solveFib(count)
	fmt.Fprintf(w, "current -> %s\n", current)
	mu.Unlock()
}

func reset(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count = 0
	fmt.Fprintf(w, "Reset back to 0\n")
	mu.Unlock()
}

func next(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	count++
	current := solveFib(count)
	fmt.Fprintf(w, "next(%d) -> %s\n", count, current)
	mu.Unlock()
}

func previous(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	previous := solveFib(count - 1)
	fmt.Fprintf(w, "previous -> %s\n", previous)
	mu.Unlock()
}

func solveFib(n int) string {
	if val, ok := past[n]; ok {
		log.Println("Using cache!")
		log.Println(val) // 100-digit Fibonacci number
		return val.String()
	}
	switch n {
	case 0:
		past[n] = big.NewInt(0)
		return "0"
	case 1:
		past[n] = big.NewInt(1)
		return "1"
	default:
		n1 := past[n-1]
		n2 := past[n-2]
		if n1 != nil && n2 != nil {
			a := new(big.Int)

			a.Add(n1, n2)
			log.Println("Using additive cache!")
			past[n] = a
			return a.String()
		}
		a := big.NewInt(0)
		b := big.NewInt(1)

		for i := 0; i < n; i++ {
			// Compute the next Fibonacci number, storing it in a.
			a.Add(a, b)
			// Swap a and b so that b is the next number in the sequence.
			a, b = b, a
		}
		log.Println(a) // 100-digit Fibonacci number
		past[n] = a
		return a.String()
	}
}
