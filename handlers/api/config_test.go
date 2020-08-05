package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/glatteis/earthwalker/domain"
)

func TestGetTileServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/config/tileserver", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	conf := domain.Config{TileServerURL: "https://tiles.wmflabs.org/osm/{z}/{x}/{y}.png"}
	handler := Root{ConfigHandler: Config{Config: conf}}

	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("got status code %v, expected %v", status, http.StatusOK)
	}
	expected := fmt.Sprintf(`{"tileserver": "%s"}`, conf.TileServerURL)
	if recorder.Body.String() != expected {
		t.Errorf("got response body %v, expected %v", recorder.Body.String(), expected)
	}
}

func TestGetNoLabelTileServer(t *testing.T) {
	req, err := http.NewRequest("GET", "/config/nolabeltileserver", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	conf := domain.Config{NoLabelTileServerURL: "https://tiles.wmflabs.org/osm-no-labels/{z}/{x}/{y}.png"}
	handler := Root{ConfigHandler: Config{Config: conf}}

	handler.ServeHTTP(recorder, req)
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("got status code %v, expected %v", status, http.StatusOK)
	}
	expected := fmt.Sprintf(`{"tileserver": "%s"}`, conf.NoLabelTileServerURL)
	if recorder.Body.String() != expected {
		t.Errorf("got response body %v, expected %v", recorder.Body.String(), expected)
	}
}
