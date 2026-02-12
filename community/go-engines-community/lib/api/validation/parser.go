package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type duplicateErrorParser struct {
	dupErrorRegexp *regexp.Regexp
}

type DuplicateErrorParser interface {
	Parse(err error, validatedStruct any) error
}

func NewDuplicateErrorParser() DuplicateErrorParser {
	return &duplicateErrorParser{
		dupErrorRegexp: regexp.MustCompile(`{ ([^:]+)`),
	}
}

func (p *duplicateErrorParser) Parse(err error, validatedStruct any) error {
	match := p.dupErrorRegexp.FindStringSubmatch(err.Error())
	if len(match) > 1 {
		field := match[1]
		rt := reflect.ValueOf(validatedStruct).Type()
		rf := p.findStructFieldByTag(rt, "bson", field)
		if rf != nil {
			return NewSingleError("exist", rf.Name, rf.Name, validatedStruct)
		}
	}

	return fmt.Errorf("can't parse duplication error: %w", err)
}

func (p *duplicateErrorParser) findStructFieldByTag(rt reflect.Type, tagKey, tagVal string) *reflect.StructField {
	switch rt.Kind() {
	case reflect.Ptr, reflect.Interface:
		return p.findStructFieldByTag(rt.Elem(), tagKey, tagVal)
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			rf := rt.Field(i)
			if rf.Anonymous {
				ft := p.findStructFieldByTag(rf.Type, tagKey, tagVal)
				if ft != nil {
					return ft
				}

				continue
			}

			ft := strings.Split(rf.Tag.Get(tagKey), ",")
			if ft[0] == tagVal {
				return &rf
			}
		}
	}

	return nil
}
