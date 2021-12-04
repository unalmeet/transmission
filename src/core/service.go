package core

type ClientService interface {
	List(IdMeeting string) ([]*Client, error)
	Store(client *Client) error
	Delete(IdMeeting, IdSession string) error
	Audio(client *Client) (*Broadcast, error)
	Image(client *Client) (*Broadcast, error)
}
