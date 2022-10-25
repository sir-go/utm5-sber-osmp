package main

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"unicode"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const DefaultConfFile = "config.toml"

func RightPadID(prefix int, id int, totalLen int) (int64, error) {
	idStr := strconv.Itoa(id)
	prefixStr := strconv.Itoa(prefix)
	gapLen := totalLen - len(prefixStr) - len(idStr)
	if gapLen < 0 {
		return 0, errors.New("can't pad id - result is too big: " + prefixStr + idStr)
	}
	gapFormat := fmt.Sprintf("%%0%dd", gapLen)
	resStr := prefixStr + fmt.Sprintf(gapFormat, 0) + idStr
	return strconv.ParseInt(resStr, 10, 64)
}

func HasDigitsOnly(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func RoundBalance(b float64) int {
	// -3999.0212813620074 -> -4000
	// -3999.36895 -> -4000
	// -3999.01234 -> -4000
	// 3999.36895 -> 3999
	// 3999.01234 -> 3999
	// -0.00021 -> -1
	// 0.00021 -> 0

	if b < 0 {
		return int(math.Floor(b))
	}
	return int(math.Trunc(b))
}

func RoundRecSum(b float64) int {
	// -3999.0212813620074 -> 4000
	// -3999.36895 -> 4000
	// -3999.01234 -> 4000
	// 3999.36895 -> 3999
	// 3999.01234 -> 3999
	// -0.00021 -> 1
	// 0.00021 -> 1

	return int(math.Ceil(math.Abs(b)))
}

func MakeUUID() string {
	strUUID := uuid.NewString()
	if id, err := uuid.NewRandom(); err != nil {
		return strUUID
	} else {
		if binID, err := id.MarshalBinary(); err != nil {
			return strUUID
		} else {
			return hex.EncodeToString(binID)
		}
	}
}
