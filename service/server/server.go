package server

import (
	"apparatus/config"
	"apparatus/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func prompt(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var response service.ApparatusAnswer
	defer json.NewEncoder(writer).Encode(&response)

	userID := request.URL.Query().Get("user-id")
	question := request.URL.Query().Get("question")

	if userID == "" || question == "" {
		http.Error(writer, "Missing user_id or question", http.StatusBadRequest)
		return
	}

	answer, err := service.ApparatusService.Query(&service.ApparatusQuestion{
		UserID: userID,
		Prompt: question,
	})
	response.Answer = answer.Answer
	response.Score = answer.Score

	if err != nil {
		http.Error(writer, fmt.Sprintf("Error processing prompt: %v", err), http.StatusInternalServerError)
		return
	}

}

func Start() {
	address := fmt.Sprintf(":%s", config.Config.ServerPort)
	router := mux.NewRouter()
	router.HandleFunc("/prompt", prompt).Methods(http.MethodGet)
	http.ListenAndServe(fmt.Sprintf(address, config.Config.ServerPort), router)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
