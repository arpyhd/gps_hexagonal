package domain

import (
	"errors"
	logger "gps_hexagonal/helpers/logger"
	"gps_hexagonal/ports"
)

type GpsMap struct {
	gmap map[int]*Gps
}

func NewGpsMap() *GpsMap {
	g := &GpsMap{gmap: make(map[int]*Gps)}
	return g
}

func (gm *GpsMap) Add(id int, gp ports.GpsPersistent) error {
	if _, ok := gm.gmap[id]; ok {
		return errors.New("Gps exist")
	}
	logger.Log.Println("Gps_map Add: ", id)
	gm.gmap[id] = NewGps(gp)
	return nil
}

func (gm *GpsMap) Del(id int) error {
	if _, ok := gm.gmap[id]; ok {
		delete(gm.gmap, id)
		return nil
	}
	return errors.New("Key not found")
}
func (gm *GpsMap) Get(id int) (ports.Gps, error) {
	if val, ok := gm.gmap[id]; ok {
		return val, nil
	} else {
		error := errors.New("Gps id not found")
		return nil, error
	}
}
