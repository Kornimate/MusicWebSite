package services

import (
	_ "archive/zip"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/beevik/guid"
)

const PATH = "." + string(filepath.Separator) + "temp"
const NAME = "MusicApp_"
const COMPRESSED_EXTENSION = ".zip"

func DownloadMusic(url string, actionType string) (string, string, bool, error) {
	tempFolderName := guid.New()

	folderPath := fmt.Sprintf("%v%v%v", PATH, string(filepath.Separator), tempFolderName)
	zipFile := NAME + guid.NewString() + COMPRESSED_EXTENSION

	os.MkdirAll(folderPath, 0777)

	cmd := exec.Command(fmt.Sprintf("goffy -%v %v -d %v", actionType, url, folderPath))

	if err := cmd.Run(); err != nil {
		os.Remove(folderPath)
		return "", "", false, err
	}

	createZipFile(folderPath, zipFile)

	return folderPath, zipFile, false, nil
}

func createZipFile(path string, name string) {
	//making zip file: https://golang.cafe/blog/golang-zip-file-example.html
}

func CleanUp(path string) {
	os.RemoveAll(path)
}
