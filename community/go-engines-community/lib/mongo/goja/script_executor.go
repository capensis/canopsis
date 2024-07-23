package goja

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"
	"github.com/dop251/goja"
)

func NewScriptExecutor(dbClient mongo.DbClient) mongo.ScriptExecutor {
	return &scriptExecutor{
		dbClient: dbClient,
	}
}

type scriptExecutor struct {
	dbClient mongo.DbClient
}

func (e *scriptExecutor) Exec(ctx context.Context, file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read file %q: %w", file, err)
	}

	jsCode := string(b)
	prg, err := goja.Compile(file, jsCode, true)
	if err != nil {
		return fmt.Errorf("cannot compile js %q: %w", file, err)
	}

	collectionNames, err := e.findUsedCollection(jsCode)
	if err != nil {
		return fmt.Errorf("cannot parse js %q: %w", file, err)
	}

	vm := goja.New()
	globalFuncs := e.getGlobalFuncs()
	for name, f := range globalFuncs {
		err = vm.GlobalObject().Set(name, f)
		if err != nil {
			return fmt.Errorf("cannot set js %s global func in %q: %w", name, file, err)
		}
	}

	dbClient := &jsDatabase{
		dbClient: e.dbClient,
		vm:       vm,
	}
	err = vm.GlobalObject().Set("db", dbClient.getMethods(ctx, collectionNames))
	if err != nil {
		return fmt.Errorf("cannot set js global var in %q: %w", file, err)
	}

	_, err = vm.RunProgram(prg)
	if err != nil {
		return fmt.Errorf("cannot execute js %q: %w", file, err)
	}

	return nil
}

func (e *scriptExecutor) findUsedCollection(jsCode string) ([]string, error) {
	re, err := regexp.Compile(`db\.(\w+)\.`)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllString(jsCode, -1)
	collectionNames := make([]string, 0, len(matches))
	exists := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		name := strings.TrimSuffix(strings.TrimPrefix(match, "db."), ".")
		if _, ok := exists[name]; !ok {
			exists[name] = struct{}{}
			collectionNames = append(collectionNames, name)
		}
	}

	return collectionNames, nil
}

func (e *scriptExecutor) getGlobalFuncs() map[string]any {
	return map[string]any{
		"genID": func() string {
			return utils.NewID()
		},
		"isInt": func(v any) bool {
			switch v.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				return true
			}

			return false
		},
		"toInt": func(v any) int64 {
			switch castedV := v.(type) {
			case int:
				return int64(castedV)
			case int8:
				return int64(castedV)
			case int16:
				return int64(castedV)
			case int32:
				return int64(castedV)
			case int64:
				return castedV
			case uint:
				return int64(castedV)
			case uint8:
				return int64(castedV)
			case uint16:
				return int64(castedV)
			case uint32:
				return int64(castedV)
			case uint64:
				return int64(castedV)
			case float32:
				return int64(castedV)
			case float64:
				return int64(castedV)
			case string:
				i, err := strconv.ParseInt(castedV, 10, 64)
				if err == nil {
					return i
				}

				f, err := strconv.ParseFloat(castedV, 64)
				if err == nil {
					return int64(f)
				}
			case bool:
				if castedV {
					return 1
				}

				return 0
			}

			return 0
		},
	}
}
