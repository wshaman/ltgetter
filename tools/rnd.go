package tools

import (
	"math/rand"
	"time"
)

func RndString(l int) string {
	sym := "abcdefghijklmnopqrstuvwxyz0123456789"
	maxLen := len(sym)
	s := ""
	for i := 0; i < l; i++ {
		r := rand.Intn(maxLen)
		s += string(sym[r])
	}
	return s
}

func RndTime(mid, deviation time.Duration) time.Duration {
	d := int64(deviation) - rand.Int63n(2*int64(deviation))
	return mid + time.Duration(d)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
