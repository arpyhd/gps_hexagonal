package ports

type GpsPersistent interface {
	GetId(device_id string) (int, string, error)
	SetImei(gps_id int, imei string) error
	InsertEvent(gps_id int, message map[string]string) error
	ReadCommand(gps_id int) (string, error)
}
