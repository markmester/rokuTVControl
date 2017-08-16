package rokuAPI

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strings"
	"sync"
)

func PowerEndpoint(w http.ResponseWriter, req *http.Request) {
	out := map[string]string{
		"success": "true",
		"data": "",
	}

	redis_ctx := *NewRedisClient()
	roku_addr := redis_ctx.Get("roku_address")

	if roku_addr != "" {
		// first try out API endoint
		resp := Request(roku_addr, "/keypress/Power", "post")

		if !strings.Contains(resp, "HTTP/1.1 200 OK") {
			out["success"] = "false"
		}
		out["data"] = resp

	// if that fails (likely b/c the tv powered off and disconnected from the network)
	// turn on using LIRC

	}

	json.NewEncoder(w).Encode(out)
}

func MuxServer(wg *sync.WaitGroup) {
	defer wg.Done()

	router := mux.NewRouter()
	router.HandleFunc("/power", PowerEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", router))
}