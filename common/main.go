package common

import "fmt"

func FormatByte(B uint64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case B >= GB:
		return fmt.Sprintf("%.2f GB", float64(B)/float64(GB))
	case B >= MB:
		return fmt.Sprintf("%.2f MB", float64(B)/float64(MB))
	case B >= KB:
		return fmt.Sprintf("%.2f KB", float64(B)/float64(KB))
	default:
		return fmt.Sprintf("%d B", B)
	}
}
