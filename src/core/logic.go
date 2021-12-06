package core

import (
	"errors"
	"log"
	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
    "crypto/md5"
    "encoding/hex"
    "strconv"
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

func (service *clientService) List(token string) ([]*Client, error) {
	return service.repository.List(token)
}

func (service *clientService) Update(token string, idSession int, media []byte) error {
	return service.repository.Update(token, idSession, media)
}

func (service *clientService) Validate(token string) (bool, error) {
	return service.repository.Validate(token), nil
}

func (service *clientService) Store(client *Client) (*Client, error) {
	if err := validate.Validate(client)
	err != nil {
		log.Println("ERROR", "Error en Service Store", err)
		return nil, errs.Wrap(ErrClientInvalid, "service.Store")
	}
	hasher := md5.New()
    hasher.Write([]byte(client.IdMeeting + strconv.Itoa(client.IdUser)))
	client.Token = hex.EncodeToString(hasher.Sum(nil))
	err := service.repository.Store(client)
	if err != nil {
		log.Println("ERROR", "Error en Service Store", err)
		return nil, errs.Wrap(err, "service.Store")
	}
	return client, nil
}

func (service *clientService) Delete(token string) error {
	return service.repository.Delete(token)
}

func (service *clientService) Audio(broadcast *Broadcast) (*Broadcast, error) {
	clients, err := service.repository.List(broadcast.Token)
	if err != nil {
		log.Println("ERROR", "Error en Service Audio", err)
		return nil, errs.Wrap(ErrClientInvalid, "service.Audio")
	}
	for _, element := range clients {
		if broadcast.Token != element.Token{
			broadcast.IdSession = append(broadcast.IdSession, element.IdSession)
		}
	}
	return broadcast, nil
}

func (service *clientService) Image(broadcast *Broadcast) (*Broadcast, error) {
	var images [][]byte
	clients, err := service.repository.List(broadcast.Token)
	if err != nil {
		log.Println("ERROR", "Error en Service Audio", err)
		return nil, errs.Wrap(ErrClientInvalid, "service.Audio")
	}
	for _, element := range clients {
		broadcast.IdSession = append(broadcast.IdSession, element.IdSession)
		images = append(images, element.Media)
	}
	// change broadcast.media for processed image
	return broadcast, nil
}
