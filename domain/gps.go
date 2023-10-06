package domain

import (
	"encoding/hex"
	"errors"
	"gps_hexagonal/helpers/gpshelper"
	"gps_hexagonal/helpers/logger"
	"gps_hexagonal/ports"
	"strings"
)

type Gps struct {
	gps_id    int
	status    int
	imei      string
	device_id string
	gs        ports.GpsPersistent
}

func NewGps(gs ports.GpsPersistent) *Gps {
	g := &Gps{
		gps_id:    0,
		status:    0,
		imei:      "",
		device_id: "",
		gs:        gs,
	}

	return g
}

func (g *Gps) Read(buf []byte) error {

	str_buf := string(buf)
	logger.Log.Println("Raw msg", str_buf)
	if len(str_buf) > 0 {
		message := map[string]string{
			"device_id":      "",
			"cmd":            "",
			"gps_date":       "",
			"data_valid":     "",
			"lat":            "",
			"dir_lat":        "",
			"lon":            "",
			"dir_lon":        "",
			"speed":          "",
			"direction":      "",
			"vehicle_status": "",
			"mcc":            "",
			"mnc":            "",
			"cell_id":        "",
			"battery":        "",
			"mileage":        "0",
			"imei":           "",
			"alarms":         "",
		}
		slice_str_buf := strings.Split(str_buf, ",")
		if str_buf[0] == '*' {
			if g.gps_id == 0 {
				gps_id, imei, err := g.gs.GetId(slice_str_buf[1])
				if err == nil {
					g.gps_id = gps_id
					g.imei = imei
					g.status = 1
				} else {
					return err
				}
			}
			if slice_str_buf[2] == "V6" {

				message = gpshelper.ProcessLiteral(str_buf, message)
				message["imei"] = message["imei"][:15]
				if g.imei == "" && g.gps_id > 0 {
					g.gs.SetImei(g.gps_id, message["imei"])
					g.imei = message["imei"]
					g.status = 2
				}

			} else if slice_str_buf[2] == "V5" {
				gpshelper.ProcessLiteral(str_buf, message)
			} else {
				return nil
			}
			message = gpshelper.VehicleStatus(message)
			err := g.gs.InsertEvent(g.gps_id, message)
			if err != nil {
				logger.Log.Println("Gps Insert Event error : ", err)
			}

		} else if str_buf[0] == '$' {
			hex_string := hex.EncodeToString(buf)
			messages := gpshelper.ProcessHex(hex_string, message)
			for _, message := range messages {
				if g.gps_id == 0 {
					gps_id, imei, err := g.gs.GetId(message["device_id"])
					if err == nil {
						g.gps_id = gps_id
						g.imei = imei
						g.status = 1
					} else {
						return err
					}
				}
				message = gpshelper.VehicleStatus(message)
				err := g.gs.InsertEvent(g.gps_id, message)
				if err != nil {
					logger.Log.Println("Gps Insert Event Hex error : ", err)
				}
			}
		} else if str_buf[0] == '%' {
			hex_string := hex.EncodeToString(buf)
			logger.Log.Println("Raw AMR", hex_string)

		}

	} else {
		if g.status == 0 {
			return nil
		}
		logger.Log.Println(hex.EncodeToString(buf))
	}

	return nil
}

func (g *Gps) Send() error {

	if g.gps_id > 0 && g.imei != "" {
		g.gs.ReadCommand(g.gps_id)
		return nil
	} else {
		err := errors.New("Gps id not found")
		return err
	}
}
