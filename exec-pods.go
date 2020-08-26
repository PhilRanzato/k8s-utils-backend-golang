package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	v1 "k8s.io/api/core/v1"
)

type Connection struct {
	From          string
	FromNamespace string
	To            string
	ToNamespace   string
}

var result bool = false

type ExecCommand struct {
	PodName      string
	PodNamespace string
	Command      string
}

type ExecResult struct {
	CommandOutput string
	Error         bool
}

func podExecHandler(w http.ResponseWriter, r *http.Request) {
	var exec ExecCommand
	err := json.NewDecoder(r.Body).Decode(&exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
	}

	clientset := getClientSet()
	pods := getPods(clientset)

	// Search for pod and service to test connection on
	var pod v1.Pod

	for _, p := range pods {
		if exec.PodName == p.Name && exec.PodNamespace == p.Namespace {
			pod = p
		}
	}

	stdout, stderr, _ := execIntoPod(clientset, &pod, exec.Command, nil)

	var output ExecResult

	if stdout != "" {
		output.CommandOutput = string(stdout)
		fmt.Println(output.CommandOutput)
		output.Error = false
	} else {
		output.CommandOutput = string(stderr)
		output.Error = true
	}

	formattedOutput, _ := json.Marshal(output)
	fmt.Println(string(formattedOutput))
	json.NewEncoder(w).Encode(string(formattedOutput))

}

func TestConnection(w http.ResponseWriter, r *http.Request) {

	var conn Connection
	err := json.NewDecoder(r.Body).Decode(&conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err.Error())
		// return
	}

	clientset := getClientSet()
	pods := getPods(clientset)
	svcs := getServices(clientset)

	// Search for pod and service to test connection on
	var podFrom v1.Pod
	var serviceTo v1.Service

	for _, pod := range pods {
		if conn.From == pod.Name && conn.FromNamespace == pod.Namespace {
			podFrom = pod
		}
	}
	for _, svc := range svcs {
		if conn.To == svc.Name && conn.ToNamespace == svc.Namespace {
			serviceTo = svc
		}
	}
	result = testConnectionPodToService(podFrom, serviceTo)
	connResult, _ := json.Marshal(result)

	json.NewEncoder(w).Encode(string(connResult))
}
