package pattern

import (
	"errors"
	"fmt"
	"regexp/syntax"
)

const (
	literalLenLimit  = 255
	groupSizeLimit   = 10
	groupNumberLimit = 10
)

var ErrTooLongLiteral = fmt.Errorf("literal is too long, max length is %d", literalLenLimit)
var ErrTooLongGroup = fmt.Errorf("literal group is too long, max length is %d", groupSizeLimit)
var ErrTooManyGroups = fmt.Errorf("too many literal groups, max number is %d", groupNumberLimit)
var ErrRegexpNotSupported = errors.New("regexp is not supported")

// parseLiterals extracts literal strings from a parsed regular expression syntax tree.
//
// This function should only be used in the context of entity pattern regexp conditions optimization,
// where the extracted literal groups are used to generate faster conditions (such as exact matches).
// That can leverage database indexes before applying the full regexp match.
//
// The function returns a slice of slices ([][]string) where:
//   - The outer slice represents different sets of literals, each set may be related to a specific entity field
//   - Each inner slice contains all possible literal values for that set
func parseLiterals(re *syntax.Regexp) ([][]string, error) {
	var literalGroups [][]string

	switch re.Op {
	case syntax.OpLiteral:
		literal := string(re.Rune)
		if len(literal) > literalLenLimit {
			return nil, ErrTooLongLiteral
		}

		literalGroups = append(literalGroups, []string{literal})
	case syntax.OpConcat:
		var buildingGroup bool

		for _, sub := range re.Sub {
			subLiteralGroups, err := parseLiterals(sub)
			if err != nil {
				return nil, err
			}

			if len(subLiteralGroups) != 0 {
				if !buildingGroup {
					// index = 0, because it seems not possible to have concatenations inside a concatenation,
					// so there can't be more than a single group from concatenation's subexpressions.
					literalGroups = append(literalGroups, subLiteralGroups[0])
					if len(literalGroups) > groupNumberLimit {
						return nil, ErrTooManyGroups
					}

					buildingGroup = true
				} else {
					var expandedGroup []string

					for _, literalGroup := range literalGroups[len(literalGroups)-1] {
						for _, subLiteralGroup := range subLiteralGroups {
							for _, subLiteral := range subLiteralGroup {
								expandedLiteral := literalGroup + subLiteral
								if len(expandedLiteral) > literalLenLimit {
									return nil, ErrTooLongLiteral
								}

								expandedGroup = append(expandedGroup, expandedLiteral)
								if len(expandedGroup) > groupSizeLimit {
									return nil, ErrTooLongGroup
								}
							}
						}
					}

					literalGroups[len(literalGroups)-1] = expandedGroup
				}
			} else {
				// if no groups returned, means that no more literals or alterations found in concatenation for the group,
				// then building is complete.
				buildingGroup = false
			}
		}
	case syntax.OpCapture:
		if re.Op == syntax.OpCapture && len(re.Sub) > 0 {
			var err error

			literalGroups, err = parseLiterals(re.Sub[0])
			if err != nil {
				return nil, err
			}
		}
	case syntax.OpAlternate:
		// alternate subexpressions might have concatenations, those concatenations might have the same literals inside
		// to avoid literal duplication, check existence in the map
		uniqueLiteralsByGroup := make(map[string]bool)

		for _, sub := range re.Sub {
			subGroups, err := parseLiterals(sub)
			if err != nil {
				return nil, err
			}

			for i := len(literalGroups); i < len(subGroups); i++ {
				literalGroups = append(literalGroups, []string{})
				if len(literalGroups) > groupNumberLimit {
					return nil, ErrTooManyGroups
				}
			}

			for subGroupIdx, subGroup := range subGroups {
				for _, subLiteral := range subGroup {
					if uniqueLiteralsByGroup[subLiteral] {
						continue
					}

					uniqueLiteralsByGroup[subLiteral] = true
					literalGroups[subGroupIdx] = append(literalGroups[subGroupIdx], subLiteral)
					if len(literalGroups[subGroupIdx]) > groupSizeLimit {
						return nil, ErrTooLongGroup
					}
				}
			}
		}
	case syntax.OpStar, syntax.OpPlus:
		if len(re.Sub) == 1 && re.Sub[0].Op == syntax.OpAnyCharNotNL {
			return nil, nil
		}

		return nil, ErrRegexpNotSupported
	case syntax.OpNoMatch, syntax.OpEmptyMatch, syntax.OpAnyChar,
		syntax.OpAnyCharNotNL, syntax.OpQuest, syntax.OpWordBoundary, syntax.OpNoWordBoundary, syntax.OpRepeat:
		return nil, ErrRegexpNotSupported
	case syntax.OpCharClass:
		runeLen := len(re.Rune)
		if runeLen%2 != 0 {
			// shouldn't happen, the OpCharClass usually contains a list of pairs, but just in case to avoid a surprise "index out of range" panic
			return nil, ErrRegexpNotSupported
		}

		literalGroups = append(literalGroups, []string{})

		for i := 0; i < runeLen; i += 2 {
			for j := re.Rune[i]; j <= re.Rune[i+1]; j++ {
				literalGroups[0] = append(literalGroups[0], string(j))
				if len(literalGroups[0]) > groupSizeLimit {
					return literalGroups, ErrTooLongGroup
				}
			}
		}
	}

	return literalGroups, nil
}
