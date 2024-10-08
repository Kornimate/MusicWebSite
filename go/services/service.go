package services

import (
	"archive/zip"
	"fmt"
	"io"
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

	err := os.MkdirAll(folderPath, 0777)

	if err != nil {
		defer CleanUp(folderPath)
		return "", "", false, err
	}

	cmd := exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "-o", folderPath+string(filepath.Separator)+"%(title)s.%(ext)s", url)

	if err := cmd.Run(); err != nil {
		defer CleanUp(folderPath)
		return "", "", false, err
	}

	fmt.Println("command executed")

	dir, err := os.Open(folderPath)

	if err != nil {
		defer CleanUp(folderPath)
		return "", "", false, err
	}

	dirInfoFiles, _ := dir.Readdir(0)

	if err := dir.Close(); err != nil {
		defer CleanUp(folderPath)
		return "", "", false, err
	}

	size := len(dirInfoFiles)

	if size > 1 {
		if err := createZipFile(folderPath, zipFile); err != nil {
			return "", "", false, err
		}

		fmt.Println("zip file created")

		return folderPath, zipFile, true, nil
	} else if size == 1 {
		fmt.Println("single file detected")

		return folderPath, dirInfoFiles[0].Name(), true, nil
	} else {
		fmt.Println("no file detected")

		return "", "", false, fmt.Errorf("no file detected for download")
	}
}

func createZipFile(path string, name string) error {
	zipFilePath := path + string(filepath.Separator) + name

	archive, err := os.Create(zipFilePath)

	if err != nil {
		defer CleanUp(path)
		return err
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	musicDir, err := os.Open(path)

	if err != nil {
		defer CleanUp(path)
		return err
	}

	files, err := musicDir.Readdir(0)

	if err != nil {
		defer CleanUp(path)
		return err
	}

	for _, file := range files {

		if file.Name() == name {
			continue
		}

		tempFileOpen, err := os.Open(path + string(filepath.Separator) + file.Name())

		if err != nil {
			defer CleanUp(path)
			return err
		}

		zipEntry, err := zipWriter.Create(file.Name())

		if err != nil {
			defer CleanUp(path)
			return err
		}

		if _, err := io.Copy(zipEntry, tempFileOpen); err != nil {
			defer CleanUp(path)
			return err
		}

		if err := tempFileOpen.Close(); err != nil {
			defer CleanUp(path)
			return err
		}
	}

	if err := zipWriter.Close(); err != nil {
		defer CleanUp(path)
		return err
	}

	return nil
}

func CleanUp(path string) {
	os.RemoveAll(path)
}
