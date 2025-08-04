package main

import (
	"downloader/model"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		os.Stderr.WriteString("Usage: downloader <DGII|SB|SIMV> <destination-folder>\n")
		for _, inst := range model.Institutions {
			code, _ := inst.Code.String()
			fmt.Printf("\t%s: %s\n", code, inst.Name)
		}
		os.Exit(-1)
	}

	institutionCode, ok := model.InstitutionAcronyms[args[0]]
	if !ok {
		os.Stderr.WriteString("Invalid institution code.:\n")
		os.Exit(-1)
	}

	destinationFolder := args[1]
	if _, err := os.Stat(destinationFolder); os.IsNotExist(err) {
		os.Stderr.WriteString("Destination folder does not exist.\n")
		os.Exit(-1)
	}

	info, err := os.Stat(destinationFolder)
	if err != nil || !info.IsDir() {
		os.Stderr.WriteString("Destination must be a valid directory.\n")
		os.Exit(-1)
	}
	institutionInfo := model.Institutions[institutionCode]
	institutionInfo.Download(destinationFolder)
	fmt.Printf("Files for %s downloaded successfully to %s\n", institutionInfo.Name, destinationFolder)
}
