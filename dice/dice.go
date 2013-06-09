package dice

import (
	"math/rand"
	"time"
)

func RollDice(numDice int) (dice []int)	 {
	dice = make([]int, numDice)
	for i:=0;i<numDice;i++ {
		dice[i] = rand.Intn(6) + 1
	}
	return dice
}

func Randomize() {
	blowOnDice()
}

// Blowing on dice both randomizes them AND
// provides good luck
func blowOnDice() {
	rand.Seed(time.Now().UTC().UnixNano())
}