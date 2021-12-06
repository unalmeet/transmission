package core

type DataSerializer interface {
	DecodeClient(input []byte) (*Client, error)
	DecodeBroadcast(input []byte) (*Broadcast, error)
	Encode(input interface{}) ([]byte, error)
}
