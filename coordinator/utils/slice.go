package utils

func InStringSlice(sl []string, str string) bool {
	for _, s := range sl {
		if s == str {
			return true
		}
	}

	return false
}
