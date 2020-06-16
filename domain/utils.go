package domain

import "math/rand"

// RandAlpha generates a length n pseudo-random string of ascii letters
// We use these as IDs.  Keep in mind that collisions, while unlikely,
// are possible.
func RandAlpha(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
