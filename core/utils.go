package core

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
)

const upperhex = "0123456789ABCDEF"

func GetEnvOrDefault(key string, defaultValue int) int {
	if envValue, exists := os.LookupEnv(key); exists {
		if envValueInt, err := strconv.Atoi(envValue); err == nil {
			return envValueInt
		}
	}
	return defaultValue
}

func Escape(s string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		b.WriteString("%")
		b.WriteByte(upperhex[s[i]>>4])
		b.WriteByte(upperhex[s[i]&15])
	}
	return b.String()
}

func RandomHex(n int, suffix []byte) (string, error) {
	_bytes := make([]byte, n)
	if _, err := rand.Read(_bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(append(_bytes, suffix...)), nil
}

func NewLine() string {
	var PS = fmt.Sprintf("%v", os.PathSeparator)
	var LineBreak = "\n"
	if PS != "/" {
		LineBreak = "\r\n"
	}
	return LineBreak
}
