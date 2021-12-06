package json

import (
	"encoding/json"
	"github.com/pkg/errors"
	"ms/transmission/core"
	"log"
)

type JsonSerializer struct{}

func NewSerializer() core.DataSerializer {
	serializer := &JsonSerializer{}
	return serializer
}

func (serializer *JsonSerializer) DecodeClient(input []byte) (*core.Client, error) {
	client := &core.Client{}
	err := json.Unmarshal(input, client)
	if err != nil {
		log.Println("ERROR", "Error decodificando body json", err)
		return nil, errors.Wrap(err, "serializer.Client.Decode")
	}
	return client, nil
}

func (serializer *JsonSerializer) DecodeBroadcast(input []byte) (*core.Broadcast, error) {
	broadcast := &core.Broadcast{}
	err := json.Unmarshal(input, broadcast)
	if err != nil {
		log.Println("ERROR", "Error decodificando body json", err)
		return nil, errors.Wrap(err, "serializer.Broadcast.Decode")
	}
	return broadcast, nil
}

func (serializer *JsonSerializer) Encode(input interface{}) ([]byte, error) {
	rawMsg, err := json.Marshal(&input)
	if err != nil {
		log.Println("ERROR", "Error codificando body json", err)
		return nil, errors.Wrap(err, "serializer.Encode")
	}
	return rawMsg, nil
}