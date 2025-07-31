package io

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

func DownloadAsync(uri, filePath string, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	fmt.Printf("Downloading %s to %s\n", uri, filePath)
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.88 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("failed to download %s: \n%v\n", uri, err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("failed to create file %s: \n%v\n", filePath, err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("failed to save file %s: \n%v\n", filePath, err)
		return
	}
	fmt.Printf("File %s downloaded successfully.\n", filePath)
}
