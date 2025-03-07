package router

type hasher func(string) int

func generateRollingHasher(prime int, mod int) hasher {
	powers := make(map[int]int, 100)
	powers[0] = 1

	hasher := func(s string) int {
		hash := 0
		for i, c := range s {
			power, exists := powers[i]
			if !exists {
				power = (powers[i-1] * prime) % mod
				powers[i] = power
			}

			charHash := (int(c) * power) % mod
			hash = (hash + charHash) % mod
		}
		return hash
	}

	return hasher
}

func getHasher(prime int) hasher {
	// 1e9+7 is a large well known hashingPrime
	// used widely for modulo operations in hashing
	primeMod := int(1e9 + 7)
	return generateRollingHasher(prime, primeMod)
}
