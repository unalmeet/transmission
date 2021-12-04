package core

type ClientSerializer interface {
	Decode(input []byte) (*Client, error)
	Encode(input interface{}) ([]byte, error)
}
