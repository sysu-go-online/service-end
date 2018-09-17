package controller

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/service-end/model"
)

// WebSocketTermHandler is a middle way handler to connect web app with docker service
func WebSocketTermHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	// Set TextMessage as default
	msgType := websocket.TextMessage
	clientMsg := make(chan RequestCommand)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	// Open a goroutine to receive message from client connection
	go readFromClient(clientMsg, ws)

	// keep connection
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
		conn := handlerClientTTYMsg(&isFirst, ws, sConn, msgType, &msg)
		sConn = conn
	}
	sConn.Close()
}

// HandlerClientMsg handle message from client and send it to docker service
func handlerClientTTYMsg(isFirst *bool, ws *websocket.Conn, sConn *websocket.Conn, msgType int, connectContext *RequestCommand) (conn *websocket.Conn) {
	// Init the connection to the docker serveice
	type res struct {
		Err string `json:"err"`
		Msg string `json:"msg"`
	}
	r := res{}
	if *isFirst {
		// check token
		ok, username := GetUserNameFromToken(connectContext.JWT)
		connectContext.username = username
		if !ok {
			fmt.Fprintln(os.Stderr, "Can not get user token information")
			r.Err = "Invalid token"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}

		// Get project information
		session := MysqlEngine.NewSession()
		u := model.User{Username: username}
		ok, err := u.GetWithUsername(session)
		if !ok {
			fmt.Fprintln(os.Stderr, "Can not get user information")
			r.Err = "Invalid user information"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.Err = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		p := model.Project{Name: connectContext.Project, UserID: u.ID}
		has, err := p.GetWithUserIDAndName(session)
		if !has {
			fmt.Fprintln(os.Stderr, "Can not get project information")
			r.Err = "Can not get project information"
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			r.Err = err.Error()
			ws.WriteJSON(r)
			ws.Close()
			conn = nil
			return
		}
		connectContext.projectType = p.Language

		tmp, err := initDockerConnection("tty")
		sConn = tmp
		if err != nil {
			fmt.Println("Can not connect to the docker service")
			r.Err = "Server error"
			ws.WriteJSON(r)
			ws.Close()
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
	handleTTYMessage(msgType, sConn, *isFirst, connectContext)
	*isFirst = false
	conn = sConn
	return
}

// SendMsgToClient send message to client
func sendTTYMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn) {
	type res struct {
		Err string `json:"err"`
		Msg string `json:"msg"`
	}
	for {
		_, msg, err := sConn.ReadMessage()
		fmt.Print(string(msg))
		r := res{}
		if err != nil {
			// Server closed connection
			r.Err = err.Error()
			cConn.WriteJSON(r)
			cConn.Close()
			return
		}
		r.Msg = string(msg)
		cConn.WriteJSON(r)
	}
}

// HandleMessage decide different operation according to the given json message
func handleTTYMessage(mType int, conn *websocket.Conn, isFirst bool, connectContext *RequestCommand) error {
	var workSpace *Command
	var err error
	if isFirst {
		projectName := connectContext.Project
		username := connectContext.username
		pwd := getPwd(projectName, username, connectContext.projectType)
		env := getEnv(projectName, username, connectContext.projectType)
		workSpace = &Command{
			Command:     connectContext.Command,
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
		err = conn.WriteMessage(mType, []byte(connectContext.Command))
	}
	if err != nil {
		return err
	}
	return nil
}
