package crosscutting

type Collection[ItemType any] struct {
	items []ItemType
}

func (collection *Collection[ItemType]) AddItem(items ...ItemType) *Collection[ItemType] {
	collection.items = append(collection.items, items...)
	return collection
}

func (collection Collection[ItemType]) GetItems() []ItemType {
	if collection.items == nil {
		return make([]ItemType, 0)
	}
	return collection.items
}
