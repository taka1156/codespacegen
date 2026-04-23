package entity

func stringValue(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}

func boolValue(v *bool) bool {
	if v == nil {
		return false
	}

	return *v
}
