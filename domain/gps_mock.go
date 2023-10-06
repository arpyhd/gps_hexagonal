package domain

import (
	"gps_hexagonal/ports"
)

type GpsMock struct{}

func (t GpsMock) Read(buf []byte) error { return nil }
func (t GpsMock) Send() error           { return nil }

func NewGpsMock() ports.Gps {
	return GpsMock{}
}
