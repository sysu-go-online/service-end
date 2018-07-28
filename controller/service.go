package controller

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/sysu-go-online/service-end/model"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/rs/xid"
	"github.com/sysu-go-online/service-end/types"
)

// JWTKey defines the token key
var JWTKey = "go-online"

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

// CheckEmail check if the email is valid
func CheckEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

// HashPassword return hash of password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CompasePassword compare raw password with hashed one
func CompasePassword(raw, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw)) == nil
}

// GenerateUserName generate unique userid
func GenerateUserName() string {
	return "user_" + generateUUID()
}

func generateUUID() string {
	guid := xid.New()
	return guid.String()
}

// CheckUsername check if the username is valid
func CheckUsername(username string) bool {
	if len(username) > 0 {
		return true
	}
	return false
}

// CheckJWT check whether the jwt is valid and if it is in the invalid database
func CheckJWT(jwtString string) (bool, error) {
	isValid, err := ValidateToken(jwtString)
	if err != nil {
		return false, err
	}
	if !isValid {
		return false, nil
	}

	has, err := model.IsJWTExist(jwtString, RedisClient)
	return !has, err
}

// ValidateToken check the format of token
func ValidateToken(jwtString string) (bool, error) {
	// validate jwt
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWTKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// parse time from jwt
	var exp int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expired := claims["exp"]
		if expired == nil {
			return false, nil
		}
		exp = int64(expired.(float64))
		if time.Now().Unix() > exp {
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

// GenerateJWT generate token for user
func GenerateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		"sub": email,
		"iat": time.Now().Unix(),
		"jti": generateUUID(),
	})

	return token.SignedString([]byte(JWTKey))
}
