package controller

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/sysu-go-online/service-end/model/service"

	"github.com/sysu-go-online/service-end/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	dao "github.com/sysu-go-online/service-end/model/service"
)

var username = "golang"

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
		err := dao.UpdateFileContent(projectID, username, filePath, string(body), false, false)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// CreateFileHandler is a handler for create file
func CreateFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	type IsDir struct {
		dir bool
	}
	// Judge if it is dir from body
	isDir := IsDir{}
	err = json.Unmarshal(body, &isDir)
	if err != nil {
		return err
	}
	dir := isDir.dir

	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)

	if ok {
		// Save file
		err := dao.UpdateFileContent(projectID, username, filePath, "", true, dir)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}
	return nil
}

// GetFileContentHandler is a handler for reading file content
func GetFileContentHandler(w http.ResponseWriter, r *http.Request) error {
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)
	if ok {
		// Load file
		content, err := dao.GetFileContent(projectID, username, filePath)
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

// GetFileContentHandler is a handler for read file content
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) error {
	// Read project id and file path from uri
	vars := mux.Vars(r)
	projectID := vars["projectid"]
	filePath := vars["filepath"]

	// Check if the file path is valid
	ok := checkFilePath(filePath)
	if ok {
		// Load file
		err := dao.DeleteFile(projectID, username, filePath)
		if err != nil {
			return err
		}
		w.WriteHeader(200)
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
	structure, err := dao.GetFileStructure(projectID, username)
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
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	// Set TextMessage as default
	msgType := websocket.TextMessage
	clientMsg := make(chan []byte)
	if err != nil {
		panic(err)
	}
	defer ws.Close()

	// Open a goroutine to receive message from client connection
	go readFromClient(clientMsg, ws)

	go func() {
		for {
			timer := time.NewTimer(time.Second * 2)
			<-timer.C
			err := ws.WriteControl(websocket.PingMessage, []byte("ping"), time.Time{})
			if err != nil {
				timer.Stop()
				return
			}
		}
	}()

	// Handle messages from the channel
	isFirst := true
	var sConn *websocket.Conn
	for msg := range clientMsg {
		conn := handlerClientMsg(&isFirst, ws, sConn, msgType, msg)
		sConn = conn
	}
	sConn.Close()
}

func AuthUserHandler(w http.ResponseWriter, r *http.Request) error {
	// Get code and state from client
	r.ParseForm()
	code := r.FormValue("code")
	state := r.FormValue("state")
	if len(code)*len(state) == 0 {
		return errors.New("Incomplete form value")
	}
	accessToken, err := GetAccessToken(code, state)
	if err != nil {
		return nil
	}
	userInfo, err := GetUserMessage(accessToken)
	if err != nil {
		return err
	}
	// Check user data in the database
	user := service.GetUserInformation(userInfo.Username)
	if user.Name == "" {
		// TODO: Add this user to the db
	}
	ret := types.AuthResponse{
		Name: user.Name,
		Icon: user.Icon,
	}
	byteRetBody, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	w.Write(byteRetBody)
	return nil
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}
