package util

import (
	"bytes"
	"io"
	"testing"
)

func TestPeek(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expected    byte
		expectError bool
	}{
		{"single byte", []byte{'a'}, 'a', false},
		{"multiple bytes", []byte{'H', 'e', 'l', 'l', 'o'}, 'H', false},
		{"empty reader", []byte{}, 0, true},
		{"null byte", []byte{0x00}, 0x00, false},
		{"space", []byte{' '}, ' ', false},
		{"digit", []byte{'5'}, '5', false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reader := bytes.NewReader(tc.data)
			result, err := Peek(reader)

			if tc.expectError {
				if err == nil {
					t.Errorf("Peek() expected error, got nil")
				}
				if err != io.EOF {
					t.Errorf("Peek() expected io.EOF, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Peek() unexpected error: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Peek() = %q, want %q", result, tc.expected)
				}
			}
		})
	}
}

func TestPeekDoesNotConsume(t *testing.T) {
	data := []byte{'H', 'e', 'l', 'l', 'o'}
	reader := bytes.NewReader(data)

	// First peek
	first, err := Peek(reader)
	if err != nil {
		t.Fatalf("Peek() unexpected error: %v", err)
	}
	if first != 'H' {
		t.Errorf("Peek() = %q, want 'H'", first)
	}

	// Second peek should return the same byte
	second, err := Peek(reader)
	if err != nil {
		t.Fatalf("Peek() unexpected error on second call: %v", err)
	}
	if second != 'H' {
		t.Errorf("Peek() second call = %q, want 'H'", second)
	}

	// ReadByte should still return 'H' since peek didn't consume it
	third, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte() unexpected error: %v", err)
	}
	if third != 'H' {
		t.Errorf("ReadByte() = %q, want 'H'", third)
	}

	// Next byte should be 'e'
	fourth, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte() unexpected error: %v", err)
	}
	if fourth != 'e' {
		t.Errorf("ReadByte() = %q, want 'e'", fourth)
	}
}

func TestPeekMultipleCalls(t *testing.T) {
	data := []byte{'a', 'b', 'c'}
	reader := bytes.NewReader(data)

	// Multiple peeks should all return 'a'
	for i := 0; i < 5; i++ {
		result, err := Peek(reader)
		if err != nil {
			t.Fatalf("Peek() call %d unexpected error: %v", i+1, err)
		}
		if result != 'a' {
			t.Errorf("Peek() call %d = %q, want 'a'", i+1, result)
		}
	}

	// After peeking, reading should still work
	first, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte() unexpected error: %v", err)
	}
	if first != 'a' {
		t.Errorf("ReadByte() = %q, want 'a'", first)
	}
}

func TestPeekAfterRead(t *testing.T) {
	data := []byte{'a', 'b', 'c'}
	reader := bytes.NewReader(data)

	// Read first byte
	first, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte() unexpected error: %v", err)
	}
	if first != 'a' {
		t.Errorf("ReadByte() = %q, want 'a'", first)
	}

	// Peek should now return 'b'
	peeked, err := Peek(reader)
	if err != nil {
		t.Fatalf("Peek() unexpected error: %v", err)
	}
	if peeked != 'b' {
		t.Errorf("Peek() = %q, want 'b'", peeked)
	}

	// Read should still return 'b'
	second, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte() unexpected error: %v", err)
	}
	if second != 'b' {
		t.Errorf("ReadByte() = %q, want 'b'", second)
	}
}

func TestPeekEmptyReader(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	result, err := Peek(reader)

	if err == nil {
		t.Error("Peek() on empty reader expected error, got nil")
	}
	if err != io.EOF {
		t.Errorf("Peek() expected io.EOF, got %v", err)
	}
	if result != 0 {
		t.Errorf("Peek() on empty reader = %q, want 0", result)
	}
}

func TestPeekAtEnd(t *testing.T) {
	reader := bytes.NewReader([]byte{'a'})

	// Read the only byte
	_, err := reader.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte() unexpected error: %v", err)
	}

	// Peek should now return EOF
	result, err := Peek(reader)
	if err == nil {
		t.Error("Peek() at end expected error, got nil")
	}
	if err != io.EOF {
		t.Errorf("Peek() expected io.EOF, got %v", err)
	}
	if result != 0 {
		t.Errorf("Peek() at end = %q, want 0", result)
	}
}
