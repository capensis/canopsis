package watcher_test

import (
	"encoding/json"
	"fmt"

	"git.canopsis.net/canopsis/go-engines/lib/canopsis/watcher"
)

func ExampleNewImpactDiffFromImpacts() {
	previousImpacts := []string{"a", "b"}
	currentImpacts := []string{"b", "c"}

	impacts := watcher.NewImpactDiffFromImpacts(
		previousImpacts,
		currentImpacts)

	json, err := json.MarshalIndent(impacts, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(json))
	// Output:
	// {
	//   "All": [
	//     "a",
	//     "b",
	//     "c"
	//   ],
	//   "Previous": {
	//     "a": true,
	//     "b": true
	//   },
	//   "Current": {
	//     "b": true,
	//     "c": true
	//   }
	// }
}
