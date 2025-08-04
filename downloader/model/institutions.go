package model

import (
	"crypto/tls"
	"downloader/io"
	_ "embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type InstitutionInfo struct {
	Code     InstitutionCode
	Name     string
	FileList []string
}

var Institutions = map[InstitutionCode]InstitutionInfo{
	DGII: {DGII, "Dirección General de Impuestos Internos", []string{"resources/dgii-list.txt"}},
	SB:   {SB, "Superintendencia de Bancos", []string{"resources/sb-list.txt"}},
	SIMV: {SIMV, "Superintendencia de Valores y Seguros", []string{"resources/simv-list.txt"}},
}

func (institution *InstitutionInfo) Download(destinationFolder string) {
	if len(institution.FileList) == 0 {
		fmt.Printf("no files to download for institution %s", institution.Name)
		return
	}
	var waitGroup sync.WaitGroup
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	for _, uri := range institution.FileList {
		fileName := filepath.Base(uri)
		filePath := filepath.Join(destinationFolder, fileName)

		waitGroup.Add(1)
		go io.DownloadAsync(uri, filePath, &waitGroup)
	}
	waitGroup.Wait()

}

func init() {
	var institutions = &Institutions
	// get folder of current executable
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)

	// and set the path for the resources
	for code, info := range *institutions {
		FileList, err := io.ReadLines(filepath.Join(exeDir, info.FileList[0]))
		if err != nil {
			panic(err)
		}
		info.FileList = FileList
		(*institutions)[code] = info

		fmt.Println("INFO: model/institutions.go, ", info.Name, "with", len(info.FileList), "files")
	}
}
