package skits

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

// GenerateSkitHash returns a hash
func GenerateSkitHash(s string) (hash string) {
	time := time.Now().String()
	hasher := md5.New()
	io.WriteString(hasher, s+time)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
