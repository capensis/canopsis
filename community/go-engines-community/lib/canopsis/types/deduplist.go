package types

import (
	"encoding/json"
	"fmt"
)

// DedupList is not to be used directly, use NewDedupList instead.
type DedupList struct {
	dedup     map[string]bool
	cachelist []string
}

// NewDedupList creates a new DedupList struct.
func NewDedupList(rawlist ...string) DedupList {
	i := DedupList{}
	i.reset()
	i.Add(rawlist...)
	return i
}

func (i *DedupList) reset() {
	i.dedup = make(map[string]bool)
	i.cachelist = make([]string, 0)
}

// EnsureInitialized ensure internal states to be initialized
func (i *DedupList) EnsureInitialized() {
	if i.dedup == nil {
		i.reset()
	}

	if i.cachelist == nil {
		i.recompute()
	}
}

// Add some items to the list
func (i *DedupList) Add(items ...string) {
	for _, item := range items {
		if !i.Exists(item) {
			i.dedup[item] = true
			i.cachelist = append(i.cachelist, item)
		}
	}
}

// recompute cachelist from dedup map
func (i *DedupList) recompute() {
	i.cachelist = make([]string, len(i.dedup))
	idx := 0
	for key := range i.dedup {
		i.cachelist[idx] = key
		i.dedup[key] = true
		idx++
	}
}

// Del some items from the list.
// This function is VERY slow, avoid making multiple calls since for each call
// if the DedupList has to be modified, the internal cached list is recomputed.
// Prefer call Del(items...) where items is a []string
func (i *DedupList) Del(items ...string) {
	modified := false
	for _, imp := range items {
		if i.Exists(imp) {
			delete(i.dedup, imp)
			modified = true
		}
	}

	if modified {
		i.recompute()
	}
}

// Exists returns true if the item exists in the list
func (i DedupList) Exists(item string) bool {
	_, exists := i.dedup[item]
	return exists
}

// List returns an unordered slice
func (i DedupList) List() []string {
	return i.cachelist
}

// Map returns a copy of the underlying map[string]bool
// The bool value has no meaning, it's only here so we can deduplicate using a standard map.
func (i DedupList) Map() map[string]bool {
	return i.dedup
}

// MarshalJSON implements json.Encoder
func (i DedupList) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(i.List())
	if err != nil {
		return b, fmt.Errorf("DedupList.MarshalJSON: %v", err)
	}
	return b, err
}

// UnmarshalJSON implements json.Decoder
func (i *DedupList) UnmarshalJSON(in []byte) error {
	var list []string
	i.reset()
	err := json.Unmarshal(in, &list)
	if err != nil {
		return fmt.Errorf("DedupList.UnmarshalJSON: %v", err)
	}
	i.Add(list...)
	return err
}
