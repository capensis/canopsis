package schema

import "strings"

func TransformToEntityMetaField(v, prefix string) string {
	if prefix != "" {
		prefix += "."
	}

	infosKey, found := strings.CutPrefix(v, "infos.")
	if found {
		return prefix + "infos->>" + "'" + infosKey + "'"
	} else {
		infosKey, found = strings.CutPrefix(v, "component_infos.")
		if found {
			return prefix + "component_infos->>" + "'" + infosKey + "'"
		}
	}

	return prefix + v
}

func GetAllowedSearchEntityMetaFields() map[string]bool {
	return map[string]bool{
		"custom_id": true,
		"name":      true,
		"type":      true,
		"category":  true,
		"component": true,
		"connector": true,
	}
}
