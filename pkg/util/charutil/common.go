package charutil

// isInRange checks whether the given byte falls between the specified minimum and maximum (inclusive).
func isInRange(b, min, max byte) bool {
	return b >= min && b <= max
}
