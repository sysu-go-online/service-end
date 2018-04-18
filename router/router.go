package router

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sysu-go-online/service-end/controller"
	"github.com/urfave/negroni"
)

var upgrader = websocket.Upgrader{}

// GetServer return web server
func GetServer() *negroni.Negroni {
	// formatter := render.New(render.Options{
	// 	IndentJSON: true,
	// })

	r := mux.NewRouter()
	// static := "static"
	// Define static service
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(static))))

	// /projects router
	// r.HandleFunc("/projects", controller.CreateProjects).Methods("POST")
	// r.HandleFunc("/projects", controller.GetProjectsID(formatter)).Methods("GET")
	// r.HandleFunc("/test", controller.TestGetProjectsID(formatter)).Methods("GET")

	// file router
	r.Handle("/{projectid}/tree/{filepath:.*}", controller.ErrorHandler(controller.UpdateFileHandler)).Methods("POST")
	r.Handle("/{projectid}/tree/{filepath:.*}", controller.ErrorHandler(controller.GetFileContentHandler)).Methods("GET")
	r.Handle("/{projectid}", controller.ErrorHandler(controller.GetFileStructureHandler)).Methods("POST")
	r.HandleFunc("/ws", controller.WebSocketTermHandler)

	// Use classic server and return it
	s := negroni.Classic()
	s.UseHandler(r)
	return s
}
