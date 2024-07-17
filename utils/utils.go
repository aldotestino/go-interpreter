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
