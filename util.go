package main

import (
	"math/rand"
)

// weighted coin toss.  odds is the change
// that the toss will appear true
func weightedRandBool(odds int) bool {
	return rand.Intn(odds) != 0
}

// a sane modulo operator that works like every other damn language
// see https://github.com/golang/go/issues/448
// https://groups.google.com/forum/#!topic/golang-nuts/xj7CV857vAg
func saneModInt(x, y int) int {
	result := x % y
	if result < 0 {
		result += y
	}
	return result
}

// go doesn't have a Int.max function
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
