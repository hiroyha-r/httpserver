package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	portFlag := flag.String("port", "8080", "comma-separated list of ports to listen on (e.g. 8080,8081,8082)")
	flag.Parse()

	ports, err := parsePorts(*portFlag)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("Ok"))
	})

	var wg sync.WaitGroup
	for _, p := range ports {
		wg.Add(1)
		addr := fmt.Sprintf(":%d", p)
		go func(addr string) {
			defer wg.Done()
			log.Printf("listening on %s", addr)
			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Fatalf("server on %s: %v", addr, err)
			}
		}(addr)
	}
	wg.Wait()
}

func parsePorts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	var ports []int
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: %w", part, err)
		}
		if n < 1 || n > 65535 {
			return nil, fmt.Errorf("port out of range [1,65535]: %d", n)
		}
		ports = append(ports, n)
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports specified")
	}
	return ports, nil
}
