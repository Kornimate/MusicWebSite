package services

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	model "example.com/music-app/models"
	"github.com/beevik/guid"
	"gopkg.in/yaml.v3"
)

var PATH = ""

const CONFIG_PATH = "./config.yaml"

const NAME = "MusicApp_"
const COMPRESSED_EXTENSION = ".zip"

func InitializePath() {
	PATH = os.Getenv("DOWNLOAD_PATH")

	if PATH == "" {
		PATH = "." + string(filepath.Separator) + "temp"
	}
}

func DownloadMusic(url string) (string, string, bool, string, error) {
	tempFolderName := guid.New()

	folderPath := fmt.Sprintf("%v%v%v", PATH, string(filepath.Separator), tempFolderName)
	zipFile := NAME + guid.NewString() + COMPRESSED_EXTENSION

	err := os.MkdirAll(folderPath, 0777)

	if err != nil {
		defer CleanUp(folderPath)
		return "", "", false, "", fmt.Errorf("error while creating dir:%v", err)
	}

	cmd := exec.Command("python3", "./python/downloader.py", url, folderPath)

	output, err := cmd.CombinedOutput()

	if err != nil {
		defer CleanUp(folderPath)
		return "", "", false, string(output), fmt.Errorf("error while running command:%v\n%v", err, string(output))
	}

	fmt.Printf("Ouput: %v\n", string(output))

	fmt.Println("command executed")

	dir, err := os.Open(folderPath)

	if err != nil {
		defer CleanUp(folderPath)
		return "", "", false, string(output), fmt.Errorf("error while opening dir:%v", err)
	}

	dirInfoFiles, _ := dir.Readdir(0)

	if err := dir.Close(); err != nil {
		defer CleanUp(folderPath)
		return "", "", false, string(output), fmt.Errorf("error while reading dir content:%v", err)
	}

	size := len(dirInfoFiles)

	serverSideInfo := fmt.Sprintf("\nGo app looked at %v and found %v files", folderPath, size)

	if size > 1 {
		if err := createZipFile(folderPath, zipFile); err != nil {
			return "", "", false, string(output) + serverSideInfo, fmt.Errorf("failed to create zip file:%v", err)
		}

		fmt.Println("zip file created")

		return folderPath, zipFile, true, string(output), nil
	} else if size == 1 {
		fmt.Println("single file detected")

		return folderPath, dirInfoFiles[0].Name(), true, string(output) + serverSideInfo, nil
	} else {
		fmt.Println("no file detected")

		return "", "", false, string(output) + serverSideInfo, fmt.Errorf("no file detected for download")
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

func ScheduledCleanUp() {
	for {
		fmt.Println("clean up started")

		temp, err := os.ReadDir(PATH)

		if err == nil {
			for _, entry := range temp {
				info, err := entry.Info()

				if err != nil {
					continue
				}

				if info.ModTime().After(time.Now().Add(10 * time.Minute)) {
					os.RemoveAll(PATH + string(filepath.Separator) + entry.Name())
				}
			}
		}

		fmt.Println("clean up finished")

		time.Sleep(10 * time.Minute)
	}
}

func GetAppVersion() (string, error) {
	var config model.Config

	file, err := os.ReadFile(CONFIG_PATH)

	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(file, &config)

	if err != nil {
		return "", err
	}

	return config.App.Version, nil
}
