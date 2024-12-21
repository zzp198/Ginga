package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

func worker(ports chan int, wg *sync.WaitGroup, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("192.168.1.1:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()

		results <- p
	}
}

func main() {

	ports := make(chan int, 1000)
	results := make(chan int)
	var openports []int
	var closedports []int

	var wg sync.WaitGroup

	for i := 0; i < cap(ports); i++ {
		go worker(ports, &wg, results)
	}

	go func() {
		for i := 1; i < 65535; i++ {
			wg.Add(1)
			ports <- i
		}
	}()

	for i := 1; i < 65535; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		} else {
			closedports = append(closedports, port)
		}
	}

	close(ports)
	close(results)

	sort.Ints(openports)
	sort.Ints(closedports)

	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

}
