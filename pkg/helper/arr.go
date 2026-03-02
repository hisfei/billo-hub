package helper

// Deduplicate removes duplicate elements from a slice of any comparable type using generics.
// It preserves the original order of the elements.
//
// Example:
//
//	Deduplicate([]string{"a", "b", "a"}) -> []string{"a", "b"}
//	Deduplicate([]int{1, 2, 1}) -> []int{1, 2}
func Deduplicate[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return slice
	}

	encountered := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, ok := encountered[v]; !ok {
			encountered[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Contains checks if a slice contains the specified element.
// It uses generics and is applicable to any comparable type.
func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}
