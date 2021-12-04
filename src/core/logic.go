package core

import (
	"errors"
	"log"
	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrClientInvalid  = errors.New("Client data Invalid")
)

type clientService struct {
	repository ClientRepository
}

func NewRedirectService(clientRepo ClientRepository) ClientService {
	return &clientService{
		repository: clientRepo,
	}
}

func (service *clientService) List(IdMeeting string) ([]*Client, error) {
	return service.repository.List(IdMeeting)
}

func (service *clientService) Store(client *Client) error {
	if err := validate.Validate(client)
	err != nil {
		log.Println("ERROR", "Error en Service Store", err)
		return errs.Wrap(ErrClientInvalid, "service.Store")
	}
	return service.repository.Store(client)
}

func (service *clientService) Delete(IdMeeting, IdSession string) error {
	return service.repository.Delete(IdMeeting, IdSession)
}

func (service *clientService) Audio(client *Client) (*Broadcast, error) {
	var broadcast = new(Broadcast)
	clients, err := service.repository.List(client.IdMeeting)
	if err != nil {
		log.Println("ERROR", "Error en Service Audio", err)
		return nil, errs.Wrap(ErrClientInvalid, "service.Audio")
	}
	for _, element := range clients {
		if client.IdSession != element.IdSession{
			broadcast.IdSession = append(broadcast.IdSession, element.IdSession)
		}
	}
	broadcast.IdMeeting = client.IdMeeting
	broadcast.Media = client.Media
	return broadcast, nil
}

func (service *clientService) Image(client *Client) (*Broadcast, error) {
	var broadcast = new(Broadcast)
	var images [][]byte
	clients, err := service.repository.List(client.IdMeeting)
	if err != nil {
		log.Println("ERROR", "Error en Service Audio", err)
		return nil, errs.Wrap(ErrClientInvalid, "service.Audio")
	}
	for _, element := range clients {
		broadcast.IdSession = append(broadcast.IdSession, element.IdSession)
		images = append(images, client.Media)
	}
	broadcast.IdMeeting = client.IdMeeting
	//process media
	broadcast.Media = client.Media
	return broadcast, nil
}
