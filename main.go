package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func tlsConfig() *tls.Config {

	cert, err := tls.LoadX509KeyPair("/client.crt", "/client.key")
	if err != nil {
		log.Fatal(err.Error())
	}
	//Truststore
	caCert, err := ioutil.ReadFile("/server.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true,
	}
	return config
}

func getMechanism(name, user, pw string) (sasl.Mechanism, error) {
	switch name {
	case "none":
		return nil, nil
	case "scram-sha-256":
		return scram.Mechanism(scram.SHA256, user, pw)
	default:
		return scram.Mechanism(scram.SHA512, user, pw)
	}
}

func getDialer() *kafka.Dialer {

	mechanism, err := getMechanism("scram-sha-256", "user", "E6tOW0zM9F")
	if err != nil {
		log.Fatal("failed to get mechanism:", err)
	}
	dialer := &kafka.Dialer{
		Timeout:       5 * time.Second,
		DualStack:     true,
		TLS:           tlsConfig(),
		SASLMechanism: mechanism,
	}

	return dialer
}

func main() {
	// to produce messages
	topic := "bmo-events"
	partition := 0
	dialer := getDialer()
	ctx := context.Background()
	conn, err := dialer.DialLeader(ctx, "tcp", "10.43.147.235:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatal("failed to set writer deadline:", err)
	}
	log.Println("Writing messages")
	out, err := conn.WriteMessages(kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")})
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	log.Println("Successfully pushed message", out)
	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
