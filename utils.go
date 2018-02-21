package main

import (
	"strings"
)

// GetEnemy finds enemy for player
func GetEnemy(m map[string]string, s string) string {
	for k, v := range m {
		if k == s {
			return v
		}
		if v == s {
			return k
		}
	}
	return ""
}

// HaveAvailableGame checks has 'm' 's'. If not return ""
func HaveAvailableGame(m map[string]string, s string) string {
	for k, v := range m {
		if k == s {
			continue
		}
		if k != "" && v == "" {
			return k
		}
	}
	return ""
}

// ChangeLetter changes letter - 'e' to 'h' in all struct's fields
func ChangeLetter(s *StrickenShips) {
	for i, v := range s.Ambient {
		s.Ambient[i] = strings.Replace(v, "e", "h", 1)
	}
	if s.Hitted != "" {
		s.Hitted = strings.Replace(s.Hitted, "e", "h", 1)
	}
}
