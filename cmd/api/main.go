package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	kafkaBootstrap := os.Getenv("BOOTSTRAP_SERVERS")
	kafkaTopics := os.Getenv("TOPICS")
	kafkaGroup := os.Getenv("GROUP_ID")

	if kafkaBootstrap == "" || kafkaTopics == "" {
		log.Fatal("BOOTSTRAP_SERVERS and TOPICS must be set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, kafkaBootstrap, kafkaTopics, kafkaGroup)
	})

	server := &http.Server{
		Addr: ":8000",
	}

	go func() {
		log.Println("Server started on :8000")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server shutdown error:", err)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request, broker, topic, groupID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("WebSocket close error:", err)
		}
	}(conn)

	// Create Kafka reader (consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // 1B
		MaxBytes: 10e6, // 10MB
	})
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			log.Println("Kafka close error:", err)
		}
	}(reader)

	ctx := context.Background()
	log.Println("WebSocket connected and Kafka consumer started")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println("Kafka read error:", err)
			time.Sleep(time.Second)
			continue
		}

		log.Printf("Consumed: topic=%s partition=%d offset=%d value=%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

		if err := conn.WriteMessage(websocket.TextMessage, msg.Value); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}
