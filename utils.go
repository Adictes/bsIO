package main

// FindEnemy finds enemy for player
func FindEnemy(m map[string]string, s string) string {
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
