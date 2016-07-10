package main

import (
	//	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkNewBoard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBoard(200, 200)
	}
}
