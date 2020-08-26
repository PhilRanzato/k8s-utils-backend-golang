package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Pvc struct {
	Name    string
	Status  string
	Storage string
}

type PageVars struct {
	Pvcs []Pvc
}

func pvcHandler(w http.ResponseWriter, r *http.Request) {

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create an rest client not targeting specific API version
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	var ns, label, field string

	api := clientset.CoreV1()
	// setup list options
	listOptions := metav1.ListOptions{
		LabelSelector: label,
		FieldSelector: field,
	}
	pvcs, err := api.PersistentVolumeClaims(ns).List(listOptions)
	if err != nil {
		log.Fatal(err)
	}
	var pvcsList []Pvc
	for _, pvc := range pvcs.Items {
		quant := pvc.Spec.Resources.Requests[v1.ResourceStorage]
		pvcObj := Pvc{
			pvc.Name,
			string(pvc.Status.Phase),
			quant.String()}
		pvcsList = append(pvcsList, pvcObj)
	}

	pvcVars := PageVars{
		Pvcs: pvcsList,
	}

	t, err := template.ParseFiles("./assets/pvc.html") //parse the html file homepage.html
	if err != nil {                                    // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, pvcVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {             // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
