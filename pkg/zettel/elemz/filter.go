package elemz

func OfType[T Elem](in []Elem) []T {
	var res []T
	for _, el := range in {
		if val, ok := el.(T); ok {
			res = append(res, val)
			continue
		}
	}
	return res
}
