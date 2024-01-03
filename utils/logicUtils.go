package utils

func SliceContains(element string, content []string) bool {
	for _, item := range content {
		if element == item {
			return true
		}
	}
	return false
}
