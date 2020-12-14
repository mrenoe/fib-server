package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"math/big"
	"os"
)

var mu sync.Mutex
var past []int
var count = 0

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/current", current)
	http.HandleFunc("/next", next)
	http.HandleFunc("/previous", previous)
	port := os.Getenv("PORT")


	log.Fatal(http.ListenAndServe(":"+port,nil))
}

func handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "This is not a valid endpoint.\nValid endpoints are /current, /next, /previous.")
	http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
}

func current(w http.ResponseWriter, r *http.Request){
	mu.Lock()
	current := solveFib(count)
	fmt.Fprintf(w, "current -> %s\n", current)
	mu.Unlock()
}

func next(w http.ResponseWriter, r *http.Request){
	mu.Lock()
	count++
	current := solveFib(count)
	fmt.Fprintf(w, "next(%d) -> %s\n", count, current)
	mu.Unlock()
}

func previous(w http.ResponseWriter, r *http.Request){
	mu.Lock()
	previous := solveFib(count-1)
	fmt.Fprintf(w, "previous -> %s\n", previous)
	mu.Unlock()
}

func solveFib(n int) string {
	// Initialize two big ints with the first two numbers in the sequence.
	a := big.NewInt(0)
	b := big.NewInt(1)

	// Loop while a is smaller than 1e100.
	for i:=0; i < n; i++ {
		// Compute the next Fibonacci number, storing it in a.
		a.Add(a, b)
		// Swap a and b so that b is the next number in the sequence.
		a, b = b, a
	}
	fmt.Println(a) // 100-digit Fibonacci number
	return a.String()
}