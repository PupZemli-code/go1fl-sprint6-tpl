package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

func Service(msg string) (string, error) {

	if msg == "" {
		return "", errors.New("файл пуст")
	}

	if strings.ContainsFunc(msg, func(r rune) bool { return r == '.' || r == '-' }) {
		return morse.ToText(msg), nil
	} else {
		return morse.ToMorse(msg), nil
	}
}
