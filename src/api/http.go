package api

import (
	"io/ioutil"
	"io"
	"log"
	"net/http"
	"github.com/go-chi/chi"
	"mime/multipart"
	"bytes"
	js "ms/transmission/serializer/json"
	"ms/transmission/core"
)

type ClientHandler interface {
	GetSession(http.ResponseWriter, *http.Request)
	PostSession(http.ResponseWriter, *http.Request)
	PutSession(http.ResponseWriter, *http.Request)
	DeleteSession(http.ResponseWriter, *http.Request)
	PostAudio(http.ResponseWriter, *http.Request)
	PostImage(http.ResponseWriter, *http.Request)
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

func (handler *handler) serializer(contentType string) core.DataSerializer {
	return js.NewSerializer()
}

func (handler *handler) GetSession(writer http.ResponseWriter, req *http.Request) {
	log.Println("INFO", "GetSession START")
	token := chi.URLParam(req, "token")
	contentType := "application/json"
	client, err := handler.service.List(token)
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
	client, err := handler.serializer(contentType).DecodeClient(requestBody)
	if err != nil {
		log.Println("ERROR", "Error deserializando body")
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	client, err = handler.service.Store(client)
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
	log.Println("INFO", "PostSession FINISH")
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}

func (handler *handler) PutSession(writer http.ResponseWriter, req *http.Request) {
	log.Println("INFO", "PutSession START")
	contentType := req.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("ERROR", "Error leyendo body")
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	client, err := handler.serializer(contentType).DecodeClient(requestBody)
	if err != nil {
		log.Println("ERROR", "Error deserializando body")
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = handler.service.Update(client.Token, client.IdSession, nil)
	if err != nil {
		log.Println("ERROR", "Error actualizando entidad")
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
	log.Println("INFO", "PutSession FINISH")
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}

func (handler *handler) DeleteSession(writer http.ResponseWriter, req *http.Request) {
	log.Println("INFO", "DeleteSession START")
	token := chi.URLParam(req, "token")
	err := handler.service.Delete(token)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("INFO", "DeleteSession FINISH")
	writer.WriteHeader(http.StatusAccepted)
}

func (handler *handler) PostAudio(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestBody, err := req.MultipartReader()
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	broadcast := new(core.Broadcast)
	for {
		part, err := requestBody.NextPart()

		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "token":
			broadcast.Token = string(readPart(part))
		case "data":
			broadcast.Media = readPart(part)
		default:
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
	}
	res, err := handler.service.Audio(broadcast)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("DEBUG", "data", res)
	responseBody, err := handler.serializer(contentType).Encode(res)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(writer, contentType, responseBody, http.StatusCreated)
}

func (handler *handler) PostImage(writer http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	requestBody, err := req.MultipartReader()
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	broadcast := new(core.Broadcast)
	for {
		part, err := requestBody.NextPart()

		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "token":
			broadcast.Token = string(readPart(part))
		case "data":
			broadcast.Media = readPart(part)
		default:
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		
	}
	res, err := handler.service.Image(broadcast)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Println("DEBUG", "data", res)
	responseBody, err := handler.serializer(contentType).Encode(res)
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
