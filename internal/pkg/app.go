package app

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"safechildhood/internal/app/handler"
	"safechildhood/internal/app/repository"
	"safechildhood/internal/app/service"
	"safechildhood/pkg/storage"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/schollz/progressbar/v3"
)

type App struct {
	repository *repository.Repository
	service    *service.Service
	handler    *handler.Handler
}

func New() *App {
	app := new(App)

	conn, err := pgx.Connect(context.Background(), "postgres://egzbmlsh:LShKLL4NFye8XhQdPx1I3jltNMkYLifH@cornelius.db.elephantsql.com/egzbmlsh")
	if err != nil {
		panic(err)
	}

	app.repository = repository.New(conn)

	playgrounds := service.NewPlaygroundsService(7 * time.Hour * 24)

	complaints := service.NewComplaintsService(app.repository.Complaints, playgrounds)

	googleDrive, err := storage.NewGoogleDrive(context.Background(), "./key.json")
	if err != nil {
		panic(err)
	}

	storage, err := service.NewStorageService(googleDrive)
	if err != nil {
		panic(err)
	}

	app.service = service.New(
		complaints,
		playgrounds,
		storage,
	)

	if errs := app.initPlaygroundsMap("./resources/playgrounds.csv"); len(errs) != 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
	}

	go app.service.Playgrounds.AutoCheckerFeatureTime()

	app.handler = handler.New(app.service)

	app.handler.Init()

	return app
}

func (a *App) initPlaygroundsMap(pathToResource string) []error {
	playgroundsMap := make(map[string]*service.MapProperties)

	errorsSlice := make([]error, 0)

	file, err := os.OpenFile(pathToResource, os.O_RDONLY, 0777)
	if err != nil {
		errorsSlice = append(errorsSlice, err)

		return errorsSlice
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

			errorsSlice = append(errorsSlice, err)

			return errorsSlice
		}

		playgroundsMap[data[0]] = &service.MapProperties{
			ID:      data[2],
			Color:   "green",
			Address: data[1],
		}
	}

	a.service.Playgrounds.SetPlaygroundsMap(playgroundsMap)

	a.initFoldersIdsMap(context.Background(), playgroundsMap)

	complaints, err := a.repository.Complaints.GetEarly(context.Background())
	if err != nil {
		errorsSlice = append(errorsSlice, err)

		return errorsSlice
	}

	a.service.Playgrounds.UpdatePlaygroundsMap(complaints)

	return []error{}
}

func (a *App) initFoldersIdsMap(ctx context.Context, playgroundsMap map[string]*service.MapProperties) []error {
	errors := make([]error, 0)

	bar := progressbar.Default(int64(len(playgroundsMap)))
	bar.Describe("initializtion coordinates folders")

	foldersIdsMap := make(map[string]string)

	for coordinates := range playgroundsMap {
		folder, err := a.service.Storage.Create(ctx, storage.GoogleDriveParameters{
			Name:       coordinates,
			ObjectMode: storage.FOLDER,
			ParentId:   "1JtHjonTau-gSkQd3Wj7wZ1Db7A8xHZDw",
		})
		if err != nil {
			errors = append(errors, err)

			foundFolder, err := a.service.Storage.GetByName(context.Background(), storage.GoogleDriveParameters{
				Name:       coordinates,
				ObjectMode: storage.FOLDER,
				ParentId:   "1JtHjonTau-gSkQd3Wj7wZ1Db7A8xHZDw",
			})
			if err != nil {
				errors = append(errors, err)

				goto updateBar
			}

			folder.Id = foundFolder[0].Id
		}

		foldersIdsMap[coordinates] = folder.Id

	updateBar:
		bar.Add(1)
	}

	fmt.Println(foldersIdsMap)

	a.service.Storage.SetFoldersIdsMap(foldersIdsMap)

	return errors
}

func (a *App) Run() {
	a.handler.Run(":8080")
}
