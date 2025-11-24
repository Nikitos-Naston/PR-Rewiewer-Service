package internal

func DeleteElemInSlice(slice []string, elem string) []string {
	for i, el := range slice {
		if el == elem {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
