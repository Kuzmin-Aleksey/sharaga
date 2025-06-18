package ui

import (
	"fmt"
	"strconv"
)

func formatPrice(price int) string {
	rub := price / 100
	kop := price % 100

	kopStr := strconv.Itoa(kop)
	if len(kopStr) == 1 {
		kopStr = "0" + kopStr
	}

	return fmt.Sprintf("%d,%s", rub, kopStr)

}
