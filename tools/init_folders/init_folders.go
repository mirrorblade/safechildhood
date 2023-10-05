package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"safechildhood/pkg/storage"

	"github.com/joho/godotenv"
	"github.com/schollz/progressbar/v3"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	coordinatesSlice := make([]string, 0)

	file, err := os.OpenFile(os.Getenv("PATH_TO_PLAYGROUNDS_FILE"), os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	csvReader.Read() //skip header

	for {
		data, err := csvReader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			panic(err)
		}

		coordinatesSlice = append(coordinatesSlice, data[0])
	}

	googleDrive, err := storage.NewGoogleDrive(context.Background(), os.Getenv("PATH_TO_GOOGLE_SERVICE_ACCOUNT"))
	if err != nil {
		panic(err)
	}

	progressBar := progressbar.Default(int64(len(coordinatesSlice)))
	progressBar.Describe("folders initialization")

	errorsSlice := make([]error, 0)

	for _, coordinates := range coordinatesSlice {
		if _, err := googleDrive.Create(context.Background(), storage.GoogleDriveParameters{
			Name:       coordinates,
			ObjectMode: storage.FOLDER,
			ParentId:   os.Getenv("MEDIA_FOLDER_ID"),
		}); err != nil {
			if errors.Is(err, storage.ErrObjectAlreadyExists) {
				errorsSlice = append(errorsSlice, err)

				goto updateBar
			}

			panic(err)
		}

	updateBar:
		progressBar.Add(1)
	}

	if len(errorsSlice) != 0 {
		fmt.Println(errorsSlice)
	}

	fmt.Println("successful folders initialization")
}
