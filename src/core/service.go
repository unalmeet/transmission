package core

type ClientService interface {
	Store(client *Client) (*Client, error)
	List(token string) ([]*Client, error)
	Update(token string, idSession int, media []byte) error
	Delete(token string) error
	Validate(token string) (bool, error)
	Audio(broadcast *Broadcast) (*Broadcast, error)
	Image(broadcast *Broadcast) (*Broadcast, error)
}
