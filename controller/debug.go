package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// ClientDebugMessage stores the data received from user
type ClientDebugMessage struct {
	BreakPoints []string
	Command     string
}

// ResponseDebugMessage stores the data to be sent to the client
type ResponseDebugMessage struct {
	Event       string
	CurrentLint string
	Info        string
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

// GDBOutput stores gdb output
type GDBOutput struct {
	Type    string `json:"type"`
	Class   string `json:"class"`
	Payload string `json:"payload"`
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
	msgType int, msg []byte) (conn *websocket.Conn) {
	// Init the connection to the docker serveice
	if *isFirst {
		tmp, err := initDockerConnection(string(msg), "debug")
		sConn = tmp
		if err != nil {
			panic(err)
		}
		// Listen message from docker service and send to client connection
		go sendDebugMsgToClient(ws, sConn)
	}

	if sConn == nil {
		fmt.Fprintf(os.Stderr, "Invalid command.")
		ws.WriteMessage(msgType, []byte("Invalid Command"))
		ws.Close()
		conn = nil
		return
	}

	// Send message to docker service
	handleDebugMessage(msgType, msg, sConn, *isFirst)
	*isFirst = false
	conn = sConn
	return
}

// handleDebugMessage parse data and send to docker server
func handleDebugMessage(mType int, msg []byte, conn *websocket.Conn, isFirst bool) error {
	var workSpace *Command
	var err error
	// send project data to docker server for preparing debug environment
	if isFirst {
		projectName := "test"
		username := "golang"
		workSpace = &Command{
			UserName:    username,
			ProjectName: projectName,
			Type:        "debug",
		}
		err = conn.WriteJSON(*workSpace)
		if err != nil {
			fmt.Println(err)
			return err
		}
		// err = conn.WriteMessage(mType, []byte("file-exec-and-symbols Debug/"+projectName+"\n"))
		// err = conn.WriteMessage(mType, []byte("file-exec-and-symbols Debug/"+"main"+"\n"))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	// parse client json data
	clientJSON := ClientDebugMessage{}
	err = json.Unmarshal(msg, &clientJSON)
	if err != nil {
		return err
	}
	// Send message
	msgsToBeSent := parseClientMessage(&clientJSON)
	for _, v := range msgsToBeSent {
		v += "\n"
		err = conn.WriteMessage(mType, []byte(v))
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func parseClientMessage(msg *ClientDebugMessage) []string {
	command := msg.Command
	switch command {
	case "run":
		return []string{"exec-run"}
	case "quit":
		return []string{"quit"}
	case "next":
		return []string{"exec-next"}
	case "continue":
		return []string{"exec-continue"}
	case "set":
		ret := []string{}
		// TODO: check breakpoint format
		for _, v := range msg.BreakPoints {
			ret = append(ret, "break-insert "+v)
		}
		return ret
	case "delete":
		// TODO: get break list from debug server and delete line
		return []string{}
	default:
		return []string{}
	}
}

func parseDebugInterface(msg map[string]interface{}, t string) []byte {
	var response ResponseDebugMessage
	if t == "output" {
		// output event
		response.Event = "output"
		response.Info = msg["msg"].(string)
	} else if t == "error" {
		// error event
		response.Event = "error"
		response.Info = msg["error"].(string)
	} else if t == "gdb" {
		if msg["class"] != nil {
			if msg["class"].(string) == "stopped" {
				payload := msg["payload"].(map[string]interface{})
				if len(payload) != 0 {
					reason := payload["reason"].(string)
					if reason == "exited-normally" {
						// finish event
						response.Event = "finish"
					} else if reason == "breakpoint-hit" {
						frame := payload["fullname"].(map[string]interface{})
						// TODO: parse relative path with fullname and file
						// fullname := frame["fullname"].(string)
						file := frame["fullname"].(string)
						line := frame["line"].(string)
						// stop event
						response.Event = "stop"
						response.CurrentLint = file + ":" + line
					}
				}
			} else if msg["class"].(string) == "thread-group-added" {
				// start event
				response.Event = "start"
			}
		}
	}
	if response.Event != "" {
		byteResponse, _ := json.Marshal(response)
		return byteResponse
	}
	return []byte{}
}
