package util

func ArrayFindIndex(items []string, findItem string) int {
	for index, item := range items {
		if item == findItem {
			return index
		}
	}
	return -1
}

func ArrayRemoveItem(items []string, findItem string) []string {
	newArray := make([]string, 0)
	for _, item := range items {
		if item == findItem {
			continue
		}
		newArray = append(newArray, item)
	}
	return newArray
}
