package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"main/utils"
	"main/dispatch"
	"os"
	"gopkg.in/hlandau/passlib.v1"
)


func fileHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat("./static" + r.URL.Path); err != nil {
		http.ServeFile(w, r, "./static/index.html")
		return
	}
	http.ServeFile(w, r, "./static"+r.URL.Path)
}

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
	r.HandleFunc("/api/users", CheckJWT(CheckRole(Doctor, HandleListUsers))).Methods("GET")
	r.HandleFunc("/api/users/{userId}", CheckJWT(CheckRole(Doctor, HandleReadUser))).Methods("GET")
	r.HandleFunc("/api/users/{userId}", CheckJWT(CheckRole(Doctor, HandleUpdateUser))).Methods("PUT")
	r.HandleFunc("/api/users/{userId}", CheckJWT(CheckRole(Doctor, HandleDeleteUser))).Methods("DELETE")

	r.HandleFunc("/api/users/{userId}/doses", CheckJWT(CheckRole(Doctor, HandleCreateDose))).Methods("POST")
	r.HandleFunc("/api/users/{userId}/doses", CheckJWT(HandleListDoses)).Methods("GET")
	r.HandleFunc("/api/users/{userId}/doses/{doseId}", CheckJWT(HandleReadDose)).Methods("GET")
	r.HandleFunc("/api/users/{userId}/doses/{doseId}", CheckJWT(CheckRole(Doctor, HandleUpdateDose))).Methods("PUT")
	r.HandleFunc("/api/users/{userId}/doses/{doseId}", CheckJWT(CheckRole(Doctor, HandleDeleteDose))).Methods("DELETE")

	r.HandleFunc("/api/users/{userId}/dosehistory", CheckJWT(CheckRole(Dispenser, HandleCreateDoseHistoryEntry))).Methods("POST")
	r.HandleFunc("/api/users/{userId}/dosehistory", CheckJWT(CheckRole(Doctor, HandleListDoseHistoryEntries))).Methods("GET")
	r.HandleFunc("/api/users/{userId}/dosehistory/{doseHistoryEntryId}", CheckJWT(CheckRole(Doctor, HandleReadDoseHistoryEntry))).Methods("GET")

	r.HandleFunc("/api/users/{userId}/dosesummaries", CheckJWT(CheckRole(Doctor, HandleListDoseSummaries))).Methods("GET")
	r.HandleFunc("/api/users/{userId}/dosesummaries/{date}", CheckJWT(CheckRole(Doctor, HandleReadDoseSummary))).Methods("GET")

	r.HandleFunc("/api/dispatcher", dispatch.CreateDispatchHandler(dispatcher)).Methods("GET")

	r.PathPrefix("/").HandlerFunc(fileHandler)

	// Start the dispatcher
	go dispatcher.Start()

	// Start web server
	log.Printf("Listening on %s:%s", config.Host.Host, config.Host.Port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host.Host, config.Host.Port), r)

	if err != nil {
		utils.LogErrorMessageFatal(err.Error())
	}
}
