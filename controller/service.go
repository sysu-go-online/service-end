package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/service-end/types"
)

func checkFilePath(path string) bool {
	return true
}

// InitDockerConnection inits the connection to the docker service with the first message received from client
func initDockerConnection(msg string) (*websocket.Conn, error) {
	// Just handle command start with `go`
	conn, err := dialDockerService()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialDockerService create connection between web server and docker server
func dialDockerService() (*websocket.Conn, error) {
	// Set up websocket connection
	dockerAddr := os.Getenv("DOCKER_ADDRESS")
	dockerPort := os.Getenv("DOCKER_PORT")
	if len(dockerAddr) == 0 {
		dockerAddr = "localhost"
	}
	if len(dockerPort) == 0 {
		dockerPort = "8888"
	}
	dockerPort = ":" + dockerPort
	dockerAddr = dockerAddr + dockerPort
	url := url.URL{Scheme: "ws", Host: dockerAddr, Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// HandleMessage decide different operation according to the given json message
func handleMessage(mType int, msg []byte, conn *websocket.Conn, isFirst bool) error {
	var workSpace *Command
	var err error
	if isFirst {
		projectName := "test"
		username := "golang"
		pwd := getPwd(projectName, username)
		env := getEnv(projectName, username)
		workSpace = &Command{
			Command:     string(msg),
			PWD:         pwd,
			ENV:         env,
			UserName:    username,
			ProjectName: projectName,
		}
	}

	// Send message
	if isFirst {
		err = conn.WriteJSON(*workSpace)
	} else {
		err = conn.WriteMessage(mType, msg)
	}
	if err != nil {
		return err
	}
	return nil
}

// ReadFromClient receive message from client connection
func readFromClient(clientChan chan<- []byte, ws *websocket.Conn) {
	for {
		_, b, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Fprintln(os.Stderr, "Remote user closed the connection")
				ws.Close()
				close(clientChan)
				break
			}
			close(clientChan)
			fmt.Fprintln(os.Stderr, "Can not read message.")
			return
		}
		// fmt.Println(string(b))
		clientChan <- b
	}
}

// HandlerClientMsg handle message from client and send it to docker service
func handlerClientMsg(isFirst *bool, ws *websocket.Conn, sConn *websocket.Conn, msgType int, msg []byte) (conn *websocket.Conn) {
	// Init the connection to the docker serveice
	if *isFirst {
		tmp, err := initDockerConnection(string(msg))
		sConn = tmp
		if err != nil {
			panic(err)
		}
		// Listen message from docker service and send to client connection
		go sendMsgToClient(ws, sConn)
	}

	if sConn == nil {
		fmt.Fprintf(os.Stderr, "Invalid command.")
		ws.WriteMessage(msgType, []byte("Invalid Command"))
		ws.Close()
		conn = nil
		return
	}

	// Send message to docker service
	handleMessage(msgType, msg, sConn, *isFirst)
	*isFirst = false
	conn = sConn
	return
}

// SendMsgToClient send message to client
func sendMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn) {
	for {
		mType, msg, err := sConn.ReadMessage()
		if err != nil {
			// Server closed connection
			fmt.Fprintln(os.Stderr, "Docker service closed the connection")
			cConn.Close()
			return
		}
		cConn.WriteMessage(mType, msg)
	}
}

// getPwd return current path of given username
func getPwd(projectName string, username string) string {
	// Return user root in test version
	return ""
}

func getEnv(projectName string, username string) []string {
	env := []string{}
	env = append(env, "GOPATH")
	env = append(env, filepath.Join("/go", "src"))
	return env
}

// GetGithubAppMessages get github app id and secret and return it
func GetGithubAppMessages() (string, string) {
	ID := "635cfc9655dbb5bca4d4"
	Secret := "f9c08ae0227b645c88d87b999202d7f52828af97"
	return ID, Secret
}

// GetAccessToken get access_token from github and return it
func GetAccessToken(code, state string) (string, error) {
	// Post data to github for returned value
	id, secret := GetGithubAppMessages()
	client := &http.Client{}
	url := "https://github.com/login/oauth/access_token"
	jsonBody := types.GithubRequestBody{
		ClientID:     id,
		ClientSecret: secret,
		Code:         code,
		State:        state,
	}
	byteBody, err := json.Marshal(jsonBody)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	// Read data from response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	jsonResBody := new(types.GithubResponseBody)
	err = json.Unmarshal(body, jsonResBody)
	if err != nil {
		return "", err
	}
	return jsonResBody.AccessToken, nil
}

// GetUserMessage get user basic infomation from github and return it
func GetUserMessage(accessToken string) (*types.GithubUserDataResponse, error) {
	// Get user message from github service
	url := "https://api.github.com/user?access_token=" + accessToken
	req, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	jsonResBody := new(types.GithubUserDataResponse)
	err = json.Unmarshal(body, jsonResBody)
	if err != nil {
		return nil, err
	}
	return jsonResBody, nil
}
