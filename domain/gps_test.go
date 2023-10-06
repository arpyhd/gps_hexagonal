package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGpsRead(t *testing.T) {
	gps := NewGps(NewGpsPersistMock())
	err := gps.Read(make([]byte, 5))
	assert.NoError(t, err)
}

func TestGpsSend(t *testing.T) {
	gps := NewGps(NewGpsPersistMock())
	err := gps.Send()
	assert.Error(t, err)
}
