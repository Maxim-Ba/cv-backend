package router

func optionalStringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
