package app

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"safechildhood/internal/app/domain"
	"safechildhood/internal/app/handler"
	"safechildhood/internal/app/repository"
	"safechildhood/internal/app/service"
	"safechildhood/pkg/storage"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	repository *repository.Repository
	service    *service.Service
	handler    *handler.Handler
}

func New() *App {
	app := new(App)

	pool, err := pgxpool.New(context.Background(), "postgres://egzbmlsh:LShKLL4NFye8XhQdPx1I3jltNMkYLifH@cornelius.db.elephantsql.com/egzbmlsh")
	if err != nil {
		panic(err)
	}

	app.repository = repository.New(pool)

	playgrounds := service.NewPlaygroundsService(7 * time.Hour * 24)

	complaints := service.NewComplaintsService(app.repository.Complaints)

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

	app.handler = handler.New(app.service)

	app.handler.Init()

	go app.autoUpdatePlaygroundsMap()

	return app
}

func (a *App) autoUpdatePlaygroundsMap() {
	ticker := time.NewTicker(1 * time.Second)

	defer func(t *time.Ticker) {
		ticker.Stop()
	}(ticker)

	for range ticker.C {
		complaints, err := a.service.GetEarly(context.Background())
		if err != nil {
			log.Println(err)
		}

		a.service.Playgrounds.UpdatePlaygroundsMap(a.createPlaygroundsMapFromComplaints(complaints))
	}
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

	a.initFoldersIdsMap(context.Background())

	complaints, err := a.service.Complaints.GetEarly(context.Background())
	if err != nil {
		errorsSlice = append(errorsSlice, err)

		return errorsSlice
	}

	a.service.Playgrounds.UpdatePlaygroundsMap(a.createPlaygroundsMapFromComplaints(complaints))

	return []error{}
}

func (a *App) createPlaygroundsMapFromComplaints(complaints []domain.Complaint) map[string]*service.MapProperties {
	playgroundsMap := make(map[string]*service.MapProperties)

	for _, complaint := range complaints {
		playgroundsMap[complaint.Coordinates] = &service.MapProperties{
			Time: complaint.CreatedAt,
		}
	}

	return playgroundsMap
}

func (a *App) initFoldersIdsMap(ctx context.Context) error {
	foldersIdsMap := make(map[string]string)

	files, err := a.service.Storage.GetByParams(context.Background(), storage.GoogleDriveParameters{
		ParentId: "1cM704evigVIu8gAssGFmohdoHo5MH8Gs",
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		foldersIdsMap[file.Name] = file.Id
	}

	a.service.Storage.SetFoldersIdsMap(foldersIdsMap)

	return nil
}

func (a *App) Run() {
	a.handler.Run(":8080")
}
