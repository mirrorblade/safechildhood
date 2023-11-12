package service

import (
	"safechildhood/pkg/converter"
	"time"

	geojson "github.com/paulmach/go.geojson"
)

type RefreshFunction func()

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
	refreshChannel  chan any
}

func NewPlaygroundsService(changeColorTime time.Duration) *PlaygroundsService {
	return &PlaygroundsService{
		changeColorTime: changeColorTime,
		playgrounds:     geojson.NewFeatureCollection(),
		refreshChannel:  make(chan any),
	}
}

func (p *PlaygroundsService) GetPlaygrounds() *geojson.FeatureCollection {
	return p.playgrounds
}

func (p *PlaygroundsService) Refresh() chan any {
	return p.refreshChannel
}

func (p *PlaygroundsService) SetPlaygroundsMap(playgroundsMap map[string]*MapProperties) RefreshFunction {
	newPlaygroundsMap := make(map[string]*MapProperties)

	for k, v := range playgroundsMap {
		newPlaygroundsMap[k] = v
	}

	p.playgroundsMap = newPlaygroundsMap

	p.updatePlaygrounds()

	return func() {
		p.refreshChannel <- struct{}{}
	}
}

func (p *PlaygroundsService) UpdatePlaygroundsMap(playgroundsMap map[string]*MapProperties) RefreshFunction {
	update := false

	for k, v := range p.playgroundsMap {
		properties, ok := playgroundsMap[k]
		if !ok || properties.Time == (time.Time{}) {
			if v.Color != "green" {
				v.Color = "green"

				update = true
			}

			continue
		}

		color := "yellow"

		if time.Since(properties.Time) > p.changeColorTime {
			color = "red"
		}

		if v.Color != color {
			v.Color = color

			update = true
		}
	}

	if update {
		p.updatePlaygrounds()

		return func() {
			p.refreshChannel <- struct{}{}
		}
	}

	return nil
}

func (p *PlaygroundsService) updatePlaygrounds() {
	playgrounds := geojson.NewFeatureCollection()

	for stringCoordinates, props := range p.playgroundsMap {
		coordinates, _ := converter.StringToCoordinates(stringCoordinates)

		feature := geojson.NewPointFeature([]float64{coordinates.Longitude, coordinates.Latitude})

		feature.Properties = map[string]any{
			"id":      props.ID,
			"color":   props.Color,
			"address": props.Address,
		}

		playgrounds.AddFeature(feature)
	}

	p.playgrounds = playgrounds
}
