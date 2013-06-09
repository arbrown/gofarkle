package farkle

type GameState struct {
	Players []FarkleDecider
	PlayerNames []string
	PlayerScores []int
	PlayerFarkles []int	
}
