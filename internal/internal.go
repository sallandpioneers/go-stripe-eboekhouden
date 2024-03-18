package internal

import (
	"bytes"
	"crypto/rand"
	"fmt"
)

type ServicesAvailable struct {
	DB bool
	FS bool
}

func Concat(values ...string) string {
	var buffer bytes.Buffer
	for _, s := range values {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func ConcatBytes(values ...[]byte) []byte {
	var buffer bytes.Buffer
	for _, s := range values {
		buffer.Write(s)
	}
	return buffer.Bytes()
}

func ModeUnknown(mode, givenMode string, possibleModes ...string) error {
	var data string
	for _, v := range possibleModes {
		data = Concat(data, " ['", v, "']")
	}
	data = Concat(data, " ")
	return fmt.Errorf("%s mode unknown. \t possible modes: %s\t given mode: %s", mode, data, givenMode)
}

func GetRandomChars(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytess, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytess {
		bytess[i] = letters[b%byte(len(letters))]
	}
	return string(bytess), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
