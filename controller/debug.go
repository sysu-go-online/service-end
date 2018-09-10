package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// ClientDebugMessage stores the data received from user
type ClientDebugMessage struct {
	BreakPoints string
	Command     string
}

// ResponseDebugMessage stores the data to be sent to the client
type ResponseDebugMessage struct {
	Event       string
	CurrentLine string
	Message     string
	AddOn       string
}

// VaribleInformation stores information of varible
type VaribleInformation struct {
	Name  string
	Type  string
	Value string
}

// DebugOutPut stores message and type from debug service
type DebugOutPut struct {
	Type    string                 `json:"type"`
	Message map[string]interface{} `json:"msg"`
}

// DebugHandler is a websocket connection and handle debug information
func DebugHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	// Set TextMessage as default
	msgType := websocket.TextMessage
	clientMsg := make(chan RequestCommand)
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

	// Init docker service connection
	isFirst := true
	sConn := handleClientDebugMessage(&isFirst, ws, nil, msgType, []byte{}, RequestCommand{})

	// Handle messages from the channel
	for msg := range clientMsg {
		conn := handleClientDebugMessage(&isFirst, ws, sConn, msgType, []byte(msg.Command), msg)
		sConn = conn
	}
	sConn.Close()
}

// SendDebugMsgToClient send debug message to client
func sendDebugMsgToClient(cConn *websocket.Conn, sConn *websocket.Conn) {
	for {
		mType, msg, err := sConn.ReadMessage()
		if err != nil {
			// Server closed connection
			fmt.Fprintln(os.Stderr, "Docker service closed the connection")
			cConn.Close()
			return
		}
		// parse, handle data and sent to the client
		msgArr := strings.Split(string(msg), "\n")
		for _, v := range msgArr {
			if len(v) >= 2 {
				debugMsg := DebugOutPut{}
				err = json.Unmarshal([]byte(v), &debugMsg)
				if err != nil {
					fmt.Println(err)
					continue
				}
				message := parseDebugInterface(debugMsg.Message, debugMsg.Type)
				if len(message) == 0 {
					continue
				}
				cConn.WriteMessage(mType, message)
			}
		}
	}
}

func handleClientDebugMessage(isFirst *bool, ws *websocket.Conn, sConn *websocket.Conn,
	msgType int, msg []byte, connectContext RequestCommand) (conn *websocket.Conn) {
	// Init the connection to the docker serveice
	if *isFirst {
		tmp, err := initDockerConnection(string(msg), "debug")
		if err != nil {
			panic(err)
		}
		sConn = tmp
		// Listen message from docker service and send to client connection
		go sendDebugMsgToClient(ws, sConn)
	}
	// fmt.Println(string(msg))

	if sConn == nil {
		fmt.Fprintf(os.Stderr, "Invalid command.")
		ws.WriteMessage(msgType, []byte("Invalid Command"))
		ws.Close()
		conn = nil
		return
	}

	// Send message to docker service
	handleDebugMessage(msgType, msg, sConn, *isFirst, connectContext)
	*isFirst = false
	conn = sConn
	return
}

