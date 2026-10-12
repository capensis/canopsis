package pattern

import (
	"sort"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/pattern"
)

const (
	maxReturnedSuggestions   = 10
	maxCalculatedSuggestions = 20
)

// getLiteralToFieldCombinations generates all possible combinations of literal-to-field mappings.
//
// Given a map where each literal has multiple candidate fields (sorted by relevance),
// this function produces combinations that pair each literal with one of its candidate fields.
//
// The generation uses a "progressive index" strategy to prioritize better suggestions:
//   - First, generate combinations using only index 0 (the best field) for all literals
//   - Then, generate combinations where at least one literal uses index 1
//   - Then, combinations where at least one uses index 2, and so on
//
// This ensures that combinations with higher-ranked fields appear first in the result.
func getLiteralToFieldCombinations(literalToFieldStats map[string][]LiteralFieldStats, takenFields map[string]bool) []map[string]string {
	if len(literalToFieldStats) == 0 {
		return []map[string]string{}
	}

	literals := make([]string, 0, len(literalToFieldStats))
	for literal := range literalToFieldStats {
		literals = append(literals, literal)
	}

	// ensure order
	sort.Strings(literals)

	maxFieldIndex := 0
	for _, literal := range literals {
		maxFieldIndex = max(len(literalToFieldStats[literal]), maxFieldIndex)
	}

	var combinations []map[string]string

	for currentMaxIndex := 0; currentMaxIndex < maxFieldIndex; currentMaxIndex++ {
		combinations = append(combinations, generateLiteralFieldCombinations(
			literalToFieldStats,
			literals,
			takenFields,
			currentMaxIndex,
			maxCalculatedSuggestions-len(combinations),
		)...)

		if len(combinations) >= maxCalculatedSuggestions {
			break
		}
	}

	return combinations
}

// getPatternsCombinations generates all combinations of field conditions across OR branches.
//
// For entity patterns with multiple OR branches, where each branch may have multiple
// optimization suggestions, this function combines them into complete pattern suggestions.
//
// Uses the same progressive index strategy as getLiteralToFieldCombinations to prioritize
// combinations where each branch uses its best (lowest index) suggestion.
func getPatternsCombinations(suggestions [][][]pattern.FieldCondition) []pattern.Entity {
	if len(suggestions) == 0 {
		return []pattern.Entity{}
	}

	maxSuggestionIndex := 0
	for _, suggestions := range suggestions {
		maxSuggestionIndex = max(len(suggestions)-1, maxSuggestionIndex)
	}

	var combinations []pattern.Entity

	for currentMaxIndex := 0; currentMaxIndex <= maxSuggestionIndex; currentMaxIndex++ {
		combinations = append(combinations, generatePatternCombinations(
			suggestions,
			currentMaxIndex,
			maxCalculatedSuggestions-len(combinations),
		)...)

		if len(combinations) >= maxCalculatedSuggestions {
			return combinations[:maxCalculatedSuggestions]
		}
	}

	return combinations
}

func generateLiteralFieldCombinations(
	literalToFieldStats map[string][]LiteralFieldStats,
	literals []string,
	takenFields map[string]bool,
	maxIdx,
	maxCombinations int,
) []map[string]string {
	currentIndices := make([]int, len(literals))
	var combinations []map[string]string

	var generate func(currentPos int, hasUsedMaxIdx bool) bool
	generate = func(currentPos int, hasUsedMaxIdx bool) bool {
		if len(combinations) >= maxCombinations {
			return false
		}

		if len(literals) == currentPos {
			if hasUsedMaxIdx {
				combination := make(map[string]string)
				for i := range literals {
					fieldStat := literalToFieldStats[literals[i]][currentIndices[i]]
					if !takenFields[fieldStat.FieldName] {
						combination[literals[i]] = fieldStat.FieldName
					}
				}

				if len(combination) > 0 {
					combinations = append(combinations, combination)
				}
			}

			return len(combinations) < maxCombinations
		}

		for idx := 0; idx <= min(maxIdx, len(literalToFieldStats[literals[currentPos]])-1); idx++ {
			currentIndices[currentPos] = idx
			if !generate(currentPos+1, hasUsedMaxIdx || (idx == maxIdx)) {
				return false
			}
		}

		return true
	}

	generate(0, false)

	return combinations
}

func generatePatternCombinations(
	suggestions [][][]pattern.FieldCondition,
	maxIdx,
	maxCombinations int,
) []pattern.Entity {
	currentPattern := make(pattern.Entity, len(suggestions))
	var combinations []pattern.Entity

	var generate func(currentPos int, hasUsedMaxIdx bool) bool
	generate = func(currentPos int, hasUsedMaxIdx bool) bool {
		if len(combinations) >= maxCombinations {
			return false
		}

		if len(suggestions) == currentPos {
			if hasUsedMaxIdx {
				combination := make(pattern.Entity, len(currentPattern))
				copy(combination, currentPattern)
				combinations = append(combinations, combination)
			}

			return len(combinations) < maxCombinations
		}

		for idx := 0; idx <= min(maxIdx, len(suggestions[currentPos])-1); idx++ {
			currentPattern[currentPos] = suggestions[currentPos][idx]
			if !generate(currentPos+1, hasUsedMaxIdx || (idx == maxIdx)) {
				return false
			}
		}

		return true
	}

	generate(0, false)

	return combinations
}
