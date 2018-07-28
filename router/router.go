package router

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/sysu-go-online/service-end/controller"
	"github.com/urfave/negroni"
)

var upgrader = websocket.Upgrader{}

// GetServer return web server
func GetServer() *negroni.Negroni {
	r := mux.NewRouter()

	// websocket handler
	ws := r.PathPrefix("/ws").Subrouter()
	ws.HandleFunc("/tty", controller.WebSocketTermHandler)
	ws.HandleFunc("/debug", controller.DebugHandler)

	// subrouter
	users := r.PathPrefix("/users").Subrouter()
	sessions := r.PathPrefix("/sessions").Subrouter()
	projects := users.PathPrefix("/{username}/projects").Subrouter()
	files := projects.PathPrefix("/{projectname}/files").Subrouter()

	// user collection
	users.Handle("", controller.ErrorHandler(controller.CreateUserHandler)).Methods("POST")
	users.Handle("/{username}", controller.ErrorHandler(controller.GetUserMessageHandler)).Methods("GET")

	// session handler
	sessions.Handle("", controller.ErrorHandler(controller.LogInHandler)).Methods("POST")
	sessions.Handle("", controller.ErrorHandler(controller.LogOutHandler)).Methods("DELETE")
	// project collection

	// file collection
	files.Handle("", controller.ErrorHandler(controller.GetFileStructureHandler)).Methods("GET")
	files.Handle("/{filepath:.*}", controller.ErrorHandler(controller.GetFileContentHandler)).Methods("GET")
	files.Handle("/{filepath:.*}", controller.ErrorHandler(controller.UpdateFileHandler)).Methods("POST")
	files.Handle("/{filepath:.*}", controller.ErrorHandler(controller.CreateFileHandler)).Methods("PUT")
	files.Handle("/{filepath:.*}", controller.ErrorHandler(controller.DeleteFileHandler)).Methods("DELETE")

	// Use classic server and return it
	handler := cors.Default().Handler(r)
	s := negroni.Classic()
	s.UseHandler(handler)
	return s
}
