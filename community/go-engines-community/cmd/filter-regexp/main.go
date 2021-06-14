package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/types"
	"github.com/rs/zerolog"
)

const (
	condAnd   = "$and"
	condOr    = "$or"
	condRegex = "$regex"
	condId    = "_id"
)

// decode Pbhavior and check Filter has $regex
func hasFilter(decoder *json.Decoder) (bool, *types.PBehaviorLegacy, error) {
	pbh := new(types.PBehaviorLegacy)
	err := decoder.Decode(pbh)
	if err != nil {
		return false, nil, err
	}
	return strings.Contains(pbh.Filter, condRegex), pbh, nil
}

// decode Pbhavior.Filter expression, replace query
func replace(pbh *types.PBehaviorLegacy, logger zerolog.Logger) (result string, err error) {
	result = pbh.Filter
	var query interface{}
	err = json.Unmarshal([]byte(pbh.Filter), &query)
	if err != nil {
		return result, err
	}

	qmap, ok := query.(map[string]interface{})
	if !ok {
		return result, fmt.Errorf("unknown type of query %+v", query)
	}

	err = scanQuery(qmap, logger)
	if err != nil {
		return result, err
	}

	if bq, err := json.Marshal(qmap); err != nil {
		return result, err
	} else {
		result = string(bq)
	}

	logger.Info().Str("filter", fmt.Sprintf("%+v", result)).Msg("replace")
	return result, err
}

// scan for top key $or, $and
func scanQuery(query map[string]interface{}, logger zerolog.Logger) error {
	for k, v := range query {
		switch k {
		case condOr, condAnd:
			err := cond(v.([]interface{}))
			if err != nil {
				logger.Warn().Err(err).Msg("scanQuery")
			}
		default:
			return fmt.Errorf("unknown key %s %+v", k, v)
		}
	}
	return nil
}

// $and, $or query by _id conditions, replace with new matching rules
func cond(query []interface{}) (err error) {
	var newMatch *match
	for i := 0; i < len(query); i++ {
		switch qmap := query[i].(type) {
		case map[string]interface{}:
			if _, ok := qmap[condId]; ok {
				newMatch, err = getIdMatch(qmap[condId])
				if newMatch != nil {
					delete(qmap, condId)
					qmap[newMatch.k] = newMatch.v
					query[i] = qmap
				}
			} else {
				for _, condVal := range []string{condAnd, condOr} {
					conditions, ok := qmap[condVal].([]interface{})
					if ok {
						err = cond(conditions)
						qmap[condVal] = conditions
					}
				}
			}
		}
	}
	return err
}

type match struct {
	k string
	v interface{}
}

// replacement for {_id: {$regex:...}} with {component:...}
func getIdMatch(query interface{}) (*match, error) {
	switch qmap := query.(type) {
	case map[string]interface{}:
		if idRegex, ok := qmap[condRegex]; ok {
			patterns, err := getPatterns(idRegex)
			if err != nil {
				return nil, err
			}
			return &match{
				k: "component",
				v: patterns,
			}, nil
		}
	}
	return nil, fmt.Errorf("getIdMatch unsupported value %+v", query)
}

func getPatterns(regex interface{}) (s interface{}, err error) {
	switch v := regex.(type) {
	case string:
		p := strings.TrimSuffix(v, "$")
		hasSuffix := len(p) < len(v)
		p = strings.TrimPrefix(p, "/")
		sl, ok := splitList(p) // replace (a|b|c) with map[string][]string{"$in": {"a", "b", "c"}}
		if hasSuffix {
			if ok {
				if len(sl) == 1 {
					return sl[0], nil
				}
				return map[string][]string{"$in": sl}, nil
			} else if len(p)+2 == len(v) {
				return p, nil
			} else {
				// {\"_id\":{\"$regex\":\"sv-mixoraprd-91$\"}}
				return map[string]string{"$regex": p + "$"}, nil
			}
		} else {
			return map[string]string{"$regex": "^" + p}, nil
		}
	default:
		err = fmt.Errorf("unknown patterns %+v", v)
	}
	return nil, err
}

// split incoming string "(a|b|c)" as []string{a, b, c}, trim "/" items' prefixes
func splitList(v string) ([]string, bool) {
	s := strings.TrimSuffix(v, ")")
	s = strings.TrimPrefix(s, "(")
	if len(s)+2 != len(v) {
		return nil, false
	}
	l := strings.Split(s, "|")
	for i := 0; i < len(l); i++ {
		l[i] = strings.TrimPrefix(l[i], "/")
	}
	return l, true
}

// process incoming json from Reader, out processed items into Writer
func process(r io.Reader, w io.Writer, logger zerolog.Logger) error {
	decoder := json.NewDecoder(r)
	encoder := json.NewEncoder(w)
	for decoder.More() {
		ok, pbh, err := hasFilter(decoder)
		if err != nil {
			logger.Error().Err(err).Msg("decode json")
			continue
		}
		if ok {
			res, err := replace(pbh, logger)
			if err != nil {
				logger.Error().Err(err).Str("filter", pbh.Filter).Msg("replace filter")
				return err
			}

			pbh.Filter = res
		}
		err = encoder.Encode(pbh)
		if err != nil {
			logger.Error().Err(err).Msg("encode json")
		}
	}
	return nil
}

func newLogger(logLevel zerolog.Level) zerolog.Logger {
	return zerolog.New(os.Stdout).Level(logLevel).With().Timestamp().Caller().Logger()
}

func main() {
	var (
		err      error
		outFile  *os.File
		logLevel zerolog.Level
	)
	if len(os.Args) > 1 {
		if len(os.Args) > 2 {
			logLevel, err = zerolog.ParseLevel(os.Args[2])
			if err != nil {
				logLevel = zerolog.WarnLevel
			}
		}
		fileName := os.Args[1]
		outFile, err = os.Create(fileName)
		if err != nil {
			fmt.Printf("open file %s error %s\n\n", fileName, err)
			fmt.Printf("Usage: cat in_file.json | %s out_file.json warn\n", os.Args[0])
			fmt.Println("where in_file.json Pbhaviors dump")
			fmt.Println("out_file.json is file to write updated Pbehaviors")
			fmt.Println("info|warn|error - logging level, default warn")
			os.Exit(1)
		}
		defer outFile.Close()
	} else {
		outFile = os.Stdout
	}

	logger := newLogger(logLevel)
	err = process(os.Stdin, outFile, logger)
	if err != nil {
		os.Exit(1)
	}
}
