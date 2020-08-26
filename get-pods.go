package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
)

type PageVariables struct {
	Pods []v1.Pod
}

func podHandler(w http.ResponseWriter, r *http.Request) {
	cs := getClientSet()

	// podVars := PageVariables{
	// 	Pods: getPods(cs),
	// }

	// t, err := template.ParseFiles("./assets/pod-list.html") //parse the html file homepage.html
	// if err != nil {                                         // if there is an error
	// 	log.Print("template parsing error: ", err) // log it
	// }
	// err = t.Execute(w, podVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	// if err != nil {             // if there is an error
	// 	log.Print("template executing error: ", err) //log it
	// }
	fmt.Println("Get pods")
	pods, _ := json.Marshal(getPods(cs))

	json.NewEncoder(w).Encode(string(pods))

}
