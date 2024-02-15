package helper

import (
	"fmt"
	"regexp"
	"strings"
)

func FormatDateString(dateStr string) string {
	parts := strings.Split(dateStr, " ")
	dateParts := strings.Split(parts[0], "-")
	return fmt.Sprintf("%s/%s/%s", dateParts[2], dateParts[1], dateParts[0])
}

func IsValidDateFormat(dateStr string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`)
	return re.MatchString(dateStr)
}
