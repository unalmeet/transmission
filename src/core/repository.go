package core

type ClientRepository interface {
	Store(client *Client) error
	List(token string) ([]*Client, error)
	Update(token string, idSession int, media []byte) error
	Delete(token string) error
	Validate(token string) bool
}
