package utils

func Unique(s []string) []string {
	has := make(map[string]struct{}, len(s))
	var k = 0
	for i, v := range s {
		if _, ok := has[v]; !ok {
			has[v] = struct{}{}
			s[i] = s[k]
			s[k] = v
			k++
		}
	}

	return s[:k]
}
