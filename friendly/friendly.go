package friendly

import (
	"fmt"
	"math"
)

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// taken from https://github.com/dustin/go-humanize/blob/v1.0.0/bytes.go#L68
func Bytes(n int64) string {
	if n < 10 {
		return fmt.Sprintf("%d B", n)
	}
	e := math.Floor(math.Log(float64(n)) / math.Log(1000))
	unit := byteUnits[int(e)]
	val := math.Floor(float64(n)/math.Pow(1000, e)*10+0.5) / 10
	format := "%.0f %s"
	if val < 10 {
		format = "%.1f %s"
	}

	return fmt.Sprintf(format, val, unit)
}
