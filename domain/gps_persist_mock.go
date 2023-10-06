package domain

import (
	"gps_hexagonal/ports"
)

type GpsPersistMock struct{}

func (t GpsPersistMock) GetId(device_id string) (int, string, error)             { return 1, "test", nil }
func (t GpsPersistMock) SetImei(gps_id int, imei string) error                   { return nil }
func (t GpsPersistMock) InsertEvent(gps_id int, message map[string]string) error { return nil }
func (t GpsPersistMock) ReadCommand(gps_id int) (string, error)                  { return "", nil }
func NewGpsPersistMock() ports.GpsPersistent {
	return GpsPersistMock{}
}
