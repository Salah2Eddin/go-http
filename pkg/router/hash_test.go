package router

import (
	"testing"
)

func TestGenerateRollingHasherConsistency(t *testing.T) {
	hasher := generateRollingHasher(31, 1000000007)
	hash1 := hasher("test")
	hash2 := hasher("test")
	if hash1 != hash2 {
		t.Errorf("hasher(\"test\") = %d, then %d, want same hash", hash1, hash2)
	}
}

func TestGenerateRollingHasherDifferentStrings(t *testing.T) {
	hasher := generateRollingHasher(31, 1000000007)
	hash1 := hasher("test1")
	hash2 := hasher("test2")
	if hash1 == hash2 {
		t.Errorf("different strings produced same hash")
	}
}

func TestGenerateRollingHasherEmptyString(t *testing.T) {
	hasher := generateRollingHasher(31, 1000000007)
	result := hasher("")
	if result < 0 {
		t.Errorf("hasher(\"\") = %d, want non-negative", result)
	}
}

func TestGetHasher(t *testing.T) {
	hasher := getHasher(31)
	if hasher == nil {
		t.Errorf("getHasher returned nil")
	}
}

func TestGetHasherConsistency(t *testing.T) {
	hasher := getHasher(31)
	hash1 := hasher("example")
	hash2 := hasher("example")
	if hash1 != hash2 {
		t.Errorf("hasher(\"example\") = %d, then %d, want same hash", hash1, hash2)
	}
}
