package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tenzi-muraru/paymentApp/payment"
)

const paymentServerHost = ":8080"

func main() {
	startApp(context.Background(), paymentServerHost)
}

func startApp(ctx context.Context, host string) {
	mongoURI := getEnvVar("MONGO_URI")
	repository, err := payment.NewDBPaymentRepository(mongoURI)
	if err != nil {
		log.Fatal(err)
	}
	handler := payment.NewHandler(repository)

	server := payment.NewServer(handler)
	go server.Run(host)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	// Block until signal received or context is closed
	select {
	case <-ch:
	case <-ctx.Done():
	}

	log.Println("Server shutting down...")
}

func getEnvVar(key string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		log.Fatalf("Missing mandatory environment variable: %s", key)
	}
	return value
}
