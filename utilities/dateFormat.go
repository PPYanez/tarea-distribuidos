package utilities

import "strings"

func DateFormat(date string) string {
	separated := strings.Split(date, "/")
	separated[2] = "20" + separated[2]
	formatted := strings.Join(reverseSlice(separated), "-")

	return formatted
}

func reverseSlice(slice []string) []string {
	newSlice := make([]string, len(slice))

	for i, v := range slice {
		// add element to new slice
		newSlice[len(slice)-1-i] = v
	}

	return newSlice
}
