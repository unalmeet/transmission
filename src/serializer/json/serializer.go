package json

import (
	"encoding/json"
	"github.com/pkg/errors"
	"ms/transmission/core"
	"log"
)

type Client struct{}

func (c *Client) Decode(input []byte) (*core.Client, error) {
	client := &core.Client{}
	err := json.Unmarshal(input, client)
	if err != nil {
		log.Println("ERROR", "Error decodificando body json", err)
		return nil, errors.Wrap(err, "serializer.Client.Decode")
	}
	log.Println("DEBUG", client)
	return client, nil
}

func (c *Client) Encode(input interface{}) ([]byte, error) {
	rawMsg, err := json.Marshal(&input)
	if err != nil {
		log.Println("ERROR", "Error codificando body json", err)
		return nil, errors.Wrap(err, "serializer.Client.Encode")
	}
	return rawMsg, nil
}