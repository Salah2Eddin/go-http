package router

import (
	"testing"
)

func TestGenerateIncIDGenerator(t *testing.T) {
	gen := generateIncIDGenerator()
	if gen == nil {
		t.Errorf("generateIncIDGenerator returned nil")
	}
}

func TestGenerateIncIDGeneratorSequential(t *testing.T) {
	gen := generateIncIDGenerator()
	id1 := gen()
	id2 := gen()
	id3 := gen()

	if id1 != 1 {
		t.Errorf("first call = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("second call = %d, want 2", id2)
	}
	if id3 != 3 {
		t.Errorf("third call = %d, want 3", id3)
	}
}

func TestGenerateIncIDGeneratorIndependent(t *testing.T) {
	gen1 := generateIncIDGenerator()
	gen2 := generateIncIDGenerator()

	id1 := gen1()
	id2 := gen2()

	if id1 != 1 || id2 != 1 {
		t.Errorf("independent generators should start at 1, got %d and %d", id1, id2)
	}
}

func TestGetIDGenerator(t *testing.T) {
	gen := getIDGenerator()
	result := gen()
	if result != 1 {
		t.Errorf("getIDGenerator() first call = %d, want 1", result)
	}
}
