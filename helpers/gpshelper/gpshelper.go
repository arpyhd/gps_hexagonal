package gpshelper

import (
	"fmt"
	"gps_hexagonal/helpers/logger"
	"math"
	"strconv"
	"strings"
	"time"
)

var Alarm = map[string]string{
	"31": "Illegal door open alarm",
	"30": "SOS",
	"29": "Speed alarm",
	"28": "Illegal ignition alarm",
	"27": "Entering alarm",
	"26": "Antena GPS diconnect alarm",
	"25": "Antena GPS diconnect alarm",
	"24": "Out alarm",
	"23": "Door open",
	"22": "Armed",
	"21": "Acc OFF",
	"20": "Crash alarm",
	"19": "Keep",
	"18": "Pump",
	"17": "Custim alarm",
	"16": "Over speed",
	"15": "Mistake GPS alarm",
	"14": "Shock alarm",
	"13": "Tilt alarm",
	"12": "Backup user battery",
	"11": "Battery remove alarm",
	"10": "Antena gps disconnect",
	"9":  "Antena gps short circuit",
	"8":  "Low level sensor 2",
	"7":  "Temple alarm",
	"6":  "Move alaram",
	"5":  "Blind recover alarm",
	"4":  "Oil cut off",
	"3":  "Battery demolition",
	"2":  "Home SOS alarm",
	"1":  "Office SOS alarm",
	"0":  "Low level sensor 1",
}

func ProcessLiteral(msg string, formated_msg map[string]string) map[string]string {
	slice_str_buf := strings.Split(msg, ",")

	if slice_str_buf[2] == "V6" {
		imei := strings.TrimSuffix(slice_str_buf[len(slice_str_buf)-1], "#")
		formated_msg["imei"] = imei
	}
	formated_msg["device_id"] = slice_str_buf[1]
	formated_msg["cmd"] = slice_str_buf[2]
	string_date := "20" + slice_str_buf[11][4:6] + "-" + slice_str_buf[11][2:4] + "-" + slice_str_buf[11][0:2] + "T" + slice_str_buf[3][0:2] + ":" + slice_str_buf[3][2:4] + ":" + slice_str_buf[3][4:6] + "+00:00"

	loc, _ := time.LoadLocation("")
	t, _ := time.ParseInLocation(time.RFC3339, string_date, loc)
	location, _ := time.LoadLocation("EET")
	t = t.In(location)
	gps_date := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	formated_msg["gps_date"] = gps_date
	formated_msg["data_valid"] = slice_str_buf[4]
	formated_msg["lat"] = getLat(slice_str_buf[5])
	formated_msg["dir_lat"] = slice_str_buf[6]
	formated_msg["lon"] = getLon(slice_str_buf[7])
	formated_msg["dir_lon"] = slice_str_buf[8]
	formated_msg["speed"] = getSpeed(slice_str_buf[9])
	formated_msg["direction"] = slice_str_buf[10]
	formated_msg["vehicle_status"] = slice_str_buf[12]
	formated_msg["mcc"] = slice_str_buf[13]
	formated_msg["mnc"] = slice_str_buf[14]
	formated_msg["cell_id"] = slice_str_buf[15]
	logger.Log.Println("Process Literal", formated_msg)
	return formated_msg
}

func ProcessHex(msg string, formated_msg map[string]string) []map[string]string {
	var response []map[string]string
	messages := strings.Split(msg, msg[:12])
	for i, message := range messages {
		if i > 0 {
			response = append(response, processBaseHex(msg[:12]+message, formated_msg))
		}
	}
	return response
}

func processBaseHex(msg string, formated_msg map[string]string) map[string]string {
	bin_types, _ := HexToBin(msg[43:44])
	var data_valid = "V"
	var dir_lon = "W"
	var dir_lat = "S"
	if bin_types[13] == '1' {
		dir_lon = "E"
	}
	if bin_types[14] == '1' {
		dir_lat = "N"
	}

	if bin_types[15] == '1' {
		data_valid = "A"
	}

	formated_msg["device_id"] = msg[2:12]
	formated_msg["cmd"] = ""
	string_date := "20" + msg[22:24] + "-" + msg[20:22] + "-" + msg[18:20] + "T" + msg[12:14] + ":" + msg[14:16] + ":" + msg[16:18] + "+00:00"
	loc, _ := time.LoadLocation("")
	t, _ := time.ParseInLocation(time.RFC3339, string_date, loc)
	location, _ := time.LoadLocation("EET")
	t = t.In(location)
	gps_date := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	formated_msg["gps_date"] = gps_date
	formated_msg["data_valid"] = data_valid
	formated_msg["lat"] = getLatHex(msg[24:32])
	formated_msg["dir_lat"] = dir_lat
	formated_msg["lon"] = getLonHex(msg[34:43])
	formated_msg["dir_lon"] = dir_lon
	formated_msg["speed"] = getSpeed(msg[44:47])
	formated_msg["direction"] = getAngle(msg[47:50])
	formated_msg["vehicle_status"] = msg[50:58]
	value, _ := strconv.ParseInt(msg[74:78], 16, 64)
	formated_msg["mcc"] = strconv.Itoa(int(value))
	value, _ = strconv.ParseInt(msg[78:80], 16, 64)
	formated_msg["mnc"] = strconv.Itoa(int(value))
	value, _ = strconv.ParseInt(msg[80:84], 16, 64)
	formated_msg["cell_id"] = strconv.Itoa(int(value))
	value, _ = strconv.ParseInt(msg[58:66], 16, 64)
	formated_msg["mileage"] = strconv.Itoa(int(value))
	logger.Log.Println("Process Hex", formated_msg)
	return formated_msg
}

