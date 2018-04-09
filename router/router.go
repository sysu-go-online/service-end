package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/service-end/controller"
	"github.com/urfave/negroni"
)

var upgrader = websocket.Upgrader{}

// GetServer return web server
func GetServer() *negroni.Negroni {
	r := mux.NewRouter()
	static := "static"
	// Define static service
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	// /projects router
	r.HandleFunc("/projects", controller.CreateProjects).Methods("POST")

	// Use classic server and return it
	s := negroni.Classic()
	s.UseHandler(r)
	return s
}
