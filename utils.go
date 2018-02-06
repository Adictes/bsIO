package main

import (
	"math/rand"
	"time"
)

// GetRandomEnemy returns enemy username
func GetRandomEnemy(currentUser string) string {
	rand.Seed(time.Now().Unix())
	for {
		i := rand.Intn(len(readyToPlay))
		if currentUser == readyToPlay[i] {
			continue
		} else {
			return readyToPlay[i]
		}
	}
}

// contain says has 'ss' 's' or not
func contain(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}
