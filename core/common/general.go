package common

import (
	"math/rand"
	"time"
)

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandomValue(list *[]string) string {
	return (*list)[rand.Intn(len(*list))]
}
