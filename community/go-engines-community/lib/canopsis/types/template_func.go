package types

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"text/template"
)

func getTemplateFunc() template.FuncMap {
	return template.FuncMap{
		// json will convert an item to an JSON-compatible element,
		// ie ints will be returned as integers and strings returned as strings with quotes
		// It is mostly used for the strings to preserve their content
		// and avoid Go behavior to escape special characters strings by default
		"json": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return err.Error()
			}

			return string(b)
		},
		// Same behavior as json but remove the quotes around the string.
		"json_unquote": func(v interface{}) string {
			b, err := json.Marshal(v)
			if err != nil {
				return err.Error()
			}

			jsonStr := string(b)

			if string(jsonStr[0]) != "\"" || string(jsonStr[len(jsonStr)-1]) != "\"" {
				log.Printf("json_unquote : trying to unquote %+v. Returning value directly", jsonStr)
				return jsonStr
			}

			return jsonStr[1 : len(jsonStr)-1]
		},
		// split will split a string according a separator and returns the substring
		"split": func(sep string, index int, v interface{}) string {
			if s, ok := v.(string); ok {
				stringSlice := strings.Split(s, sep)
				if 0 <= index && index < len(stringSlice) {
					return stringSlice[index]
				}
				log.Printf("split : index %+v out of bounds", index)
			} else {
				log.Printf("split : %+v is not a string", v)
			}
			return ""
		},
		// trim will return a string with all leading and trailing white space removed
		"trim": func(v interface{}) string {
			if s, ok := v.(string); ok {
				return strings.TrimSpace(s)
			}
			log.Printf("trim : %+v is not a string", v)
			return ""
		},
		// formattedDate will return a formatted string from a CpsTime
		"formattedDate": func(format string, v interface{}) string {
			if t, ok := v.(CpsTime); ok {
				return t.Time.Format(format)
			}
			log.Printf("formattedDate : %+v is not a CpsTime", v)
			return ""
		},
		// replace will replace a string, replacing matches of the regex with the replacement string
		"replace": func(oldRegex string, new string, v interface{}) string {
			if s, ok := v.(string); ok {
				re, err := regexp.Compile(oldRegex)
				if err != nil {
					log.Printf("replace : %+v cannot be parsed by regexp, %v", oldRegex, err)
					return ""
				}
				return re.ReplaceAllString(s, new)
			}
			log.Printf("replace : %+v is not a string", v)
			return ""
		},
	}
}
