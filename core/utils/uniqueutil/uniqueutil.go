package uniqueutil

func IsKeyUniq(key string, found map[string]struct{}) bool {
	if _, ok := found[key]; !ok {
		return true
	}
	return false
}
