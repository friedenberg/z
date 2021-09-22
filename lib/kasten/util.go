package kasten

//TODO
//https://github.com/mitchellh/mapstructure

func getBoolOption(o map[string]interface{}, k string) bool {
	if s, ok := o[k]; ok {
		if sb, ok := s.(bool); ok {
			return sb
		}
	}

	return false
}

func getStringOption(o map[string]interface{}, k string) string {
	if s, ok := o[k]; ok {
		if sb, ok := s.(string); ok {
			return sb
		}
	}

	return ""
}
