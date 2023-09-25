package service

import (
	"fmt"
	"safechildhood/internal/app/domain"
	"safechildhood/tools"
	"time"

	geojson "github.com/paulmach/go.geojson"
)

type MapProperties struct {
	ID      string    `json:"id"`
	Color   string    `json:"color"`
	Address string    `json:"address"`
	Time    time.Time `json:"time"`
}

type PlaygroundsService struct {
	playgrounds    *geojson.FeatureCollection
	playgroundsMap map[string]*MapProperties

	changeColorTime time.Duration
	refreshState    bool
}

func NewPlaygroundsService(changeColorTime time.Duration) *PlaygroundsService {
	return &PlaygroundsService{
		changeColorTime: changeColorTime,
		playgrounds:     &geojson.FeatureCollection{},
	}
}

func (p *PlaygroundsService) SetPlaygroundsMap(playgroundsMap map[string]*MapProperties) {
	p.playgroundsMap = make(map[string]*MapProperties)

	for k, v := range playgroundsMap {
		p.playgroundsMap[k] = v
	}

	fmt.Println(playgroundsMap)
}

func (p *PlaygroundsService) GetPlaygrounds() *geojson.FeatureCollection {
	return p.playgrounds
}

func (p *PlaygroundsService) CheckRefreshState() bool {
	return p.refreshState
}

func (p *PlaygroundsService) RefreshPlaygrounds() {
	defer func(p *PlaygroundsService) {
		time.Sleep(3 * time.Second)

		p.refreshState = false
	}(p)

	for stringCoordinates, props := range p.playgroundsMap {
		coordinates, _ := tools.StringToCoordinates(stringCoordinates)

		feature := geojson.NewPointFeature([]float64{coordinates.Longitude, coordinates.Latitude})

		feature.Properties = map[string]any{
			"id":      props.ID,
			"color":   props.Color,
			"address": props.Address,
		}

		if props.Time != (time.Time{}) {
			feature.Properties["time"] = props.Time
		}

		p.playgrounds.AddFeature(feature)
	}

	fmt.Println(p.playgrounds)

	p.refreshState = true
}

func (p *PlaygroundsService) UpdatePlaygroundsMap(complaints []domain.Complaint) {
	for _, complaint := range complaints {
		if playground, ok := p.playgroundsMap[complaint.Coordinates]; ok {
			playground.Time = complaint.CreatedAt
			playground.Color = "yellow"

			if time.Since(complaint.CreatedAt) > p.changeColorTime {
				playground.Color = "red"
			}
		}
	}

	p.RefreshPlaygrounds()
}

func (p *PlaygroundsService) AutoCheckerFeatureTime() {
	ticker := time.NewTicker(1 * time.Second)

	update := false

	for range ticker.C {
		for _, props := range p.playgroundsMap {
			if props.Color != "yellow" {
				continue
			}

			if time.Since(props.Time) > p.changeColorTime {
				props.Color = "red"

				update = true
			}
		}

		if update {
			p.RefreshPlaygrounds()

			update = false
		}
	}
}
