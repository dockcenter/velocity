package slices

func Contains[T comparable](slice []T, searchElement T) bool {
	for _, e := range slice {
		if e == searchElement {
			return true
		}
	}
	return false
}
