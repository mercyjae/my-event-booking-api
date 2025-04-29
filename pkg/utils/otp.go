package utils

import (
	"strconv"
	"time"
    "math/rand"
	
)

func GenerateOTP() string {
    rand.Seed(time.Now().UnixNano())
    otp := rand.Intn(900000) + 100000 
    return strconv.Itoa(otp)
}