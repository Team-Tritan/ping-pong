package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/team-tritan/ping-pong/internal/status"
)

func HandleProfileRequest(w http.ResponseWriter, r *http.Request) {

	host := r.URL.Query().Get("host")
	check_type := r.URL.Query().Get("type")

	fmt.Println(host, check_type)

	if check_type == "HTTP" {
		httpStatusChecker := status.NewHTTPStatus()

		status, err := httpStatusChecker.CheckStatus("GET", host, true)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			w.Write([]byte("Bad News"))
			return
		}

		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(status)
		return
	}
}

func NewServer(addr string) error {
	m := mux.NewRouter()

	m.Path("/profile_url").Methods("GET").HandlerFunc(HandleProfileRequest)

	server := &http.Server{
		Handler: m,
		Addr:    addr,
		// TODO: enforce timeouts
	}

	log.Printf("Listening on http://%s", server.Addr)

	return server.ListenAndServe()
}