// handleDebugMessage parse data and send to docker server
func handleDebugMessage(mType int, msg []byte, conn *websocket.Conn, isFirst bool, connectContext RequestCommand) error {
	var workSpace *Command
	var err error
	// send project data to docker server for preparing debug environment
	if isFirst {
		workSpace = &Command{
			UserName:    connectContext.username,
			ProjectName: connectContext.Project,
			Type:        "debug",
			Command:     "/main",
		}
		err = conn.WriteJSON(*workSpace)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}
	// parse client json data
	clientJSON := ClientDebugMessage{}
	err = json.Unmarshal(msg, &clientJSON)
	if err != nil {
		return err
	}
	// Send message
	msgsToBeSent := parseClientMessage(&clientJSON)
	err = conn.WriteMessage(mType, []byte(msgsToBeSent+"\n"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func parseClientMessage(msg *ClientDebugMessage) string {
	command := msg.Command
	switch command {
	case "run":
		return "exec-run"
	case "quit":
		return "quit"
	case "next":
		return "exec-next"
	case "continue":
		return "exec-continue"
	case "set":
		// TODO: check breakpoint format
		return "break-insert " + msg.BreakPoints
	case "delete":
		// TODO: check breakpoint format
		return "break-delete " + msg.BreakPoints
	default:
		return ""
	}
}

func parseDebugInterface(msg map[string]interface{}, t string) []byte {
	var response ResponseDebugMessage
	var command string
	if msg["command"] != nil {
		command = msg["command"].(string)
	}
	if t == "output" {
		// output event
		response.Event = "output"
		response.Message = msg["msg"].(string)
	} else if t == "error" {
		// error event
		response.Event = "error"
		if command == "file-exec-and-symbols" {
			response.Message = "load"
		} else if len(command) >= 12 && command[:12] == "break-insert" {
			response.Message = "bk-add"
		}
	} else if t == "gdb" {
		if command == "" {
			if msg["class"] != nil {
				if msg["class"].(string) == "stopped" {
					stopDebugEvent(msg, &response)
				}
			}
		} else {
			// TODO: handle failure condition
			if len(command) >= 12 && command[:12] == "break-insert" {
				breakPointAddedEvent(msg, &response)
			} else if command == "file-exec-and-symbols" {
				if msg["class"] != nil && msg["class"].(string) == "done" {
					response.Event = "done"
					response.Message = "loaded"
				}
			} else if len(command) >= 14 && command[:12] == "break-delete" {
				if msg["class"] != nil && msg["class"].(string) == "done" {
					response.Event = "done"
					response.Message = "bk-del"
					response.AddOn = command[13:]
				}
			} else if command == "exec-run" || command == "exec-next" || command == "exec-continue" {
				if msg["class"] != nil && msg["class"].(string) == "running" {
					response.Event = "done"
					response.Message = "running"
				}
			}
		}
	}
	if response.Event != "" {
		byteResponse, _ := json.Marshal(response)
		return byteResponse
	}
	return []byte{}
}

func breakPointAddedEvent(msg map[string]interface{}, response *ResponseDebugMessage) {
	if msg["class"] == nil {
		return
	}
	if msg["class"].(string) == "done" {
		response.Event = "done"
		response.Message = "bk-add"
		payload := msg["payload"].(map[string]interface{})
		if payload != nil {
			// breakpoint event
			bkpt := payload["bkpt"].(map[string]interface{})
			if bkpt != nil {
				// TODO: parse relative path with fullname and file
				fullname := bkpt["fullname"].(string)
				// file := bkpt["file"].(string)
				line := bkpt["line"].(string)
				number := bkpt["number"].(string)
				response.CurrentLine = fullname + ":" + line
				response.AddOn = number
			}
		}
	}
}

func stopDebugEvent(msg map[string]interface{}, response *ResponseDebugMessage) {
	payload := msg["payload"].(map[string]interface{})
	if len(payload) != 0 {
		reason := payload["reason"].(string)
		if reason == "exited-normally" {
			// finish event
			response.Event = "finish"
		} else if reason == "breakpoint-hit" {
			frame := payload["frame"].(map[string]interface{})
			// TODO: parse relative path with fullname and file
			fullname := frame["fullname"].(string)
			// file := frame["file"].(string)
			line := frame["line"].(string)
			// stop event
			response.Event = "stop"
			response.CurrentLine = fullname + ":" + line
		} else if reason == "end-stepping-range" {
			// stop event
			frame := payload["frame"].(map[string]interface{})
			// TODO: parse relative path with fullname and file
			fullname := frame["fullname"].(string)
			// file := frame["file"].(string)
			line := frame["line"].(string)
			response.Event = "stop"
			response.CurrentLine = fullname + ":" + line
		}
	}
}
