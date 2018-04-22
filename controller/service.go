package controller

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
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
		pwd := getPwd("test")
		var env []string
		entrypoint := make([]string, 1) // Set `/go` as default entrypoint
		entrypoint[0] = "/go"
		username := "test"
		workSpace = &Command{
			Command:    string(msg),
			Entrypoint: entrypoint,
			PWD:        pwd,
			ENV:        env,
			UserName:   username,
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
		clientChan <- b
	}
}

// HandlerClientMsg handle message from client and send it to docker service
func handlerClientMsg(isFirst *bool, ws *websocket.Conn, msgType int, msg []byte) {
	var conn *websocket.Conn
	// Init the connection to the docker serveice
	if *isFirst {
		conn, err := initDockerConnection(string(msg))
		if err != nil {
			panic(err)
		}
		if conn == nil {
			fmt.Fprintf(os.Stderr, "Invalid command.")
			ws.WriteMessage(msgType, []byte("Invalid Command"))
			return
		}
		// Listen message from docker service and send to client connection
		go sendMsgToClient(ws, conn)
	}

	// Send message to docker service
	handleMessage(msgType, msg, conn, *isFirst)
	*isFirst = false
}

// SendMsgToClient send message to client
func sendMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn) {
	for {
		mType, msg, err := sConn.ReadMessage()
		if err != nil {
			return
		}
		cConn.WriteMessage(mType, msg)
	}
}

// GetPwd return current path of given username
func getPwd(username string) string {
	// Return user root in test version
	return "/"
}
