package main

import (
	"apparatus/initialization"
	"apparatus/model"
	"apparatus/server"
	"apparatus/service"
	"fmt"
	"os"
)

func serve() {
	initialization.Init()
	server.Start()
}

func prompt(userID, question string) {
	initialization.Init()
	if !model.IsValidClientId(userID) {
		fmt.Println("Invalid user ID.")
		os.Exit(-1)
	}
	fmt.Println("Is service.ApparatusService nil?", service.ApparatusService == nil)
	answer, error := service.ApparatusService.Query(&service.ApparatusQuestion{
		UserID: userID,
		Prompt: question,
	})

	if error != nil {
		fmt.Printf("Error processing prompt: %v\n", error)
		os.Exit(-1)
	}
	fmt.Println(answer.Answer)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go -((p <user_id> <question>)|s)")
		os.Exit(1)
	}

	mode := os.Args[1]
	switch mode {
	case "-p":
		if len(os.Args) < 4 {
			fmt.Println("Usage: go run main.go -p <user_id> <question>")
			os.Exit(1)
		}
		userID := os.Args[2]
		question := os.Args[3]
		fmt.Println("Prompt mode with user ID:", userID, "and question:", question)
		prompt(userID, question)

	case "-s":
		fmt.Println("Service mode")
		serve()

	default:
		fmt.Println("First argument must be '-p' or '-s'")
		os.Exit(-1)
	}
}
