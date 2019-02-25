package costexp

import (
	"strconv"
	"strings"
)

/*
GetInstanceHourAndNum returns instance hour and instance num.
*/
func GetInstanceHourAndNum(amount, start string) (float64, float64) {
	hrs, _ := strconv.ParseFloat(amount, 64)
	month := strings.Split(start, "-")[1]
	num := hrs / float64(24*Days[month])
	return hrs, num
}

var Days = map[string]int{
	"01": 31,
	"02": 28,
	"03": 31,
	"04": 30,
	"05": 31,
	"06": 30,
	"07": 31,
	"08": 31,
	"09": 30,
	"10": 31,
	"11": 30,
	"12": 31,
}
