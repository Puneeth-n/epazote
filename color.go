package epazote

import (
	"fmt"
)

const escape = "\x1b"

func Red(s string) string {
	return fmt.Sprintf("%s[0;31m%s%s[0;00m", escape, s, escape)
}

func Green(s string) string {
	return fmt.Sprintf("%s[0;32m%s%s[0;00m", escape, s, escape)
}
