package friendly

import (
	"fmt"
	"math"
)

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// based on https://github.com/dustin/go-humanize/blob/v1.0.0/bytes.go
// but with KB instead of kB and we hide .0
func Bytes(n int64) string {
	if n < 10 {
		return fmt.Sprintf("%d B", n)
	}
	e := math.Floor(math.Log(float64(n)) / math.Log(1000))
	unit := byteUnits[int(e)]
	val := math.Floor(float64(n)/math.Pow(1000, e)*10+0.5) / 10
	format := "%.0f %s"
	if val < 10 && val != math.Trunc(val) {
		format = "%.1f %s"
	}

	return fmt.Sprintf(format, val, unit)
}
