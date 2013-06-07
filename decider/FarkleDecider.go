package decider

type FarkleDecider interface {
	FarkleDecide(dice []int, runScore int, totalScore int, numFarkles int, opponenetScores []int) (keep []bool, rollAgain bool)
}