package main

import (
	"log"
	"net/http"
	"wyoassign/wyoassign"

	"github.com/gorilla/mux"
)

func main() {
	wyoassign.InitAssignments()

	log.Println("starting API server")
	//create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	//specify endpoints
	router.HandleFunc("/", wyoassign.MainPage).Methods("GET")
	router.HandleFunc("/api-status", wyoassign.APISTATUS).Methods("GET")
	router.HandleFunc("/all_assignments", wyoassign.GetAssignments).Methods("GET")
	router.HandleFunc("/assignment/{id}", wyoassign.GetAssignment).Methods("GET")
	router.HandleFunc("/assignment/{id}", wyoassign.DeleteAssignment).Methods("DELETE")
	//router.HandleFunc("/assignment/{id}", wyoassign.DeleteAssignment).Queries("validationKey", "Valkey").Methods("DELETE")
	router.HandleFunc("/assignment", wyoassign.CreateAssignment).Methods("POST")
	router.HandleFunc("/assignment/{id}", wyoassign.UpdateAssignment).Methods("PUT")

	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)

}
