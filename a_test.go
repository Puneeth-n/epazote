package epazote

import (
	"bytes"
	"log"
)

// for catching the log.Println
var buf *bytes.Buffer

func init() {
	buf = new(bytes.Buffer)
	log.SetOutput(buf)
}
