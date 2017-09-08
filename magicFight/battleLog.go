package main

type BattleLog struct {
	BattleID string
	Logs     map[int]string // logs[turn]result
}
