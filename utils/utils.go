package utils

import (
	"encoding/json"
	"fmt"
)

func LogAsJSON(tag string, body any) {
	jsonData, _ := json.MarshalIndent(body, "", "  ")
	fmt.Printf("[%s]:\n", tag)
	fmt.Println(string(jsonData))
	fmt.Println()
}

func BoolToNumber(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func AndNumbers(n1, n2 float64) float64 {
	if n1 > 0 && n2 > 0 {
		return 1
	}
	return 0
}

func OrNumbers(n1, n2 float64) float64 {
	if n1 > 0 || n2 > 0 {
		return 1
	}
	return 0
}
