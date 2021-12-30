package trait

import "fmt"

// HasDetail is a trait that can be applied to a struct to indicate that it has a detail field.
type HasDetail struct {
	details map[string]string
}

// Detail sets the detail for the given key.
func (i *HasDetail) Detail(key string, value string) {
	if i.details == nil {
		i.details = make(map[string]string)
	}
	i.details[key] = value
}

// GetDetail returns the detail for the given key.
func (i *HasDetail) GetDetail(key string) (string, error) {
	if i.details == nil {
		return "", fmt.Errorf("no details")
	}
	if value, ok := i.details[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("not found detail")
}

//GetDetails returns the details map.
func (i *HasDetail) GetDetails() map[string]string {
	return i.details
}
