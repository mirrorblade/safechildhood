package service

import (
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

func (p *PlaygroundsService) GetPlaygrounds() *geojson.FeatureCollection {
	return p.playgrounds
}

func (p *PlaygroundsService) CheckRefreshState() bool {
	return p.refreshState
}

func (p *PlaygroundsService) SetPlaygroundsMap(playgroundsMap map[string]*MapProperties) {
	p.playgroundsMap = make(map[string]*MapProperties)

	for k, v := range playgroundsMap {
		p.playgroundsMap[k] = v
	}

	p.updatePlaygrounds()
}

func (p *PlaygroundsService) UpdatePlaygroundsMap(playgroundsMap map[string]*MapProperties) {
	updateBool := false

	for k, v := range p.playgroundsMap {
		properties, ok := playgroundsMap[k]
		if !ok || properties.Time == (time.Time{}) {
			if v.Color != "green" {
				v.Color = "green"

				updateBool = true
			}

			continue
		}

		color := ""

		if time.Since(properties.Time) > p.changeColorTime {
			color = "red"
		} else {
			color = "yellow"
		}

		if v.Color != color {
			v.Color = color

			updateBool = true
		}
	}

	if updateBool {
		p.updatePlaygrounds()

		p.refreshState = true

		time.Sleep(2 * time.Second)

		p.refreshState = false
	}
}

func (p *PlaygroundsService) updatePlaygrounds() {
	p.playgrounds = geojson.NewFeatureCollection()

	for stringCoordinates, props := range p.playgroundsMap {
		coordinates, _ := tools.StringToCoordinates(stringCoordinates)

		feature := geojson.NewPointFeature([]float64{coordinates.Longitude, coordinates.Latitude})

		feature.Properties = map[string]any{
			"id":      props.ID,
			"color":   props.Color,
			"address": props.Address,
		}

		p.playgrounds.AddFeature(feature)
	}
}
