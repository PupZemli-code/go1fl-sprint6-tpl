package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Service(msg string) (string, error) {
	var morseText bool
	if msg == "" {
		return "", errors.New("файл пуст")
	}

	if strings.Contains(string(msg), "..") ||
		strings.Contains(string(msg), "--") ||
		strings.Contains(string(msg), "-.") ||
		strings.Contains(string(msg), ".-") {
		morseText = true
	}

	if morseText {
		out := morse.ToText(msg)
		return out, nil
	} else {
		out := morse.ToMorse(msg)
		return out, nil
	}
}
