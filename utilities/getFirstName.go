package utilities

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func GetFirstName(fullName string) string {
	parts := strings.Split(fullName, " ")
	if len(parts) > 0 {
		return parts[0]
	}
	return fullName
}

func GetFirstNameWithNumbers(fullName string) string {
	firstName := GetFirstName(fullName)
    source := rand.NewSource(time.Now().UnixNano())
    random := rand.New(source)
    randomNumbers := random.Intn(900) + 100 // menghasilkan angka acak 3 digit antara 100 dan 999
    return firstName + fmt.Sprintf("%03d", randomNumbers)
}