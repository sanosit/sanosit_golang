package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)

	wrkrs := strings.Split(os.Getenv("WORKERS"), ":")
	if wrkrs[0] == "" {
		log.Fatalln("no workers to connect to")
	}

	c := make(chan int)
	go counter(c, len(wrkrs))

	http.HandleFunc("/", proxyhandler(wrkrs, c))
	log.Fatalln(http.ListenAndServe(port, nil))
}

func counter(c chan int, n int) {
	for {
		for i := 0; i < n; i++ {
			c <- i
		}
	}
}

func proxyhandler(wrkrs []string, c chan int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		u := r.URL
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading request: %s", err)
			return
		}
		r.Body.Close()

		n := <-c
		client := &http.Client{}
		addr := fmt.Sprintf("http://localhost:%s%s", wrkrs[n], u.Path)
		req, err := http.NewRequest(method, addr, bytes.NewReader(payload))
		if err != nil {
			log.Printf("error constructing new request: %s", err)
			return
		}
		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
		res, err := client.Do(req)
		if err != nil {
			log.Printf("error processing request to worker: %s", err)
			return
		}

		w.WriteHeader(res.StatusCode)
		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Printf("error reading response from worker: %s", err)
			return
		}
		res.Body.Close()
		_, err = w.Write(b)
		if err != nil {
			log.Printf("error writing response: %s", err)
			return
		}
	}
}
