package core

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
)

const upperhex = "0123456789ABCDEF"

func escape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		b.WriteString("%")
		b.WriteByte(upperhex[s[i]>>4])
		b.WriteByte(upperhex[s[i]&15])
	}
	return b.String()
}

func RandomHex(n int, suffix []byte) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(append(bytes, suffix...)), nil
}

func NewLine() string {
	var PS = fmt.Sprintf("%v", os.PathSeparator)
	var LineBreak = "\n"
	if PS != "/" {
		LineBreak = "\r\n"
	}
	return LineBreak
}
