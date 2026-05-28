package nats

import (
	"fmt"

	natsgo "github.com/nats-io/nats.go"
)

func New(url string) (*natsgo.Conn, error) {
	conn, err := natsgo.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}

	return conn, nil
}
