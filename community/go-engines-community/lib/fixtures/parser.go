package fixtures

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/password"
	"github.com/brianvoe/gofakeit/v6"
	"gopkg.in/yaml.v2"
)

const (
	tplRegexp        = `^(?P<name>[\w\-\d]+)\s*\(template\)$`
	docRegexp        = `^(?P<name>[\w\-\d]+)$`
	docWithTplRegexp = `^(?P<name>[\w\-\d]+)\s*\(extend\s*(?P<tpl>[\w\-\d]+)\)$`
	referenceRegexp  = `^@(?P<ref>[\w\-\d]+)$`
	methodRegexp     = `^<(?P<method>\w+)\((?P<args>[^)]*)\)(\.(?P<field>\w+))?>$`
	keyCurrent       = "Current"
	keyIndex         = "Index"
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
		docRe:        regexp.MustCompile(docRegexp),
		docWithTplRe: regexp.MustCompile(docWithTplRegexp),
		tplRe:        regexp.MustCompile(tplRegexp),
		referenceRe:  regexp.MustCompile(referenceRegexp),
		methodRe:     regexp.MustCompile(methodRegexp),
	}
}

type parser struct {
	reflectFaker reflect.Value

	docRe, docWithTplRe, tplRe, referenceRe, methodRe *regexp.Regexp
}

func (p *parser) Parse(content []byte) (map[string][]interface{}, error) {
	var dataByCollection yaml.MapSlice
	err := yaml.Unmarshal(content, &dataByCollection)
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

		data, ok := collV.Value.(yaml.MapSlice)
		if !ok {
			return nil, fmt.Errorf("cannot decode content: %q must be object", collectionName)
		}

		tpls := make(map[string]yaml.MapSlice)
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

			matches := p.tplRe.FindStringSubmatch(key)
			if len(matches) > 0 {
				name := matches[p.tplRe.SubexpIndex("name")]
				tpls[name] = val

				continue
			}

			matches = p.docRe.FindStringSubmatch(key)
			if len(matches) > 0 {
				name := matches[p.docRe.SubexpIndex("name")]
				doc, err := p.processItem(index, val, references)
				index++
				if err != nil {
					return nil, fmt.Errorf("cannot process %s: %w", name, err)
				}

				references[name] = doc["_id"]
				docs = append(docs, doc)
				continue
			}

			matches = p.docWithTplRe.FindStringSubmatch(key)
			if len(matches) == 0 {
				return nil, fmt.Errorf("invalid doc key %q", key)
			}

			name := matches[p.docWithTplRe.SubexpIndex("name")]
			tplName := matches[p.docWithTplRe.SubexpIndex("tpl")]
			tpl, ok := tpls[tplName]
			if !ok {
				return nil, fmt.Errorf("unknown tpl %q", tplName)
			}

			doc, err := p.processItem(index, mergeOrderedMaps(tpl, val), references)
			index++
			if err != nil {
				return nil, fmt.Errorf("cannot process %s: %w", name, err)
			}

			references[name] = doc["_id"]
			docs = append(docs, doc)
		}

		if len(docs) > 0 {
			docsByCollection[collectionName] = docs
		}
	}

	return docsByCollection, nil
}

func (p *parser) processItem(index int, data yaml.MapSlice, references map[string]interface{}) (map[string]interface{}, error) {
	doc := make(map[string]interface{}, len(data))

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
					return nil, fmt.Errorf("not string freference %q for %q", ref, key)
				}
			} else {
				return nil, fmt.Errorf("unknown reference %q for %q", ref, key)
			}
		}

		switch val := v.Value.(type) {
		case yaml.MapSlice:
			var err error
			doc[key], err = p.processItem(index, val, references)
			if err != nil {
				return nil, fmt.Errorf("cannot process %q: %w", key, err)
			}
		case []interface{}:
			var err error
			newVal := make([]interface{}, len(val))
			for i := range val {
				if m, ok := val[i].(yaml.MapSlice); ok {
					newVal[i], err = p.processItem(index, m, references)
					if err != nil {
						return nil, fmt.Errorf("cannot process %q: %w", key, err)
					}
				} else {
					newVal[i] = val[i]
				}
			}
			doc[key] = newVal
		case string:
			matches := p.referenceRe.FindStringSubmatch(val)
			if len(matches) > 0 {
				ref := matches[p.referenceRe.SubexpIndex("ref")]
				doc[key], ok = references[ref]
				if !ok {
					return nil, fmt.Errorf("unknown reference %q for %q", ref, key)
				}
				continue
			}

			matches = p.methodRe.FindStringSubmatch(val)
			if len(matches) == 0 {
				doc[key] = v.Value
				continue
			}

			method := matches[p.methodRe.SubexpIndex("method")]
			args := matches[p.methodRe.SubexpIndex("args")]
			field := matches[p.methodRe.SubexpIndex("field")]

			switch method {
			case keyCurrent:
				if field == "" {
					return nil, fmt.Errorf("%q field not defined for %q", keyCurrent, key)
				}

				if fieldV, ok := doc[field]; ok {
					doc[key] = fieldV
				} else {
					return nil, fmt.Errorf("missing %q field for %q", keyCurrent, key)
				}
			case keyIndex:
				doc[key] = index
			default:
				var err error
				doc[key], err = callReflectMethod(p.reflectFaker, method, args)
				if err != nil {
					return nil, fmt.Errorf("cannot call faker method for %q: %w", key, err)
				}
			}
		default:
			doc[key] = v.Value
		}
	}

	return doc, nil
}

func mergeOrderedMaps(l, r yaml.MapSlice) yaml.MapSlice {
	res := make(yaml.MapSlice, len(r))
	has := make(map[interface{}]bool)

	for i, rv := range r {
		has[rv.Key] = true
		v := rv

		if rm, ok := rv.Value.(yaml.MapSlice); ok {
			for _, lv := range l {
				if lv.Key == rv.Key {
					if lm, ok := lv.Value.(yaml.MapSlice); ok {
						v = yaml.MapItem{
							Key:   rv.Key,
							Value: mergeOrderedMaps(lm, rm),
						}
						break
					}
				}
			}
		}

		res[i] = v
	}

	for _, v := range l {
		if !has[v.Key] {
			res = append(res, v)
		}
	}

	return res
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
