package app

import (
	"fmt"
	"log"
	"net/http"
	"test/app/Config"

	"github.com/gorilla/mux"
)

/*
func handleWss(wsconn *websocket.Conn) {
	log.Println("handlewss")

	conn, err := net.Dial("tcp", ":5900")
	if err != nil {
		log.Println(err)
		wsconn.Close()
	} else {
		wsconn.PayloadType = websocket.BinaryFrame
		go io.Copy(conn, wsconn)
		go io.Copy(wsconn, conn)
		select {}
	}
}



func bootHandshake(config *websocket.Config, r *http.Request) error {
	//config.Protocol = []string{"binary"}
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
	return nil
}
*/

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	})
}

func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"inline.bundle.js": staticDirectory + "/inline.bundle.js",
		"bower_components": staticDirectory + "/bower_components/",
		"img":              staticDirectory + "/img/",
		"scripts":          staticDirectory + "/scripts/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}

func ServeVNCStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"vms/console/app/styles":           staticDirectory + "app/styles/",
		"vms/console/app/bower_components": staticDirectory + "/bower_components/",
		"vms/console/app/images":           staticDirectory + "app/images/",
		"vms/console/app/scripts":          staticDirectory + "/scripts/",
		"vms/console/vendor":               staticDirectory + "/vendor/",
		"vms/console/app":                  staticDirectory + "app/",
		"vms/console/core":                 staticDirectory + "core/",
	}
	for pathName, pathValue := range staticPaths {
		pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathPrefix).Handler(http.StripPrefix(pathPrefix,
			http.FileServer(http.Dir(pathValue))))
	}
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	//staticDirectory := "app/static/"
	//ServeStatic(router, staticDirectory)
	serveSingle("/inline.bundle.js", "./app/static/inline.bundle.js")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./app/static/")))
	ServeVNCStatic(router, "app/static/noVNC/")
	//router.PathPrefix("/vms/{id}/websockify").Handler(websocket.Server{Handshake: bootHandshake, Handler: handleWss})
	return router
}

type App struct {
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware", r.Method)

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
			w.Header().Set("Content-Type", "application/json")
			return
		}

		next.ServeHTTP(w, r)
		log.Println("Executing middleware again")
	})
}

func (a *App) Run() {
	conf := Config.Conf
	//	db := connect()
	fmt.Println("Listen and Server")

	log.Fatal(http.ListenAndServe(conf.URL, corsMiddleware(NewRouter())))

}
