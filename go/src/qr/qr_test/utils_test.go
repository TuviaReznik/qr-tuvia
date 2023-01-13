package qr_test

import "math/rand"

var (
	alphaNumRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	stringRunes   = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ \t\n,./?!:;'()[]{}@#$%^&*-_=+")
)

func GenerateRandomAlphaNumericString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = alphaNumRunes[rand.Intn(len(alphaNumRunes))]
	}
	return string(b)
}

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = stringRunes[rand.Intn(len(stringRunes))]
	}
	return string(b)
}
