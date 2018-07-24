package controller

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

// HandlerClientMsg handle message from client and send it to docker service
func handlerClientTTYMsg(isFirst *bool, ws *websocket.Conn, sConn *websocket.Conn, msgType int, msg []byte) (conn *websocket.Conn) {
	// Init the connection to the docker serveice
	if *isFirst {
		tmp, err := initDockerConnection(string(msg), "tty")
		sConn = tmp
		if err != nil {
			panic(err)
		}
		// Listen message from docker service and send to client connection
		go sendTTYMsgToClient(ws, sConn)
	}

	if sConn == nil {
		fmt.Fprintf(os.Stderr, "Invalid command.")
		ws.WriteMessage(msgType, []byte("Invalid Command"))
		ws.Close()
		conn = nil
		return
	}

	// Send message to docker service
	handleTTYMessage(msgType, msg, sConn, *isFirst)
	*isFirst = false
	conn = sConn
	return
}

// SendMsgToClient send message to client
func sendTTYMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn) {
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

// HandleMessage decide different operation according to the given json message
func handleTTYMessage(mType int, msg []byte, conn *websocket.Conn, isFirst bool) error {
	var workSpace *Command
	var err error
	if isFirst {
		projectName := "test"
		username := "golang"
		pwd := getPwd(projectName, username)
		env := getEnv(projectName, username, "golang")
		workSpace = &Command{
			Command:     string(msg),
			PWD:         pwd,
			ENV:         env,
			UserName:    username,
			ProjectName: projectName,
			Type:        "tty",
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
