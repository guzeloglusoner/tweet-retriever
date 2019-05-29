package main

import (
	"log"
	"net/http"
	"os"
)

var (
	// web server configurations
	serverPort = ":9090" // pass this as an environment variable from docker compose file
	// kafka configurations
	kafkaBrokerURL = []string{"localhost:29092"}
	kafkaClientID  = "web-server-consumer"
	kafkaTopic     = "tweet"
)

var logger = log.New(os.Stdout, "main: ", log.LstdFlags)

func main() {
	hub := initHub()
	go hub.start()

	// TODO move this part into a separate file
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(hub, w, r)
	})

	go ConsumeKafkaTopic(hub)

	logger.Printf("The server is listening on port %v", serverPort)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
