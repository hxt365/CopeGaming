package utils

import (
	"fmt"
	"os"
	"strconv"
)

func MustEnv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic(fmt.Sprintf("Missing env %s", name))
	}

	return env
}

func MustStrToFloat32(val string) float32 {
	fVal, err := strconv.ParseFloat(val, 32)
	if err != nil {
		panic(fmt.Sprintf("Couldn't convert str to float32: %s", err))
	}

	return float32(fVal)
}

func MustStrToInt(val string) int {
	iVal, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Couldn't convert str to int: %s", err))
	}

	return int(iVal)
}
