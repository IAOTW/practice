package myjson

import (
	"encoding/json"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
)

type Host struct {
	IP   string
	Name string
}

func Convert2Json(h *Host) (string, error) {
	b, err := json.Marshal(h)
	return string(b), err
}

func TestConvert2Json(t *testing.T) {
	patches := gomonkey.ApplyFunc(json.Marshal, func(v interface{}) ([]byte, error) {
		return []byte(`{"IP":"192.168.23.92","Name":"Sky"}`), nil
	})

	defer patches.Reset()

	h := Host{Name: "Sky", IP: "192.168.23.92"}
	s, err := Convert2Json(&h)

	expectedString := `{"IP":"192.168.23.92","Name":"Sky"}`

	if s != expectedString || err != nil {
		t.Errorf("expected %v, got %v", expectedString, s)
	}

}
