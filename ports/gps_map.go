package ports

type GpsMap interface {
	Add(int, GpsPersistent) error
	Del(int) error
	Get(int) (Gps, error)
}
