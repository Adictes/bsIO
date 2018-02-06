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

// delete removes 's' from 'ss'
func delete(ss []string, s string) {
	i := index(ss, s)
	ss = append(ss[:i], ss[i+1:]...)
}

// index @TODO binary search
func index(ss []string, s string) int {
	for i, v := range ss {
		if v == s {
			return i
		}
	}
	return -1
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
