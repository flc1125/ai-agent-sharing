package dg

import (
	"math/rand/v2"
	"testing"
)

func TestDAG(t *testing.T) {
	for {
		n := rand.IntN(100)

		if n < 5 {
			t.Logf("random number is less than 5")
			return // exit the loop
		}

		t.Logf("random number is %d", n)
	}
}

func TestDG(t *testing.T) {
	if rand.IntN(10) < 5 {
		t.Logf("random number is less than 5")
		return
	}

	for {
		n := rand.IntN(100)

		t.Logf("random number is %d", n)
	}
}
