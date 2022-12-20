package code

import (
	"fmt"
	"math/rand"
	"time"
)

func Sms() string {
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(10000)
	return fmt.Sprintf("%04d", num)
}
