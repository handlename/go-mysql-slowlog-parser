package slowlog

import (
	"log"
	"strconv"
)

func stringToFloat32(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)

	if err != nil {
		log.Println(err)
	}

	return float32(f)
}

func stringToInt32(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)

	if err != nil {
		log.Println(err)
	}

	return int32(i)
}
