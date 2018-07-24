package controller

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/service-end/types"
)

func checkFilePath(path string) bool {
	return true
}

// InitDockerConnection inits the connection to the docker service with the first message received from client
func initDockerConnection(msg string, service string) (*websocket.Conn, error) {
	// Just handle command start with `go`
	conn, err := dialDockerService(service)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialDockerService create connection between web server and docker server
// Accept service type:
// tty debug
func dialDockerService(service string) (*websocket.Conn, error) {
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
	url := url.URL{Scheme: "ws", Host: dockerAddr, Path: "/" + service}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
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

// getPwd return current path of given username
func getPwd(projectName string, username string) string {
	// Return user root in test version
	return ""
}

func getEnv(projectName string, username string, language string) []string {
	env := []string{}
	if language == "golang" {
		env = append(env, "GOPATH=/root/go:/home/go")
	}
	return env
}

// GetConfigContent read configure file and return the content
func GetConfigContent() *types.ConfigFile {
	// Get messages from configure file
	configureFilePath := os.Getenv("CONFI_FILE_PATH")
	if len(configureFilePath) == 0 {
		configureFilePath = "/config/config.yml"
	}
	content, err := ioutil.ReadFile(configureFilePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	config := new(types.ConfigFile)
	err = yaml.Unmarshal(content, config)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return config
}
