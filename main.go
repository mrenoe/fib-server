package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math"
	"math/big"
	"sync"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

//Because we're using a global map, we need to mutex lock each write
var mu sync.Mutex

//This is the cache of past fibonacci instances
var past = make(map[uint64]*big.Int, 0)

//count is the global counter of where the server is at in calculating fibonacci numbers
var count uint64 = 0

func main() {
	flag.Parse()
	h := requestHandler
	if *compress {
		h = fasthttp.CompressHandler(h)
	}

	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	//Routes and their handlers
	switch string(ctx.Path()) {
	case "/current":
		current(ctx)
	case "/next":
		next(ctx)
	case "/previous":
		previous(ctx)
		//Added reset for development. Will exist behind the scenes for testing purposes
	case "/reset":
		reset(ctx)
	default:
		ctx.Error("Unsupported path\nSupported Paths are /current, /next, and /previous", fasthttp.StatusUnprocessableEntity)
	}

}

func current(ctx *fasthttp.RequestCtx) {
	mu.Lock()
	current := solveFib(count)
	fmt.Fprintf(ctx, "%s\n", current)
	mu.Unlock()
}

func reset(ctx *fasthttp.RequestCtx) {
	mu.Lock()
	count = 0
	fmt.Fprintf(ctx, "Reset back to 0\n")
	mu.Unlock()
}

func next(ctx *fasthttp.RequestCtx) {
	mu.Lock()
	count++
	current := solveFib(count)
	if count == math.MaxUint64 {
		count = 0
	}
	fmt.Fprintf(ctx, "fib(%d) -> %s\n", count, current)
	mu.Unlock()
}

func previous(ctx *fasthttp.RequestCtx) {
	mu.Lock()
	previous := solveFib(count - 1)
	fmt.Fprintf(ctx, "%s\n", previous)
	mu.Unlock()
}

func solveFib(n uint64) string {
	if val, ok := past[n]; ok {
		//If there's a previous cache hit here, return that
		//log.Println(val)
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
		//If there's a past instance of n-1 or n-2, use those rather than adding from 0 each time
		n1 := past[n-1]
		n2 := past[n-2]
		if n1 != nil && n2 != nil {
			a := new(big.Int)

			a.Add(n1, n2)
			//log.Println("Using additive cache!")
			past[n] = a
			return a.String()
		}
		//However if there's no fib(n-1) and fib(n-2), we're going to solve it and store it
		a := big.NewInt(0)
		b := big.NewInt(1)
		var i uint64
		for i = 0; i < n; i++ {
			// Compute the next Fibonacci number, storing it in a.
			a.Add(a, b)
			// Swap a and b so that b is the next number in the sequence.
			a, b = b, a
		}
		//And place it in the cache
		past[n] = a
		return a.String()
	}
}
