package fixtures

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	keyRangeRegexp         = `^(?P<key>[\w-]+)\{(?P<from>\d+)\.\.(?P<to>\d+)\}$`
	referenceRegexp        = `^@(?P<ref>[\w\-\d\*]+)$`
	referenceFullObjRegexp = `^!(?P<ref>[\w\-\d\*]+)$`
	methodRegexp           = `^` + methodTplRegexp + `$`
	methodTplRegexp        = `<(?P<method>\w+)\((?P<args>[^)]*)\)(\.(?P<field>\w+))?>`
	tplRegexp              = `^<\{(?P<content>.+)\}>$`
	methodCurrent          = "Current"
	methodIndex            = "Index"
	methodRangeIndex       = "RangeIndex"
	methodBytes            = "Bytes"
)

type Parser interface {
	Parse(content []byte) (map[string][]interface{}, error)
}

func NewParser(faker *Faker) Parser {
	return &parser{
		faker:              faker,
		reflectFaker:       reflect.ValueOf(faker),
		keyRangeRe:         regexp.MustCompile(keyRangeRegexp),
		referenceRe:        regexp.MustCompile(referenceRegexp),
		referenceFullObjRe: regexp.MustCompile(referenceFullObjRegexp),
		methodRe:           regexp.MustCompile(methodRegexp),
		methodTplRe:        regexp.MustCompile(methodTplRegexp),
		tplRe:              regexp.MustCompile(tplRegexp),
	}
}

type parser struct {
	faker        *Faker
	reflectFaker reflect.Value

	keyRangeRe, referenceRe, referenceFullObjRe, methodRe, methodTplRe, tplRe *regexp.Regexp
}

func (p *parser) Parse(content []byte) (map[string][]interface{}, error) {
	var dataByCollection yaml.MapSlice
	decoder := yaml.NewDecoder(bytes.NewBuffer(content), yaml.UseOrderedMap())
	err := decoder.Decode(&dataByCollection)
	if err != nil {
		return nil, fmt.Errorf("cannot decode content: %w", err)
	}

	docsByCollection := make(map[string][]interface{}, len(dataByCollection))
	referenceIDs := make(map[string]interface{})
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

			matches := p.keyRangeRe.FindStringSubmatch(key)
			if len(matches) > 0 {
				refKey := matches[p.keyRangeRe.SubexpIndex("key")]
				fromStr := matches[p.keyRangeRe.SubexpIndex("from")]
				toStr := matches[p.keyRangeRe.SubexpIndex("to")]
				from, err := strconv.Atoi(fromStr)
				if err != nil {
					return nil, fmt.Errorf("from %q must be int in range %s: %w", fromStr, key, err)
				}
				to, err := strconv.Atoi(toStr)
				if err != nil {
					return nil, fmt.Errorf("to %q must be int in range %s: %w", toStr, key, err)
				}
				if from >= to {
					return nil, fmt.Errorf("from %q must be less than to %q in range %s", fromStr, toStr, key)
				}

				ids := make([]interface{}, 0, to-from+1)
				for i := from; i <= to; i++ {
					doc, err := p.processItem(index, &i, val, referenceIDs, references)
					index++
					if err != nil {
						return nil, fmt.Errorf("cannot process %s: %w", key, err)
					}

					referenceIDs[refKey+strconv.Itoa(i)] = doc["_id"]
					references[refKey+strconv.Itoa(i)] = doc
					ids = append(ids, doc["_id"])
					docs = append(docs, doc)
				}

				referenceIDs[refKey+"*"] = ids

				continue
			}

			doc, err := p.processItem(index, nil, val, referenceIDs, references)
			index++
			if err != nil {
				return nil, fmt.Errorf("cannot process %s: %w", key, err)
			}

			referenceIDs[key] = doc["_id"]
			references[key] = doc
			docs = append(docs, doc)
		}

		docsByCollection[collectionName] = docs
		p.faker.ResetUniqueName()
	}

	return docsByCollection, nil
}

func (p *parser) processItem(
	index int,
	rangeIndex *int,
	data yaml.MapSlice,
	referenceIDs, references map[string]interface{},
) (map[string]interface{}, error) {
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
			if refV, ok := referenceIDs[ref]; ok {
				if s, ok := refV.(string); ok {
					key = s
				} else {
					return nil, fmt.Errorf("not string reference %q for %q", ref, key)
				}
			} else {
				return nil, fmt.Errorf("unknown reference %q for %q", ref, key)
			}
		}

		doc[key], err = p.processValue(v.Value, index, rangeIndex, doc, referenceIDs, references)
		if err != nil {
			return nil, fmt.Errorf("cannot process %q: %w", key, err)
		}
	}

	return doc, nil
}

