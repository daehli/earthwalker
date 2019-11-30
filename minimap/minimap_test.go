package minimap

import (
	"github.com/golang/geo/s2"
	"testing"
)

func TestRender(t *testing.T) {
	err := render(s2.LatLngFromDegrees(
		51.408731,
		19.744071,
	))
	if err != nil {
		t.Error(err)
	}
}
