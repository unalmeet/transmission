package api

import (
	"io/ioutil"
	"io"
	"log"
	"net/http"
	"github.com/go-chi/chi"
	"mime/multipart"
    "strconv"
	"bytes"
	js "ms/transmission/serializer/json"
	"ms/transmission/core"
)

type ClientHandler interface {
	DeleteSession(http.ResponseWriter, *http.Request)
	PostSession(http.ResponseWriter, *http.Request)
	PutAudio(http.ResponseWriter, *http.Request)
	PutImage(http.ResponseWriter, *http.Request)
	GetSession(http.ResponseWriter, *http.Request)
}

type handler struct {
	service core.ClientService
}

func NewHandler(clientService core.ClientService) ClientHandler {
	return &handler{service: clientService}
}

func setupResponse(writer http.ResponseWriter, contentType string, body []byte, statusCode int) {
	writer.Header().Set("Content-Type", contentType)
	writer.WriteHeader(statusCode)
	_, err := writer.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (handler *handler) serializer(contentType string) core.ClientSerializer {
	return &js.Client{}
}

func (handler *handler) GetSession(writer http.ResponseWriter, req *http.Request) {
	log.Println("INFO", "GetSession START")
	idMeeting := chi.URLParam(req, "idMeeting")
	contentType := "application/json"
	client, err := handler.service.List(idMeeting)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := handler.serializer(contentType).Encode(client)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("INFO", "GetSession FINISH")
	setupResponse(writer, contentType, responseBody, http.StatusOK)
}

func (handler *handler) PostSession(writer http.ResponseWriter, req *http.Request) {
	log.Println("INFO", "PostSession START")
	contentType := req.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("ERROR", "Error leyendo body")
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	client, err := handler.serializer(contentType).Decode(requestBody)
	if err != nil {
		log.Println("ERROR", "Error deserializando body")
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = handler.service.Store(client)
	if err != nil {
		log.Println("ERROR", "Error almacenando entidad")
		log.Println("ERROR", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := handler.serializer(contentType).Encode(client)
	if err != nil {
		log.Println("ERROR", "Error serializando body")
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("INFO", "PostSession START")
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}

func (handler *handler) DeleteSession(writer http.ResponseWriter, req *http.Request) {
	log.Println("INFO", "DeleteSession START")
	idSession := chi.URLParam(req, "idSession")
	idMeeting := chi.URLParam(req, "idMeeting")
	log.Println("DEBUG", idSession, idMeeting)
	err := handler.service.Delete(idMeeting, idSession)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("INFO", "DeleteSession START")
	writer.WriteHeader(http.StatusAccepted)
}

func (handler *handler) PutAudio(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestBody, err := req.MultipartReader()
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	client := new(core.Client)
	for {
		part, err := requestBody.NextPart()

		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "idmeeting":
			client.IdMeeting = string(readPart(part))
		case "idsession":
			client.IdSession, _ = strconv.Atoi(string(readPart(part)))
		case "data":
			client.Media = readPart(part)
		default:
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
	}
	broadcast, err := handler.service.Audio(client)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("DEBUG", "data", broadcast)
	responseBody, err := handler.serializer(contentType).Encode(broadcast)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}

func (handler *handler) PutImage(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	client, err := handler.serializer(contentType).Decode(requestBody)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = handler.service.Store(client)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := handler.serializer(contentType).Encode(client)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}

func readPart(part *multipart.Part) []byte {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, part)
	if err != nil {
		return []byte{}
	}
	return buf.Bytes()
}
