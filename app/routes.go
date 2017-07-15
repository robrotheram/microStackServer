package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"test/app/handlers"

	"github.com/gorilla/mux"
)

var wsMutex sync.Mutex
var ws *Websockify

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "TutorialEdge REST API: v0.0.1")
}

func IndexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}
func CORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELET")
}

func HandleWebsockify(ipport string, password string) string {
	var Url = "http://192.168.0.165:6080/vnc_auto.html?token=TOKEN&password=PASSWORD"
	wsMutex.Lock()
	defer wsMutex.Unlock()
	if ws == nil {
		ws = &Websockify{
			Debug:  true,
			Listen: ":6080",
		}
		ws.Run()
	}
	token := ws.Register(ipport)
	return strings.Replace(strings.Replace(Url, "TOKEN", token, 1), "PASSWORD", password, 1)
}

func VNCHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	host, port := handlers.GetVNC(uuid)
	var connection = uuid + " | " + host + ":" + strconv.Itoa(port)
	url := HandleWebsockify(host+":"+strconv.Itoa(port), "")
	fmt.Println(connection + " <br/> <a href='" + url + "'>link</a>")
	http.Redirect(w, r, url, 307)
}

var routes = Routes{
	Route{"Index", "GET", "/", IndexHandler("app/static/index.html")},
	Route{"Alltags", "GET", "/tags", handlers.HomePage},
	Route{"Register", "POST", "/users/register", handlers.Register},
	Route{"Login", "POST", "/users/login", handlers.Login},
	Route{"INFO", "GET", "/vms/info", handlers.Info},
	Route{"CREATE", "POST", "/vms/create", handlers.Create},
	Route{"STOP", "POST", "/vms/stop", handlers.Stop},
	Route{"START", "POST", "/vms/start", handlers.Start},
	Route{"Index", "GET", "/vms/console/", IndexHandler("app/static/noVNC/vnc_lite.html")},
	Route{"AUDIT", "POST", "/vms/audit", handlers.GetAudit},
	Route{"AUDIT", "OPTIONS", "/vms/audit", CORS},
	Route{"VNC", "GET", "/vms/{uuid}/vnc", VNCHandler},
}
