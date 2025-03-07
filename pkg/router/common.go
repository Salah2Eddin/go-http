package router

func isWildcard(part string) bool {
	wildcard := byte('*')
	return len(part) == 1 && part[0] == wildcard
}
