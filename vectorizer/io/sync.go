package io

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"vectorizer/embeddings"

	"github.com/heussd/pdftotext-go"
)

var articleHeader *regexp.Regexp = regexp.MustCompile(`([Vv][Ii][Ss][Tt][AaOo]\s*(\.|:))|([Cc][Oo][Nn][Sd][Ii][Dd][Ee][Rr][Aa][Nn][Dd][Oo](\.|:)?)|([Aa][Rr][Tt](([ÍíIi][Cc][Uu][Ll][Oo])|\.?)\s*(\d+|([UuÚú][Nn][Ii][Cc][Oo]))\s*(\.|:))`)
var nonAlphaNumericNeitherSpace *regexp.Regexp = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
var lawNumberExpr *regexp.Regexp = regexp.MustCompile(`(\d+)(-|_|\s)+(\d+)`)

func ReadPdf(code embeddings.InstitutionCode, path string) (*embeddings.SourceDocument, error) {
	// Stringbuilder

	var buffer strings.Builder

	pdf, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %s", err)
	}

	pages, err := pdftotext.Extract(pdf)
	if err != nil {
		return nil, fmt.Errorf("failed to extract text from PDF: %s", err)
	}

	for _, page := range pages {
		buffer.WriteString(page.Content)
	}

	result := buffer.String()
	// Perform these replacements : hn -> ón, 10s -> los, 0s -> os, Ira -> 1ra, a1 -> al, ([aeiou])rn -> $1m, ([aeiou])0 -> $1o
	result = strings.ReplaceAll(result, "hn", "ón")
	result = strings.ReplaceAll(result, "10s", "los")
	result = strings.ReplaceAll(result, "0s", "os")
	result = strings.ReplaceAll(result, "s0", "so")
	result = strings.ReplaceAll(result, "Ira", "1ra")
	result = strings.ReplaceAll(result, "a1", "al")
	result = strings.ReplaceAll(result, "[lo]", "10")
	result = strings.ReplaceAll(result, "gobio", "gobierno")
	result = strings.ReplaceAll(result, "Gobio", "Gobierno")
	result = regexp.MustCompile(`([aeiou])rn`).ReplaceAllString(result, `$1m`)
	result = regexp.MustCompile(`([aeiou])0`).ReplaceAllString(result, `$1o`)

	return &embeddings.SourceDocument{
		FileName:        path[strings.LastIndex(path, string(os.PathSeparator))+1:],
		InstitutionCode: code, // Default institution code, can be modified later
		Content:         result,
	}, nil
}

func ExtractArticles(sourceDocument *embeddings.SourceDocument) (*embeddings.ParsedDocument, error) {
	switch sourceDocument.InstitutionCode {
	case embeddings.DGII:
		return extractDGIIArticles(sourceDocument)
	case embeddings.SB:
		return extractSBArticles(sourceDocument)
	case embeddings.SIMV:
		return extractSIMVrticles(sourceDocument)
	}
	return nil, fmt.Errorf("unsupported institution code: %d", sourceDocument.InstitutionCode)

}

func extractSBArticles(sourceDocument *embeddings.SourceDocument) (*embeddings.ParsedDocument, error) {
	return extractDGIIArticles(sourceDocument) // Placeholder for SB extraction logic
}

func extractSIMVrticles(sourceDocument *embeddings.SourceDocument) (*embeddings.ParsedDocument, error) {
	return extractDGIIArticles(sourceDocument) // Placeholder for SB extraction logic
}

// https://dgii.gov.do/legislacion/leyestributarias/Paginas/default.aspx
func extractDGIIArticles(sourceDocument *embeddings.SourceDocument) (*embeddings.ParsedDocument, error) {
	// hn -> ón, 10s -> los, 0s -> os, Ira -> 1ra, a1 -> al, ([aeiou])rn -> $1m, ([aeiou])0 -> $1o
	var concernIndex, viewIndex int = 0, 0
	indices := articleHeader.FindAllStringIndex(sourceDocument.Content, -1)
	if len(indices) == 0 {
		return nil, fmt.Errorf("no articles found in the document")
	}

	lawNumber, err := lawNumberFromTitle(sourceDocument.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to extract law number from title: %s", err)
	}
	parsedDocument := &embeddings.ParsedDocument{
		FileName:        lawNumber,
		InstitutionCode: sourceDocument.InstitutionCode,
		Segments:        make([]embeddings.Segment, 0, len(indices)),
	}

	for index, idxPair := range indices {
		title := strings.ToUpper(strings.TrimSpace(sourceDocument.Content[idxPair[0]:idxPair[1]]))

		var content string
		if index+1 < len(indices) {
			content = strings.TrimSpace(sourceDocument.Content[idxPair[1]:indices[index+1][0]])
		} else {
			content = strings.TrimSpace(sourceDocument.Content[idxPair[1]:])
		}
		parsedDocument.Segments = append(parsedDocument.Segments, embeddings.Segment{
			Title:   formatTitle(title, &concernIndex, &viewIndex),
			Content: content,
		})
	}

	return parsedDocument, nil
}

func lawNumberFromTitle(title string) (string, error) {
	if !lawNumberExpr.MatchString(title) {
		return "", fmt.Errorf("no law number found in title: %s", title)
	}
	matches := lawNumberExpr.FindStringSubmatch(title)
	if len(matches) < 3 {
		return "", fmt.Errorf("invalid law number format in title: %s", title)
	}
	lawNumber := "Ley " + matches[1] + "-" + matches[3] + ","

	return strings.ToUpper(lawNumber), nil
}

func formatTitle(title string, concernIndex *int, viewIndex *int) string {
	var trimmedString string = strings.ToUpper(strings.TrimSpace(title))
	var spaceIndex int = strings.LastIndex(trimmedString, " ")

	var formattedTitle string = nonAlphaNumericNeitherSpace.ReplaceAllString(trimmedString[spaceIndex+1:], "")

	if strings.HasPrefix(title, "CONSIDERANDO") {
		*concernIndex++
		formattedTitle = fmt.Sprintf("CONSIDERANDO %d", *concernIndex)

	} else if strings.HasPrefix(title, "VISTA") {
		*viewIndex++
		formattedTitle = fmt.Sprintf("VISTA %d", *viewIndex)

	} else if strings.HasPrefix(title, "VISTO") {
		*viewIndex++
		formattedTitle = fmt.Sprintf("VISTO %d", *viewIndex)

	} else {
		formattedTitle = fmt.Sprintf("ARTICULO %s", formattedTitle)
	}

	return formattedTitle
}
