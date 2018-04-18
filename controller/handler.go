package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	dao "github.com/sysu-go-online/service-end/model/service"
)

var upgrader = websocket.Upgrader{}

// UpdateFileHandler is a handler for update file
func UpdateFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)

	if ok {
		// Save file
		err := dao.UpdateFileContent(projectID, filePath, string(body))
		if err != nil {
			return err
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// GetFileContentHandler is a handler for read file content
func GetFileContentHandler(w http.ResponseWriter, r *http.Request) error {
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)
	if ok {
		// Load file
		content, err := dao.GetFileContent(projectID, filePath)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
		w.Write(content)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// GetFileStructureHandler is handler for get project structure
func GetFileStructureHandler(w http.ResponseWriter, r *http.Request) error {
	// Read project id
	vars := mux.Vars(r)
	projectID := vars["projectid"]

	// Get file structure
	structure, err := dao.GetFileStructure(projectID)
	if err != nil {
		return err
	}
	ret, err := json.Marshal(structure)
	if err != nil {
		return err
	}
	w.Write(ret)
	return nil
}

// WebSocketTermHandler is a middle way handler to connect web app with docker service
func WebSocketTermHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	// Set TextMessage as default
	msgType := websocket.TextMessage
	clientMsg := make(chan []byte)
	if err != nil {
		return
	}
	defer ws.Close()

	// Open a goroutine to receive message from client connection
	go readFromClient(clientMsg, ws)

	// Handle messages from the channel
	isFirst := true
	for msg := range clientMsg {
		err := handlerClientMsg(&isFirst, ws, msgType, msg)
		if err != nil {
			return
		}
	}
}
