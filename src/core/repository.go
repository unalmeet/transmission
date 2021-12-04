package core

type ClientRepository interface {
	List(IdMeeting string) ([]*Client, error)
	Store(client *Client) error
	Delete(IdMeeting, IdSession string) error
}
