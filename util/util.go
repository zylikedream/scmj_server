package util

func RemoveSlice(sli []int, startPos int, length int) []int {
	remainCards := make([]int, 0, len(sli)-length)
	remainCards = append(remainCards, sli[:startPos]...)
	remainCards = append(remainCards, sli[startPos+length:]...)
	return remainCards
}

func CopyIntMap(src map[int]int) map[int]int {
	copyMap := make(map[int]int, len(src))
	for k, v := range src {
		copyMap[k] = v
	}
	return copyMap
}

func IntMapToIntSlice(src map[int]int) []int {
	var slice []int
	for k, v := range src {
		for i := 0; i < v; i++ {
			slice = append(slice, k)
		}
	}
	return slice
}
