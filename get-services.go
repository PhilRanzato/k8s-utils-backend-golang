package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	cs := getClientSet()

	fmt.Println("Get services")

	svcs, _ := json.Marshal(getServices(cs))

	json.NewEncoder(w).Encode(string(svcs))

}
