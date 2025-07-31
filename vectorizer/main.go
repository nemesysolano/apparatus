package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"vectorizer/config"
	"vectorizer/embeddings"
	"vectorizer/io"
	"vectorizer/repository"

	"github.com/ollama/ollama/api"
)

type Arguments struct {
	pdfFilePath     string
	institutionCode embeddings.InstitutionCode
}

func arguments() (Arguments, error) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <path_to_pdf>")
		return Arguments{}, fmt.Errorf("insufficient arguments")
	}

	pdfFilePath := os.Args[1]
	code := os.Args[2]

	institutionCode, exists := embeddings.InstitutionAcronyms[code]
	if !exists {
		return Arguments{}, fmt.Errorf("invalid institution code: %s", code)
	}

	return Arguments{
		pdfFilePath:     pdfFilePath,
		institutionCode: institutionCode,
	}, nil
}

func main() {
	args, err := arguments()
	if err != nil {
		log.Fatalf("Error parsing arguments: %s", err)
		os.Exit(1)
	}

	sourceDocument, err := io.ReadPdf(args.institutionCode, args.pdfFilePath)
	if err != nil {
		log.Printf("Failed to read PDF file: %s", err)
		os.Exit(2)
	}

	parsedDocument, err := io.ExtractArticles(sourceDocument)
	if err != nil {
		log.Printf("Failed to extract articles: %s", err)
		os.Exit(3)
	}

	url, err := url.Parse(config.Config.OllamaHost)
	if err != nil {
		log.Printf("Invalid OLLAMA_HOST URL: %s", config.Config.OllamaHost)
		os.Exit(4)
	}

	client := api.NewClient(url, http.DefaultClient)

	emb, err := embeddings.EmbedParsedDocument(context.Background(), client, parsedDocument)
	if err != nil {
		log.Printf("Failed to embed document: %s", err)
		os.Exit(5)
	}

	err = repository.SaveAllEmbeddings(emb)
	if err != nil {
		log.Printf("Failed to save embeddings to database: %s", err)
		os.Exit(6)
	}

}
