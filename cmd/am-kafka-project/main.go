package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"am-kafka-project/internal/http/server"
	"am-kafka-project/internal/kafka/producer"
	"am-kafka-project/pkg/version"
)

func main() {
	fmt.Println(version.Info())        // Print the version information
	p := producer.NewProducerService() // Create a new producer service
	err := p.Configure()               // Configure the producer service
	if err != nil {
		fmt.Println(err) // If there is an error while configuring the producer service, print the error and exit
		os.Exit(1)
	}
	defer p.Close()

	// Start a goroutine to handle producer events
	go func() {
		for e := range p.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition) // If there is an error while delivering a message, print the error
				}
			}
		}
	}()

	var errChan = make(chan error, 1) // Create a channel for errors

	// Start a goroutine to start the HTTP server
	go func() {
		fmt.Println("Starting server")
		errChan <- server.StartHTTPServer(*p) // Start the HTTP server and send any errors to the error channel
	}()

	var signalChan = make(chan os.Signal, 1)                 // Create a channel for signals
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM) // Notify the signal channel of interrupt and SIGTERM signals

	// Wait for a signal or an error
	select {
	case <-signalChan:
		fmt.Println("got an interrupt, exiting...") // If an interrupt signal is received, print a message and exit
	case err := <-errChan:
		if err != nil {
			fmt.Println("error while running api, exiting...") // If an error is received, print the error and exit
		}
	}
}
