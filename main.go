package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"math/big"
	"sync"
)

var (
	addr     = flag.String("addr", ":8080", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Whether to enable transparent response compression")
)

var mu sync.Mutex
var past = make(map[uint]*big.Int, 0)
var count uint = 0

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
	switch string(ctx.Path()) {
	case "/current":
		current(ctx)
	case "/next":
		next(ctx)
	case "/previous":
		previous(ctx)
	case "/reset":
		reset(ctx)
	default:
		ctx.Error("Unsupported path\nSupported Paths are /current, /next, and /previous", fasthttp.StatusUnprocessableEntity)
	}

}

func current(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	mu.Lock()
	current := solveFib(count)
	fmt.Fprintf(ctx, "%s\n", current)
	mu.Unlock()
}

func reset(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	mu.Lock()
	count = 0
	fmt.Fprintf(ctx, "Reset back to 0\n")
	mu.Unlock()
}

func next(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	mu.Lock()
	count++
	current := solveFib(count)
	fmt.Fprintf(ctx, "fib(%d) -> %s\n", count, current)
	mu.Unlock()
}

func previous(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")
	mu.Lock()
	previous := solveFib(count - 1)
	fmt.Fprintf(ctx, "%s\n", previous)
	mu.Unlock()
}

func solveFib(n uint) string {
	if val, ok := past[n]; ok {
		//log.Println("Using cache!")
		log.Println(val) //  Fibonacci number
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
			//log.Println("Using additive cache!")
			past[n] = a
			return a.String()
		}
		a := big.NewInt(0)
		b := big.NewInt(1)
		var i uint
		for i = 0; i < n; i++ {
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
