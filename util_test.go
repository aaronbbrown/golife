package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaxInt(t *testing.T) {
	assert.Equal(t, maxInt(2, 1), 2, "2 is bigger than 1")
	assert.Equal(t, maxInt(1, 2), 2, "2 is bigger than 1 reversed")
}

func TestSaneModInt(t *testing.T) {
	assert.Equal(t, saneModInt(4, 2), 0, "normal 0 remainder")
	assert.Equal(t, saneModInt(-1, 4), 3, "wrap around to other side")
}