func (p *parser) processValue(
	fieldVal interface{},
	index int,
	rangeIndex *int,
	doc map[string]interface{},
	referenceIDs, references map[string]interface{},
) (interface{}, error) {
	switch val := fieldVal.(type) {
	case yaml.MapSlice:
		return p.processItem(index, rangeIndex, val, referenceIDs, references)
	case []interface{}:
		var err error
		newVal := make([]interface{}, len(val))
		for i := range val {
			newVal[i], err = p.processValue(val[i], index, rangeIndex, doc, referenceIDs, references)
			if err != nil {
				return nil, err
			}
		}
		return newVal, nil
	case string:
		matches := p.referenceRe.FindStringSubmatch(val)
		if len(matches) > 0 {
			ref := matches[p.referenceRe.SubexpIndex("ref")]
			newVal, ok := referenceIDs[ref]
			if !ok {
				return nil, fmt.Errorf("unknown reference %q", ref)
			}

			return newVal, nil
		}

		matches = p.referenceFullObjRe.FindStringSubmatch(val)
		if len(matches) > 0 {
			ref := matches[p.referenceFullObjRe.SubexpIndex("ref")]
			newVal, ok := references[ref]
			if !ok {
				return nil, fmt.Errorf("unknown reference %q", ref)
			}

			return newVal, nil
		}

		return p.processMethod(val, fieldVal, index, rangeIndex, doc)
	default:
		return fieldVal, nil
	}
}

func (p *parser) processMethod(val string, fieldVal interface{}, index int, rangeIndex *int, doc map[string]interface{}) (interface{}, error) {
	matches := p.methodRe.FindStringSubmatch(val)
	if len(matches) == 0 {
		matches := p.tplRe.FindStringSubmatch(val)
		if len(matches) == 0 {
			return fieldVal, nil
		}

		content := matches[p.tplRe.SubexpIndex("content")]
		var err error
		res := p.methodTplRe.ReplaceAllStringFunc(content, func(tplV string) string {
			tplRes, tplErr := p.processMethod(tplV, tplV, index, rangeIndex, doc)
			if tplErr != nil {
				err = tplErr
			}
			return fmt.Sprintf("%v", tplRes)
		})

		return res, err
	}

	method := matches[p.methodRe.SubexpIndex("method")]
	args := matches[p.methodRe.SubexpIndex("args")]
	field := matches[p.methodRe.SubexpIndex("field")]

	switch method {
	case methodCurrent:
		if field == "" {
			return nil, fmt.Errorf("%q field not defined", methodCurrent)
		}

		if fieldV, ok := doc[field]; ok {
			return fieldV, nil
		}

		return nil, fmt.Errorf("missing %q field", methodCurrent)
	case methodIndex:
		return index + 1, nil
	case methodRangeIndex:
		if rangeIndex == nil {
			return nil, errors.New("cannot use range index")
		}
		return *rangeIndex, nil
	case methodBytes:
		if args == "" {
			return nil, fmt.Errorf("%q args not defined", methodBytes)
		}
		return []byte(args), nil
	default:
		newVal, err := callReflectMethod(p.reflectFaker, method, args)
		if err != nil {
			return nil, fmt.Errorf("cannot call faker method: %w", err)
		}

		return newVal, nil
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
		if methodReflect.Type().NumIn() != len(strs) {
			return nil, fmt.Errorf("expected %d arguments for method %q but got %d", methodReflect.Type().NumIn(), method, len(strs))
		}
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
	if len(returnReflect) == 2 {
		errVal := returnReflect[1].Interface()
		if errVal == nil {
			return returnReflect[0].Interface(), nil
		}
		if err, ok := errVal.(error); ok {
			return returnReflect[0].Interface(), fmt.Errorf("method %q returned error: %w", method, err)
		}
	}

	if len(returnReflect) != 1 {
		return nil, fmt.Errorf("unexpected count of return value for method %q : expected 1 but got %d", method, len(returnReflect))
	}

	return returnReflect[0].Interface(), nil
}
