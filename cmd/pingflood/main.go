package main

import (
	"fmt"
	"github.com/calmh/jsonrpc"
	"net"
	"os"
	"sync"
	"time"
)

type Ping struct {
	ch <-chan jsonrpc.Response
	ts time.Time
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:\n  pingflood <host:port>")
		os.Exit(2)
	}

	conn, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		panic(err)
	}

	rpc := jsonrpc.NewConnection(conn, jsonrpc.ProceraDialect)
	ping := rpc.Request("system.ping")

	var lock sync.Mutex
	var resps = make(chan Ping, 10000)
	var r, minRtt, maxRtt, avgRtt int64

	go func() {
		for p := range resps {
			<-p.ch
			rtt := time.Since(p.ts).Nanoseconds() / 100000 // 10 * milliseconds

			lock.Lock()
			if minRtt == 0 || rtt < minRtt {
				minRtt = rtt
			}
			if rtt > maxRtt {
				maxRtt = rtt
			}
			avgRtt += rtt
			r++
			lock.Unlock()
		}
	}()

	var t0 = time.Now()
	var n, t int64

	for {
		ch, err := ping(nil)
		if err != nil {
			panic(err)
		}
		resps <- Ping{ch, time.Now()}
		n++
		t++
		if d := time.Since(t0).Seconds(); d >= 1 {
			lock.Lock()
			fmt.Printf("%d requests in %.01f ms; %.01f reqs/s; %d requests outstanding;", n, d*1000, float64(n)/d, t-r)
			fmt.Printf(" rtt min/avg/max %.01f/%.01f/%.01f ms\n", float64(minRtt)/10, float64(avgRtt)/10/float64(n), float64(maxRtt)/10)
			n = 0
			minRtt = 0
			maxRtt = 0
			avgRtt = 0
			lock.Unlock()
			t0 = time.Now()
		}
	}
}
