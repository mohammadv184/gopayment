package traits

import "fmt"

type HasDetail struct {
	Details map[string]string
}

func (i *HasDetail) Detail(key string, value string) {
	if i.Details == nil {
		i.Details = make(map[string]string)
	}
	i.Details[key] = value
}
func (i *HasDetail) GetDetail(key string) (string, error) {
	if i.Details == nil {
		return "", fmt.Errorf("no details")
	}
	if value, ok := i.Details[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("not found detail")
}
func (i *HasDetail) GetDetails() map[string]string {
	return i.Details
}
