package service

import (
	"fmt"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

// RoundBalance
// -3999.0212813620074 -> -4000
// -3999.36895 -> -4000
// -3999.01234 -> -4000
// 3999.36895 -> 3999
// 3999.01234 -> 3999
// -0.00021 -> -1
// 0.00021 -> 0
func RoundBalance(b float64) int {
	if b < 0 {
		return int(math.Floor(b))
	}
	return int(math.Trunc(b))
}

// RoundRecSum
// -3999.0212813620074 -> 4000
// -3999.36895 -> 4000
// -3999.01234 -> 4000
// 3999.36895 -> 4000
// 3999.01234 -> 4000
// -0.00021 -> 1
// 0.00021 -> 1
func RoundRecSum(b float64) int {
	return int(math.Ceil(math.Abs(b)))
}

func RightPadID(prefix int, id int, totalLen int) (int64, error) {
	idStr := strconv.Itoa(id)
	prefixStr := strconv.Itoa(prefix)
	gapLen := totalLen - len(prefixStr) - len(idStr)
	if gapLen < 1 {
		return 0, errors.New("can't pad id - result is too big: " + prefixStr + idStr)
	}
	gapFormat := fmt.Sprintf("%%0%dd", gapLen)
	resStr := prefixStr + fmt.Sprintf(gapFormat, 0) + idStr
	return strconv.ParseInt(resStr, 10, 64)
}
