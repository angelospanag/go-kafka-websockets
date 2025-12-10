package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
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
	autoOffset := os.Getenv("AUTO_OFFSET_RESET") // "earliest" or "latest"

	if kafkaBootstrap == "" || kafkaTopics == "" {
		log.Fatal("BOOTSTRAP_SERVERS and TOPICS must be set")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, kafkaBootstrap, kafkaTopics, kafkaGroup, autoOffset)
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
	err = server.Shutdown(ctx)
	if err != nil {
		return
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request, bootstrap, topics, group, offset string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Failed to close WebSocket connection:", err)
		}
	}(conn)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap,
		"group.id":          group,
		"auto.offset.reset": offset,
	})
	if err != nil {
		log.Println("Failed to create Kafka consumer:", err)
		return
	}
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Println("Failed to close Kafka consumer:", err)
		}
	}(consumer)

	err = consumer.SubscribeTopics([]string{topics}, nil)
	if err != nil {
		log.Println("Failed to subscribe to topics:", err)
		return
	}

	log.Println("WebSocket connected and Kafka consumer started")

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			// Errors are informational; continue
			log.Println("Consumer error:", err)
			continue
		}

		log.Printf("Consumed: topic=%s partition=%d offset=%d key=%s value=%s timestamp=%v\n",
			*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset,
			string(msg.Key), string(msg.Value), msg.Timestamp)

		if err := conn.WriteMessage(websocket.TextMessage, msg.Value); err != nil {
			log.Println("WebSocket write error:", err)
			break
		}
	}
}
