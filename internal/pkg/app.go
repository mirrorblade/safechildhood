package app

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"safechildhood/internal/app/config"
	"safechildhood/internal/app/domain"
	"safechildhood/internal/app/handler"
	"safechildhood/internal/app/repository"
	"safechildhood/internal/app/service"
	"safechildhood/pkg/storage"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	logger     *zap.Logger
	config     *config.Config
	repository *repository.Repository
	service    *service.Service
	handler    *handler.Handler
}

func New(configsPath ...string) *App {
	app := new(App)

	app.config = config.New(configsPath...)

	pool, err := pgxpool.New(context.Background(), app.config.Database.Uri)
	if err != nil {
		panic(err)
	}

	app.repository = repository.New(pool)

	playgrounds := service.NewPlaygroundsService(app.config.Playgrounds.CriticalTimeLimit)

	complaints := service.NewComplaintsService(app.repository.Complaints)

	googleDrive, err := storage.NewGoogleDrive(context.Background(), app.config.GoogleDrive.PathToGoogleServiceAccout)
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

	if err := app.initPlaygroundsMap(app.config.Playgrounds.PathToFile); err != nil {
		panic(err)
	}

	if app.config.Server.Debug {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.OutputPaths = append(loggerConfig.OutputPaths, app.config.Logger.Output)
		loggerConfig.ErrorOutputPaths = append(loggerConfig.OutputPaths, app.config.Logger.OutputErrors)

		app.logger, err = loggerConfig.Build()
		if err != nil {
			panic(err)
		}
	} else {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.OutputPaths = append(loggerConfig.OutputPaths, app.config.Logger.Output)
		loggerConfig.ErrorOutputPaths = append(loggerConfig.OutputPaths, app.config.Logger.OutputErrors)

		app.logger, err = loggerConfig.Build()
		if err != nil {
			panic(err)
		}
	}

	app.handler = handler.New(app.service, app.logger)
	app.handler.Init(app.config)

	go app.autoUpdatePlaygroundsMap()

	return app
}

func (a *App) autoUpdatePlaygroundsMap() {
	ticker := time.NewTicker(100 * time.Millisecond)

	defer func() {
		ticker.Stop()
	}()

	for range ticker.C {
		complaints, err := a.service.GetEarly(context.Background())
		if err != nil {
			a.logger.Error(err.Error())
		} else {
			if refreshFunction := a.service.Playgrounds.UpdatePlaygroundsMap(
				a.createPlaygroundsMapFromComplaints(complaints),
			); refreshFunction != nil {
				go refreshFunction()
			}
		}
	}
}

func (a *App) initPlaygroundsMap(pathToResource string) error {
	playgroundsMap := make(map[string]*service.MapProperties)

	file, err := os.OpenFile(pathToResource, os.O_RDONLY, 0777)
	if err != nil {
		return err
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

			return err
		}
		fmt.Printf(`"%s",`, data[0])
		playgroundsMap[data[0]] = &service.MapProperties{
			ID:      data[2],
			Color:   "green",
			Address: data[1],
		}
	}

	a.service.Playgrounds.SetPlaygroundsMap(playgroundsMap)

	if err := a.initFoldersIdsMap(context.Background()); err != nil {
		return err
	}

	complaints, err := a.service.Complaints.GetEarly(context.Background())
	if err != nil {
		return err
	}

	if refreshFunction := a.service.Playgrounds.UpdatePlaygroundsMap(
		a.createPlaygroundsMapFromComplaints(complaints),
	); refreshFunction != nil {
		go refreshFunction()
	}

	return nil
}

func (a *App) createPlaygroundsMapFromComplaints(complaints *[]domain.Complaint) map[string]*service.MapProperties {
	playgroundsMap := make(map[string]*service.MapProperties)

	for _, complaint := range *complaints {
		playgroundsMap[complaint.Coordinates] = &service.MapProperties{
			Time: complaint.CreatedAt,
		}
	}

	return playgroundsMap
}

func (a *App) initFoldersIdsMap(ctx context.Context) error {
	foldersIdsMap := make(map[string]string)

	files, err := a.service.Storage.GetByParams(context.Background(), storage.GoogleDriveParameters{
		ParentId: a.config.GoogleDrive.MediaFolderId,
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
	a.handler.Run(":" + strconv.Itoa(a.config.Server.Port))
}
