package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGpsMapAdd(t *testing.T) {
	gps_domain := NewGpsMap()
	err := gps_domain.Add(1, NewGpsPersistMock())
	assert.NoError(t, err)
}

func TestGpsMapDel(t *testing.T) {
	gps_domain := NewGpsMap()
	err := gps_domain.Del(2)
	assert.Error(t, err)
}

func TestGpsMapGet(t *testing.T) {
	gps_domain := NewGpsMap()
	_, err := gps_domain.Get(3)
	assert.Error(t, err)
}
