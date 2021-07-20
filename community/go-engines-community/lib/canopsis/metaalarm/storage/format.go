package storage

import (
	"fmt"
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"math"
	"strconv"
	"strings"
)

func DecodeGroup(encodedGroup []string) (AlarmGroup, error) {
	decodedGroup := AlarmGroup{
		OpenTime: types.CpsTime{},
		Group:    make(map[string]types.CpsTime),
	}

	minTimeStamp := int64(math.MaxInt64)

	for _, groupItem := range encodedGroup {
		split := strings.Split(groupItem, ",")

		timestamp, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return AlarmGroup{
				OpenTime: types.CpsTime{},
				Group:    nil,
			}, err
		}

		if minTimeStamp > timestamp {
			minTimeStamp = timestamp
		}

		decodedGroup.Group[split[0]] = types.NewCpsTime(timestamp)
	}

	decodedGroup.OpenTime = types.NewCpsTime(minTimeStamp)

	return decodedGroup, nil
}

func EncodeGroup(decodedGroup AlarmGroup) []string {
	var encodedGroup []string

	for alarmID, alarmTime := range decodedGroup.Group {
		encodedGroup = append(encodedGroup, fmt.Sprintf("%s,%d", alarmID, alarmTime.Unix()))
	}

	return encodedGroup
}
