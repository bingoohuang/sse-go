package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

var (
	mutex      = new(sync.Mutex)
	sseChannel SSEChannel
)

// SSEChannel model.
type SSEChannel struct {
	Clients  []chan string
	Notifier chan string
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.Println("SSE-GO")

	sseChannel = SSEChannel{
		Clients:  make([]chan string, 0),
		Notifier: make(chan string),
	}

	done := make(chan interface{})
	defer close(done)

	go broadcaster(done)

	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		h.Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Connection doesnot support streaming", http.StatusBadRequest)
			return
		}

		sseChan := make(chan string)
		mutex.Lock()
		sseChannel.Clients = append(sseChannel.Clients, sseChan)
		mutex.Unlock()

		// this go routine reads and streams into the channel
		d := make(chan interface{})
		defer close(d)
		defer log.Println("Closing channel.")

		for {
			select {
			case <-d:
				close(sseChan)
				return
			case data := <-sseChan:
				log.Printf("data: %v", data)
				fmt.Fprintf(w, "data: %v \n\n", data)
				flusher.Flush()
			}
		}
	})

	http.HandleFunc("/log", logHTTPRequest)

	log.Println("Listening to 5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Panicf("failed to listen on 5000: %v", err)
	}
}

func logHTTPRequest(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		log.Printf("Error: %v", err)
	}
	method := r.Method

	logMsg := fmt.Sprintf("Method: %v, Body: %v", method, buf.String())
	log.Println(logMsg)

	sseChannel.Notifier <- logMsg
}

func broadcaster(done <-chan interface{}) {
	log.Println("Broadcaster Started.")

	for {
		select {
		case <-done:
			return
		case data := <-sseChannel.Notifier:
			mutex.Lock()
			for _, channel := range sseChannel.Clients {
				channel <- data
			}
			mutex.Unlock()
		}
	}
}
