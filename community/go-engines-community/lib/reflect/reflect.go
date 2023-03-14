package reflect

import "reflect"

func ToMap(v any) (map[string]any, bool) {
	rv := reflect.ValueOf(v)
	rv = UnwrapPointer(rv)

	switch rv.Kind() {
	case reflect.Map:
		m := make(map[string]any, rv.Len()+1)
		mi := rv.MapRange()
		for mi.Next() {
			if mi.Key().Kind() != reflect.String {
				return nil, false
			}
			m[mi.Key().String()] = mi.Value().Interface()
		}

		return m, true
	case reflect.Struct:
		m := make(map[string]any, rv.Type().NumField()+1)
		StructToMap(rv, m)

		return m, true
	default:
		return nil, false
	}
}

func UnwrapPointer(v reflect.Value) reflect.Value {
	for {
		switch v.Kind() {
		case reflect.Interface, reflect.Ptr:
			v = v.Elem()
		default:
			return v
		}
	}
}

func StructToMap(v reflect.Value, data map[string]any) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		if f.Anonymous {
			fv = UnwrapPointer(fv)
			StructToMap(fv, data)
			continue
		}

		data[f.Name] = fv.Interface()
	}
}
