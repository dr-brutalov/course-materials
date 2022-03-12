package wyoassign

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Assignments []Assignment `json:"assignments"`
}

type Assignment struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Points      int    `json:"points"`
}

var Assignments []Assignment

const Valkey string = "FooKey"

func InitAssignments() {
	var assignment Assignment
	assignment.Id = "Mike1A"
	assignment.Title = "Lab 4 "
	assignment.Description = "Some lab this guy made yesterday?"
	assignment.Points = 20
	Assignments = append(Assignments, assignment)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to my assignment page!")
}

func APISTATUS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func GetAssignments(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	var response Response

	response.Assignments = Assignments

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		return
	}

	//TODO
	w.Write(jsonResponse)
}

func GetAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	params := mux.Vars(r)

	for _, assignment := range Assignments {
		if assignment.Id == params["id"] {
			json.NewEncoder(w).Encode(assignment)
			break
		} else {
			log.Printf("This assignment does not exist. Check the ID and try again or create a new assignment.")
			fmt.Fprintf(w, "This assignment does not exist.")
		}
	}
	//TODO : Provide a response if there is no such assignment
	//w.Write(jsonResponse)
}

func DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s DELETE end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/txt")

	params := mux.Vars(r)

	//key := r.URL.Query().Get("validationkey")

	response := make(map[string]string)
	//response["validationKey"] = Valkey

	//if key == Valkey {
	response["status"] = "No Such ID to Delete"
	for index, assignment := range Assignments {
		if assignment.Id == params["id"] {
			Assignments = append(Assignments[:index], Assignments[index+1:]...)
			response["status"] = "Success"
			break
		}
	}
	//}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(jsonResponse)
}

func UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	var response Response

	params := mux.Vars(r)
	response.Assignments = Assignments

	for _, assignment := range Assignments {
		if assignment.Id == params["id"] {
			assignment.Id = r.FormValue("id")
			assignment.Title = r.FormValue("title")
			assignment.Description = r.FormValue("desc")
			assignment.Points, _ = strconv.Atoi(r.FormValue("points"))
			DeleteAssignment(w, r)
			Assignments = append(Assignments, assignment)
			//w.WriteHeader(http.StatusCreated)
		} else {
			log.Printf("This assignment does not exist. Check the ID and try again or create a new assignment.")
			fmt.Fprintf(w, "This assignment does not exist.")
		}
	}

}

func CreateAssignment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var assignment Assignment
	r.ParseForm()
	// Possible TODO: Better Error Checking!
	// Possible TODO: Better Logging
	if r.FormValue("id") != "" {
		assignment.Id = r.FormValue("id")
		assignment.Title = r.FormValue("title")
		assignment.Description = r.FormValue("desc")
		assignment.Points, _ = strconv.Atoi(r.FormValue("points"))
		Assignments = append(Assignments, assignment)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
