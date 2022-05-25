package fixtures

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/goccy/go-yaml"
)

const (
	referenceRegexp = `^@(?P<ref>[\w\-\d]+)$`
	methodRegexp    = `^<(?P<method>\w+)\((?P<args>[^)]*)\)(\.(?P<field>\w+))?>$`
	keyCurrent      = "Current"
	keyIndex        = "Index"
)

type Parser interface {
	Parse(content []byte) (map[string][]interface{}, error)
}

func NewParser(passwordEncoder password.Encoder) Parser {
	faker := Faker{
		Faker:           gofakeit.New(0),
		passwordEncoder: passwordEncoder,
	}

	return &parser{
		reflectFaker: reflect.ValueOf(faker),
		referenceRe:  regexp.MustCompile(referenceRegexp),
		methodRe:     regexp.MustCompile(methodRegexp),
	}
}

type parser struct {
	reflectFaker reflect.Value

	referenceRe, methodRe *regexp.Regexp
}

func (p *parser) Parse(content []byte) (map[string][]interface{}, error) {
	var dataByCollection yaml.MapSlice
	decoder := yaml.NewDecoder(bytes.NewBuffer(content), yaml.UseOrderedMap())
	err := decoder.Decode(&dataByCollection)
	if err != nil {
		return nil, fmt.Errorf("cannot decode content: %w", err)
	}

	docsByCollection := make(map[string][]interface{}, len(dataByCollection))
	references := make(map[string]interface{})

	for _, collV := range dataByCollection {
		collectionName, ok := collV.Key.(string)
		if !ok {
			return nil, fmt.Errorf("%+v not string key", collV.Key)
		}

		if collectionName == "template" {
			continue
		}

		data, ok := collV.Value.(yaml.MapSlice)
		if !ok {
			return nil, fmt.Errorf("cannot decode content: %q must be object", collectionName)
		}

		docs := make([]interface{}, 0, len(data))
		index := 0

		for _, v := range data {
			key, ok := v.Key.(string)
			if !ok {
				return nil, fmt.Errorf("%+v not string key", v.Key)
			}

			val, ok := v.Value.(yaml.MapSlice)
			if !ok {
				return nil, fmt.Errorf("cannot decode content: %q must be object", key)
			}

			doc, err := p.processItem(index, val, references)
			index++
			if err != nil {
				return nil, fmt.Errorf("cannot process %s: %w", key, err)
			}

			references[key] = doc["_id"]
			docs = append(docs, doc)
		}

		docsByCollection[collectionName] = docs
	}

	return docsByCollection, nil
}

func (p *parser) processItem(index int, data yaml.MapSlice, references map[string]interface{}) (map[string]interface{}, error) {
	doc := make(map[string]interface{}, len(data))
	var err error

	for _, v := range data {
		key, ok := v.Key.(string)
		if !ok {
			return nil, fmt.Errorf("%+v not string key", v.Key)
		}

		matches := p.referenceRe.FindStringSubmatch(key)
		if len(matches) > 0 {
			ref := matches[p.referenceRe.SubexpIndex("ref")]
			if refV, ok := references[ref]; ok {
				if s, ok := refV.(string); ok {
					key = s
				} else {
					return nil, fmt.Errorf("not string reference %q for %q", ref, key)
				}
			} else {
				return nil, fmt.Errorf("unknown reference %q for %q", ref, key)
			}
		}

		doc[key], err = p.processValue(v.Value, index, doc, references)
		if err != nil {
			return nil, fmt.Errorf("cannot process %q: %w", key, err)
		}
	}

	return doc, nil
}

func (p *parser) processValue(fieldVal interface{}, index int, doc map[string]interface{}, references map[string]interface{}) (interface{}, error) {
	switch val := fieldVal.(type) {
	case yaml.MapSlice:
		return p.processItem(index, val, references)
	case []interface{}:
		var err error
		newVal := make([]interface{}, len(val))
		for i := range val {
			newVal[i], err = p.processValue(val[i], index, doc, references)
			if err != nil {
				return nil, err
			}
		}
		return newVal, nil
	case string:
		matches := p.referenceRe.FindStringSubmatch(val)
		if len(matches) > 0 {
			ref := matches[p.referenceRe.SubexpIndex("ref")]
			newVal, ok := references[ref]
			if !ok {
				return nil, fmt.Errorf("unknown reference %q", ref)
			}

			return newVal, nil
		}

		matches = p.methodRe.FindStringSubmatch(val)
		if len(matches) == 0 {
			return fieldVal, nil
		}

		method := matches[p.methodRe.SubexpIndex("method")]
		args := matches[p.methodRe.SubexpIndex("args")]
		field := matches[p.methodRe.SubexpIndex("field")]

		switch method {
		case keyCurrent:
			if field == "" {
				return nil, fmt.Errorf("%q field not defined", keyCurrent)
			}

			if fieldV, ok := doc[field]; ok {
				return fieldV, nil
			}

			return nil, fmt.Errorf("missing %q field", keyCurrent)
		case keyIndex:
			return index, nil
		default:
			newVal, err := callReflectMethod(p.reflectFaker, method, args)
			if err != nil {
				return nil, fmt.Errorf("cannot call faker method: %w", err)
			}

			return newVal, nil
		}
	default:
		return fieldVal, nil
	}
}

func callReflectMethod(rv reflect.Value, method, args string) (interface{}, error) {
	methodReflect := rv.MethodByName(method)
	if !methodReflect.IsValid() {
		return nil, fmt.Errorf("unexpected method %q", method)
	}

	in := make([]reflect.Value, 0)
	if args != "" {
		strs := strings.Split(args, ",")
		for i, s := range strs {
			switch methodReflect.Type().In(i).Kind() {
			case reflect.Int:
				vi, err := strconv.Atoi(s)
				if err != nil {
					return nil, fmt.Errorf("cannot parse %q as int %d argument for method %q: %w", s, i, method, err)
				}
				in = append(in, reflect.ValueOf(vi))
			case reflect.Bool:
				b, err := strconv.ParseBool(s)
				if err != nil {
					return nil, fmt.Errorf("cannot parse %q as bool %d argument for method %q: %w", s, i, method, err)
				}
				in = append(in, reflect.ValueOf(b))
			case reflect.String:
				in = append(in, reflect.ValueOf(s))
			default:
				return nil, fmt.Errorf("unknown %d argument type %q for method %q", i, methodReflect.Type().In(i).Kind().String(), method)
			}
		}
	}

	returnReflect := methodReflect.Call(in)
	if len(returnReflect) != 1 {
		return nil, fmt.Errorf("unexpected count of return value for method %q : expected 1 but got %d", method, len(returnReflect))
	}

	return returnReflect[0].Interface(), nil
}
