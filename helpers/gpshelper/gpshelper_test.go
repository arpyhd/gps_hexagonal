package gpshelper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestProcessLiteral(t *testing.T){

	formated_msg := map[string]string{
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
	msg:="*HQ,7018495625,V6,165900,A,4545.5080,N,02254.3814,E,000.00,000,020323,FFFFFBFF,226,05,703,26034200,8940051807190947031F#"
	returned_map:=ProcessLiteral(msg,formated_msg)
	
	response:=make(map[string]string)
	response["alarms"]=""
	response["battery"]="" 
	response["cell_id"]="703"
	response["cmd"]="V6"
	response["data_valid"]="A"
	response["device_id"]="7018495625"
	response["dir_lat"]="N"
	response["dir_lon"]="E"
	response["direction"]="000"
	response["gps_date"]="2023-03-02 18:59:00"
	response["imei"]="8940051807190947031F"
	response["lat"]="45.758466667"
	response["lon"]="22.906356667"
	response["mcc"]="226"
	response["mileage"]="0" 
	response["mnc"]="05"
	response["speed"]="0"
	response["vehicle_status"]="FFFFFBFF"
	assert.Equal(t,response,returned_map)
}

func TestProcessHex(t *testing.T){
	formated_msg := map[string]string{
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
	msg:="2470184956251658080203234545507804022543814e000000fbfffbff00006a6e0000000000e20502bf000001"
	returned_map:=ProcessHex(msg,formated_msg)
	response:=make(map[string]string)
	response["alarms"]=""
	response["battery"]="" 
	response["cell_id"]="703"
	response["cmd"]="" 
	response["data_valid"]="V"
	response["device_id"]="7018495625"
	response["dir_lat"]="N"
	response["dir_lon"]="E"
	response["direction"]="0"
	response["gps_date"]="2023-03-02 18:58:08"
	response["imei"]="" 
	response["lat"]="45.758463333" 
	response["lon"]="22.906356667"
	response["mcc"]="226"
	response["mileage"]="27246"
	response["mnc"]="5"
	response["speed"]="0"
	response["vehicle_status"]="fbfffbff"
	var responses []map[string]string
	responses=append(responses,response)
	assert.Equal(t,responses,returned_map)
}
