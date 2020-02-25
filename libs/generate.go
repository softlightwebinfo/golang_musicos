package libs

import (
	"crypto/rand"
	"fmt"
)

func GenerateName(s string) string {
	n := 20
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	result := fmt.Sprintf("%X", b)
	return fmt.Sprintf("%s.%s", result, s)
}
