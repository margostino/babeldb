package common

import "log"

func Check(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err.Error())
	}
}

func GetOrDefault(index int, list []string) string {
	if len(list) >= index+1 {
		return list[index]
	}
	return ""
}

func SilentCheck(err error, message string) {
	if err != nil {
		log.Printf("Error: %s - %s\n", err.Error(), message)
	}
}

func IsError(err error, message string) bool {
	if err != nil {
		log.Printf("Error: %s - %s\n", err.Error(), message)
		return true
	}
	return false
}
