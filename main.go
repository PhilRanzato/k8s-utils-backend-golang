package main

import (
	"fmt"
	"net/http"

	// import an external library by running:
	// brew install dep
	// dep init
	// dep ensure --add github.com/gorilla/mux
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/hello", handler).Methods("GET")
	r.HandleFunc("/api/pods", podHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/services", serviceHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/pod-exec", podExecHandler).Methods("POST")
	r.HandleFunc("/api/test-connection", TestConnection).Methods("POST")
	r.HandleFunc("/api/pvc", pvcHandler).Methods("GET")
	return r
}

func main() {
	r := newRouter()
	fmt.Println("Server listening on port 80")
	// enable CORS
	router := handlers.CORS(handlers.AllowedHeaders([]string{"Accept", "X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)
	http.ListenAndServe(":80", router)
}
