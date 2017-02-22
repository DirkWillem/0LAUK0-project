package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Initialize router
	r := mux.NewRouter()

	r.HandleFunc("/api/authenticate", HandleAuthenticate).Methods("POST")
	r.HandleFunc("/api/authenticatedispenser", HandleAuthenticateDispenser).Methods("POST")

	r.HandleFunc("/api/medications", CheckJWT(CheckRole(DoctorOrPharmacist, HandleCreateMedication))).Methods("POST")
	r.HandleFunc("/api/medications", CheckJWT(CheckRole(DoctorOrPharmacist, HandleListMedications))).Methods("GET")
	r.HandleFunc("/api/medications/{medicationId}", CheckJWT(CheckRole(DoctorOrPharmacist, HandleReadMedication))).Methods("GET")
	r.HandleFunc("/api/medications/{medicationId}", CheckJWT(CheckRole(DoctorOrPharmacist, HandleUpdateMedication))).Methods("PUT")
	r.HandleFunc("/api/medications/{medicationId}", CheckJWT(CheckRole(DoctorOrPharmacist, HandleDeleteMedication))).Methods("DELETE")

	r.HandleFunc("/api/users", CheckJWT(CheckRole(Doctor, HandleCreateUser))).Methods("POST")
	r.HandleFunc("/api/users", CheckJWT(CheckRole(Admin, HandleListUsers))).Methods("GET")
	r.HandleFunc("/api/users/{userId}", CheckJWT(CheckRole(Admin, HandleReadUser))).Methods("GET")
	r.HandleFunc("/api/users/{userId}", CheckJWT(CheckRole(Doctor, HandleUpdateUser))).Methods("PUT")
	r.HandleFunc("/api/users/{userId}", CheckJWT(CheckRole(Admin, HandleDeleteUser))).Methods("DELETE")

	r.HandleFunc("/api/users/{userId}/doses", CheckJWT(CheckRole(Doctor, HandleCreateDose))).Methods("POST")
	r.HandleFunc("/api/users/{userId}/doses", CheckJWT(HandleListDoses)).Methods("GET")
	r.HandleFunc("/api/users/{userId}/doses/{doseId}", CheckJWT(HandleReadDose)).Methods("GET")
	r.HandleFunc("/api/users/{userId}/doses/{doseId}", CheckJWT(CheckRole(Doctor, HandleUpdateDose))).Methods("PUT")
	r.HandleFunc("/api/users/{userId}/doses/{doseId}", CheckJWT(CheckRole(Doctor, HandleDeleteDose))).Methods("DELETE")

	// Start web server
	log.Printf("Listening on %s:%s", config.Host.Host, config.Host.Port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host.Host, config.Host.Port), r)

	if err != nil {
		LogErrorMessageFatal(err.Error())
	}
}