func HexToBin(hex string) (string, error) {
	ui, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%016b", ui), nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func getLat(a string) string {
	i, _ := strconv.ParseFloat(a[0:2], 8)
	f, _ := strconv.ParseFloat(a[2:9], 8)
	s := i + f/60.0
	return fmt.Sprintf("%v", roundFloat(s, 9))
}

func getLon(a string) string {
	i, _ := strconv.ParseFloat(a[0:3], 8)
	f, _ := strconv.ParseFloat(a[3:10], 8)
	s := i + f/60.0
	return fmt.Sprintf("%v", roundFloat(s, 9))
}

func getLatHex(a string) string {
	i, _ := strconv.ParseFloat(a[0:2], 8)
	f, _ := strconv.ParseFloat(a[2:8], 8)
	s := i + f/(60.0*10000)
	return fmt.Sprintf("%v", roundFloat(s, 9))
}

func getLonHex(a string) string {
	i, _ := strconv.ParseFloat(a[0:3], 8)
	f, _ := strconv.ParseFloat(a[3:9], 8)
	s := i + f/(60.0*10000)
	return fmt.Sprintf("%v", roundFloat(s, 9))
}

func getSpeed(a string) string {
	var i float64

	if a == "" {
		i = 0
	}

	i, _ = strconv.ParseFloat(a, 8)
	return fmt.Sprintf("%v", roundFloat(i*1.852, 2))
}

func getAngle(a string) string {
	var i float64

	if a == "" {
		i = 0
	}

	i, _ = strconv.ParseFloat(a, 8)
	return fmt.Sprintf("%v", roundFloat(i, 2))
}

func asBits(val uint64) []uint64 {
	bits := []uint64{}
	for i := 0; i < 32; i++ {
		bits = append([]uint64{val & 0x1}, bits...)
		val = val >> 1
	}
	return bits
}

func VehicleStatus(message map[string]string) map[string]string {
	i, err := strconv.ParseUint(message["vehicle_status"], 16, 32)
	if err != nil {
		logger.Log.Println("Error procesing Vehicle ", message["vehicle_status"])
	}
	var response = ""
	events := asBits(i)
	if events[31] == 0 {
		response += Alarm["31"] + ", "
	}
	if events[30] == 0 {
		response += Alarm["30"] + ", "
	}
	if events[29] == 0 {
		response += Alarm["29"] + ", "
	}
	if events[28] == 0 {
		response += Alarm["28"] + ", "
	}
	if events[27] == 0 {
		response += Alarm["27"] + ", "
	}
	if events[26] == 0 {
		response += Alarm["26"] + ", "
	}
	if events[25] == 0 {
		response += Alarm["25"] + ", "
	}
	if events[24] == 0 {
		response += Alarm["24"] + ", "
	}
	if events[23] == 0 {
		response += Alarm["23"] + ", "
	}
	if events[22] == 0 {
		response += Alarm["22"] + ", "
	}
	if events[21] == 0 {
		response += Alarm["21"] + ", "
	}
	if events[20] == 0 {
		response += Alarm["20"] + ", "
	}
	if events[19] == 0 {
		response += Alarm["19"] + ", "
	}
	if events[18] == 0 {
		response += Alarm["18"] + ", "
	}
	if events[17] == 0 {
		response += Alarm["17"] + ", "
	}
	if events[16] == 0 {
		response += Alarm["16"] + ", "
	}
	if events[15] == 0 {
		response += Alarm["15"] + ", "
	}
	if events[14] == 0 {
		response += Alarm["14"] + ", "
	}
	if events[13] == 0 {
		response += Alarm["13"] + ", "
	}
	if events[12] == 0 {
		response += Alarm["12"] + ", "
	}
	if events[11] == 0 {
		response += Alarm["11"] + ", "
	}
	if events[10] == 0 {
		response += Alarm["10"] + ", "
	}
	if events[9] == 0 {
		response += Alarm["9"] + ", "
	}
	if events[8] == 0 {
		response += Alarm["8"] + ", "
	}
	if events[7] == 0 {
		response += Alarm["7"] + ", "
	}
	if events[6] == 0 {
		response += Alarm["6"] + ", "
	}
	if events[5] == 0 {
		response += Alarm["5"] + ", "
	}
	if events[4] == 0 {
		response += Alarm["4"] + ", "
	}
	if events[3] == 0 {
		response += Alarm["3"] + ", "
	}
	if events[2] == 0 {
		response += Alarm["2"] + ", "
	}
	if events[1] == 0 {
		response += Alarm["1"] + ", "
	}
	message["alarms"] = response
	return message
}
