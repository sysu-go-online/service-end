package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	dao "github.com/sysu-go-online/service-end/model/service"
	"github.com/sysu-go-online/service-end/tools"
)

var upgrader = websocket.Upgrader{}

// UpdateFileHandler is a handler for update file
func UpdateFileHandler(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tools.HandlerError(w, err, "Can not read user message", 400)
		return
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
			tools.HandlerError(w, err, "Can not update file content", 500)
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
}

// GetFileContentHandler is a handler for read file content
func GetFileContentHandler(w http.ResponseWriter, r *http.Request) {
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
			tools.HandlerError(w, err, "Can not read content", 500)
		}
		w.WriteHeader(200)
		w.Write(content)
	} else {
		w.WriteHeader(400)
	}
}

// GetFileStructureHandler is handler for get project structure
func GetFileStructureHandler(w http.ResponseWriter, r *http.Request) {
	// Read project id
	vars := mux.Vars(r)
	projectID := vars["projectid"]

	// Get file structure
	structure, err := dao.GetFileStructure(projectID)
	if err != nil {
		tools.HandlerError(w, err, "Can not get required file structure", 500)
	}
	ret, err := json.Marshal(structure)
	if err != nil {
		tools.HandlerError(w, err, "Can not marshal json to string", 500)
	}
	w.Write(ret)
}

// WebSocketTermHandler is a middle way handler to connect web app with docker service
func WebSocketTermHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	// Set TextMessage as default
	msgType := websocket.TextMessage
	clientMsg := make(chan []byte)
	if err != nil {
		tools.HandlerError(w, err, "Can not upgrade http connection to ws", 500)
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
