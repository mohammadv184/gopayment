package trait

// HasDetail is a trait that can be applied to a struct to indicate that it has a detail field.
type HasDetail struct {
	details map[string]string
}

// Detail sets the detail for the given key.
func (i *HasDetail) Detail(key, value string) {
	if i.details == nil {
		i.details = make(map[string]string)
	}
	i.details[key] = value
}

// GetDetail returns the detail for the given key.
func (i *HasDetail) GetDetail(key string) string {
	if value, ok := i.details[key]; ok {
		return value
	}
	return ""
}

// Has returns true if the given key is present in the details.
func (i *HasDetail) Has(key string) bool {
	_, has := i.details[key]
	return has
}

//GetDetails returns the details map.
func (i *HasDetail) GetDetails() map[string]string {
	return i.details
}
