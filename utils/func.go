package utils

import "fmt"

func StringToUint(s string) uint {
	var id uint
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return 0
	}
	return id
}
