package htmlparsing

import "net/url"

// URLValues converts a string map to URL values suitable for form submission
func URLValues(parameters map[string]string) url.Values {
	values := make(url.Values)

	for key := range parameters {
		values.Set(key, parameters[key])
	}

	return values
}
