package minimap

import (
	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

func render(position s2.LatLng) error {
	context := sm.NewContext()
	context.SetSize(400, 400)
	// context.SetCenter(position)
	context.SetBoundingBox(s2.RectFromLatLng(position))
	img, err := context.Render()
	if err != nil {
		return err
	}

	err = gg.SavePNG("test.png", img)
	if err != nil {
		return err
	}

	return nil
}
