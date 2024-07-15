package utilities

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOTP generates a 5-digit OTP
func GenerateOTP() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%05d", rand.Intn(100000))
}

// GenerateExpiryTime generates the expiration time for the OTP
func GenerateExpiryTime() time.Time {
    return time.Now().Add(10 * time.Minute)  // OTP valid for 10 minutes
}