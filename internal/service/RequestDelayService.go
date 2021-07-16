package service

import "math/rand"

func GetMaxRequestDelay(maxExecutionTime uint16, pageCount uint16) float32 {
	return float32(maxExecutionTime) / float32(pageCount)
}

func GetRandomDelay(maxRequestDelay float32) float32 {
	minRequestDelay := maxRequestDelay / 2
	return minRequestDelay + rand.Float32() * (maxRequestDelay - minRequestDelay)
}
