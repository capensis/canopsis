package storage

import (
	"encoding/json"
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"math"
	"sort"
	"strconv"
	"strings"
)

//TODO: refactor all applicators with the new timeBasedAlarmGroup

type timeBasedAlarmGroup struct {
	key   string
	times []int64
	ids   []string
}

func NewAlarmGroup(key string) TimeBasedAlarmGroup {
	return &timeBasedAlarmGroup{
		key:   key,
		times: make([]int64, 0),
		ids:   make([]string, 0),
	}
}

func (g timeBasedAlarmGroup) GetKey() string {
	return g.key
}

func (g timeBasedAlarmGroup) GetAlarmIds() []string {
	return g.ids
}

func (g timeBasedAlarmGroup) GetTimes() []int64 {
	return g.times
}

func (g *timeBasedAlarmGroup) GetGroupLength() int {
	return len(g.ids)
}

// in empty group returns MaxInt64 for easy minimal timestamp checks
func (g *timeBasedAlarmGroup) GetOpenTime() int64 {
	if len(g.times) != 0 {
		return g.times[0]
	}

	return math.MaxInt64
}

func (g timeBasedAlarmGroup) MarshalJSON() ([]byte, error) {
	encodedGroup := make([]string, len(g.ids))

	for i := 0; i < len(g.ids); i++ {
		encodedGroup[i] = fmt.Sprintf("%s,%d", g.ids[i], g.times[i])
	}

	return json.Marshal(encodedGroup)
}

func (g timeBasedAlarmGroup) MarshalBinary() ([]byte, error) {
	return json.Marshal(g)
}

func (g *timeBasedAlarmGroup) UnmarshalJSON(b []byte) error {
	var encodedGroup []string
	err := json.Unmarshal(b, &encodedGroup)
	if err != nil {
		return err
	}

	g.ids = make([]string, len(encodedGroup))
	g.times = make([]int64, len(encodedGroup))

	for idx, groupItem := range encodedGroup {
		split := strings.Split(groupItem, ",")
		if len(split) != 2 {
			return fmt.Errorf("group item should contain 2 elements, group: %v", split)
		}

		timestamp, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return err
		}

		g.ids[idx] = split[0]
		g.times[idx] = timestamp
	}

	return nil
}

func (g *timeBasedAlarmGroup) Push(newAlarm types.Alarm, ruleTimeInterval int64) {
	newAlarmTimestamp := newAlarm.Value.LastUpdateDate.Unix()

	for idx, v := range g.ids {
		if v == newAlarm.ID {
			g.times = append(g.times[:idx], g.times[idx+1:]...)
			g.ids = append(g.ids[:idx], g.ids[idx+1:]...)

			break
		}
	}

	// if length = 0 => init times and ids with single values
	if len(g.times) == 0 {
		g.times = []int64{newAlarmTimestamp}
		g.ids = []string{newAlarm.ID}

		return
	}

	// times is always sorted, so the 0 element is always open timestamp
	openTimestamp := g.times[0]

	//if alarm is late
	if newAlarmTimestamp < openTimestamp {
		//check if interval can be shifted, if any alarm in the Group will be lost => then we cannot shift time
		if g.times[len(g.times)-1] > newAlarmTimestamp+ruleTimeInterval {
			return
		}

		// Push to front, because it's new minimal value
		g.times = append([]int64{newAlarmTimestamp}, g.times...)
		g.ids = append([]string{newAlarm.ID}, g.ids...)

		return
	}

	newAlarmRuleStartTime := newAlarmTimestamp-ruleTimeInterval
	if newAlarmRuleStartTime > openTimestamp {
		idx := sort.Search(len(g.times), func(i int) bool {
			return g.times[i] >= newAlarmRuleStartTime
		})

		//remove outdated from front, push to the end, because it's the oldest value
		g.times = append(g.times[idx:], newAlarmTimestamp)
		g.ids = append(g.ids[idx:], newAlarm.ID)

		return
	}

	//insert to insertIdx to keep sort
	insertIdx := sort.Search(len(g.times), func(i int) bool {
		return g.times[i] >= newAlarmTimestamp
	})
	g.times = append(g.times[:insertIdx], append([]int64{newAlarmTimestamp}, g.times[insertIdx:]...)...)
	g.ids = append(g.ids[:insertIdx], append([]string{newAlarm.ID}, g.ids[insertIdx:]...)...)
}

func (g *timeBasedAlarmGroup) RemoveBefore(timestamp int64) {
	idx := sort.Search(len(g.times), func(i int) bool {
		return g.times[i] >= timestamp
	})

	g.times = g.times[idx:]
	g.ids = g.ids[idx:]
}
