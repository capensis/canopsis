package heartbeat

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	HeartbeatSourceDefault = ""
	HeartbeatSourceMongo   = "mongo"
)

type Heartbeat struct {
	Pattern          map[string]string `bson:"pattern"`
	ExpectedInterval string            `bson:"expected_interval"`
	Output           string            `bson:"output"`
}

func (h *Heartbeat) ToHeartBeatItem() (Item, error) {
	expectedInterval, err := time.ParseDuration(h.ExpectedInterval)
	if err != nil {
		return NewItem(0), err
	}
	hi := NewItem(expectedInterval)
	hi.Mappings = h.Pattern
	hi.Output = h.Output
	hi.Source = HeartbeatSourceMongo
	return hi, err
}

// Item contains the required informations for the HeartBeat engine
// to work.
type Item struct {
	Mappings    map[string]string `toml:"mappings"`
	MaxDuration time.Duration     `toml:"maxduration"`
	Source      string
	Output      string
}

func (li *Item) ID() string {
	return BuildID(li.Mappings)
}

func BuildID(mapping map[string]string) string {
	idvalues := make([]string, 0)

	keys := make([]string, 0)
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, field := range keys {
		idvalues = append(idvalues, field+":"+mapping[field])
	}

	return strings.Join(idvalues, ".")
}

// NewItem creates a HeartBeatItem instance for you.
func NewItem(maxDuration time.Duration) Item {
	return Item{
		Mappings:    make(map[string]string, 0),
		MaxDuration: maxDuration,
		Source:      HeartbeatSourceDefault,
	}
}

// AddMapping simplifies your life to declare a useful HeartBeatItem.
// You can safely ignore the error, it only says that the mapping
// you try to insert already exists.
func (li *Item) AddMapping(field, matches string) error {
	if _, ok := li.Mappings[field]; ok {
		return errors.New("mapping already exists")
	}
	li.Mappings[field] = matches
	return nil
}

type SafeHeartbeatItems struct {
	v   []Item
	mux sync.Mutex
}

func NewSafeHeartbeatItems() SafeHeartbeatItems {
	return SafeHeartbeatItems{v: make([]Item, 0)}
}

func (i *SafeHeartbeatItems) Value() []Item {
	i.mux.Lock()
	defer i.mux.Unlock()
	v := make([]Item, len(i.v))
	copy(v, i.v)
	return v
}

func (i *SafeHeartbeatItems) AddItem(hi Item) {
	i.mux.Lock()
	defer i.mux.Unlock()
	i.v = append(i.v, hi)
}

func (i *SafeHeartbeatItems) Clear() {
	i.mux.Lock()
	defer i.mux.Unlock()
	i.v = make([]Item, 0)
}

func (i *SafeHeartbeatItems) Len() int {
	i.mux.Lock()
	defer i.mux.Unlock()
	return len(i.v)
}
