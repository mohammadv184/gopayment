package traits

import "fmt"

type HasDetail struct {
	details map[string]string
}

func (i *HasDetail) Detail(key string, value string) {
	if i.details == nil {
		i.details = make(map[string]string)
	}
	i.details[key] = value
}
func (i *HasDetail) GetDetail(key string) (string, error) {
	if i.details == nil {
		return "", fmt.Errorf("no details")
	}
	if value, ok := i.details[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("not found detail")
}
func (i *HasDetail) GetDetails() map[string]string {
	return i.details
}
