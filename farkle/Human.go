package farkle

import (
	"regexp"
	"fmt"
	"strconv"
)

type Human struct {
	PrintDice bool
}


func (h Human) FarkleDecide(dice []int, runScore int, game GameState, player int) (keep []bool, rollAgain bool) {
	var keepstr, rollstr, pad string

	ndice := len(dice)
	keep = make([]bool, ndice)

	pad = "          "

	if h.PrintDice {
		fmt.Println(dice)
		pad = ""
	}

	zerofive := []int{0,1,2,3,4,5}

	fmt.Printf("%s%v\n", pad, zerofive[:ndice])

	fmt.Printf("Which dice to keep? ")

	fmt.Scan(&keepstr)

	re, _ := regexp.Compile(`\d`)

	digits := re.FindAllString(keepstr, ndice)

	for _, s := range digits {
		i, ierr := strconv.Atoi(s)
		if i < ndice && ierr == nil {
			keep[i] = true
		}
	}

	fmt.Printf("Continue rolling? ")
	nc, err := fmt.Scan(&rollstr)

	if err != nil {
		fmt.Printf("Read %d chars\n%s\n", nc, err.Error())
	}

	rollAgain, _ = regexp.MatchString(`(?i)\s*y(es)?\s*`, rollstr)
	return
}
