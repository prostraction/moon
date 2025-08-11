package moon

import (
	"fmt"
	"testing"
)

func TestCalcMoonNumber(t *testing.T) {
	for i := 1700; i < 2200; i++ {
		fmt.Println(i, CalcMoonNumber(i))
	}
}
