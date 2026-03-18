package crosscutting

func ContainsInSlice[ItemType comparable](slice []ItemType, item ItemType) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
