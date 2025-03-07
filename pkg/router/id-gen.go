package router

type idGenerator func() int

func generateIncIDGenerator() idGenerator {
	id := 0
	gen := func() int {
		id++
		return id
	}
	return gen
}

func getIDGenerator() idGenerator {
	return generateIncIDGenerator()
}
