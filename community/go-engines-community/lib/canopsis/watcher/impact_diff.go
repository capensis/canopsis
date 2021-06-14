package watcher

// ImpactDiff is a struct that allows to easily handle changes in the impact of
// an entity. See NewImpactDiffFromImpacts for the documentation of its field.
type ImpactDiff struct {
	All      []string
	Previous map[string]bool
	Current  map[string]bool
}

// NewImpactDiffFromImpacts takes as parameters two lists of impacts (a
// previous and a current one), and returns an ImpactDiff structs containing:
//  - All: a slice containing all the impacts (previous or current)
//  - Previous: a map such that Previous[watcherID] is true if watcherID is in
//    the previousImpacts slice
//  - Current: a map such that Current[watcherID] is true if watcherID is in
//    the currentImpacts slice
//
// See impact_diff_test.go for an example.
func NewImpactDiffFromImpacts(
	previousImpacts []string,
	currentImpacts []string,
) ImpactDiff {
	all := []string{}
	previous := map[string]bool{}
	current := map[string]bool{}

	for _, impact := range previousImpacts {
		previous[impact] = true
		all = append(all, impact)
	}

	for _, impact := range currentImpacts {
		current[impact] = true
		if !previous[impact] {
			all = append(all, impact)
		}
	}

	return ImpactDiff{
		All:      all,
		Previous: previous,
		Current:  current,
	}
}
