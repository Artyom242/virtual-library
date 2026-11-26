package validation

func IsIDMoreThanOne(id int) bool {
	return id > 1
}

func IsUniqueValues(slice []int) bool {
	values := make(map[int]bool)

	for _, value := range slice {
		if _, ok := values[value]; ok {
			return false
		}
		values[value] = true
	}

	return true
}
