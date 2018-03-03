package main

// StrickenShips is used as JSON wrapper for sending it to websocket
type StrickenShips struct {
	Ambient []string
	Hitted  string
}

// TurnWrapper sending with true
// if player moving, false otherwise
type TurnWrapper struct {
	Turn bool
}

// WinWrapper sending with true
// if player wins, false otherwise
type WinWrapper struct {
	Win bool
}

// CorrectnessWrapper sending with true
// if ships set properly, false otherwise
type CorrectnessWrapper struct {
	Correctness bool
}

// NameWrapper sending with name
type NameWrapper struct {
	Name string
}
