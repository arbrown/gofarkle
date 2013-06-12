package farkle

type GamePlayer interface {
	GamePlay(players []FarkleDecider) (winner_id, turns int)
}