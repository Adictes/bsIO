package main

import (
	"crypto/rand"
	"math/big"
)

// GetRandomEnemy returns enemy username
func GetRandomEnemy(currentUser string) string {
	for {
		i, _ := rand.Int(rand.Reader, big.NewInt(int64(len(readyToPlay))))
		if currentUser == readyToPlay[i.Int64()] {
			continue
		} else {
			return readyToPlay[i.Int64()]
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
