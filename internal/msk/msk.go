package msk

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"encoding/json"
	"fmt"
	"github/erastusk/tracer/internal/types"
	"github/erastusk/tracer/internal/utils"
	"math/rand"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type msk struct {
	topic     string
	partition int
	error     error
}

func NewMsk() *msk {
	return &msk{}
}

func (m *msk) Connectivity(opts types.PromptOptions) error {
	conn, err := m.createDialer(opts)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func (m *msk) ListTopics(opts types.PromptOptions) error {
	conn, err := m.createDialer(opts)
	if err != nil {
		return err
	}
	defer conn.Close()
	partitions, err := conn.ReadPartitions()
	if err != nil {
		// panic(err.Error())
		utils.Red.Println(err.Error())
		return err
	}

	mtopics := map[string]struct{}{}

	for _, p := range partitions {
		mtopics[p.Topic] = struct{}{}
	}
	for k := range mtopics {
		utils.Green.Println(k)
	}
	return nil
}

func (m *msk) Produce(opts types.PromptOptions, topic string) error {
	m.topic = topic
	c, err := m.createDialerLeader(opts)
	if err != nil {
		return err
	}
	c.SetWriteDeadline(time.Now().Add(10 * time.Second))
	m.topic = topic
	payload := types.UserEvent{
		UserID:    rand.Intn(100),
		Event:     "Tracer logs",
		Timestamp: time.Now(),
	}
	wg := &sync.WaitGroup{}
	produceMessage := func() {
		go func() {
			defer wg.Done()
			defer c.Close()
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				m.error = err
			}
			key := fmt.Sprintf("Key-%d", rand.Intn(10000))

			_, err = c.WriteMessages(kafka.Message{
				Key:   []byte(key),
				Value: payloadBytes,
			})
			if err != nil {
				m.error = err
			}
		}()
	}
	wg.Add(1)
	produceMessage()
	wg.Wait()
	if m.error != nil {
		return m.error
	}
	return nil
}

func (m *msk) Consume(opts types.PromptOptions, topic string) error {
	// c, err := m.createDialerReader(opts)
	m.topic = topic
	c, err := m.createDialerLeader(opts)
	if err != nil {
		utils.Red.Println("Failed to create dialer:", err)
		return err
	}
	c.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := c.ReadBatch(10e3, 1e6)
	b := make([]byte, 10e3)
	msglength := 10
	wg := &sync.WaitGroup{}
	readMessages := func() <-chan string {
		readChan := make(chan string, msglength)
		go func() {
			// r.SetOffset(12)
			defer c.Close()
			defer wg.Done()
			defer batch.Close()
			defer close(readChan)
			for {
				n, err := batch.Read(b)
				if err != nil {
					m.error = err
					break
				}
				if len(b[:n]) > 1 {
					readChan <- string(b[:n])
					break
				} else {
					m.error = fmt.Errorf("topic has no messages")
				}

			}
		}()
		return readChan
	}
	wg.Add(1)
	r := readMessages()
	wg.Wait()
	for x := range r {
		utils.Green.Println(x)
	}
	if m.error != nil {
		return m.error
	}
	return nil
}

func (m *msk) createDialer(opts types.PromptOptions) (*kafka.Conn, error) {
	mechanism, err := scram.Mechanism(scram.SHA512, opts.Username, opts.Password)
	if err != nil {
		return nil, err
	}
	tls_, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS:           tls_,
	}
	conn, err := dialer.DialContext(context.Background(), "tcp", opts.Endpoint)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (m *msk) createDialerLeader(opts types.PromptOptions) (*kafka.Conn, error) {
	mechanism, err := scram.Mechanism(scram.SHA512, opts.Username, opts.Password)
	if err != nil {
		return nil, err
	}
	tls_, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
		TLS:           tls_,
	}
	conn, err := dialer.DialLeader(context.Background(), "tcp", opts.Endpoint, m.topic, m.partition)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

//go:embed certs/AmazonRootCa.pem
var caCert string

func createTLSConfig() (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM([]byte(caCert)) {
		return nil, fmt.Errorf("failed to append CA certificate")
	}
	// Create TLS configuration
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		MinVersion: tls.VersionTLS12, // to ensure minimum TLS version
	}

	return tlsConfig, nil
}
