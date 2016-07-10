package main

import (
	//	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkStep(b *testing.B) {
	sizeX, sizeY := 200, 200
	l := NewLife("foo", sizeX, sizeY)
	for i := 0; i < b.N; i++ {
		l.generation++
		l.Step(sizeX, sizeY)
	}
}
