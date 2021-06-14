package stats

import (
	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	libmongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	mock_mongo "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/mocks/lib/mongo"
	"github.com/golang/mock/gomock"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestManager_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for name, data := range getPingDataSets() {
		mockDbCollection := mock_mongo.NewMockDbCollection(ctrl)
		mockDbCollection.
			EXPECT().
			Find(gomock.Any(), gomock.Any()).
			Return(mockCursor(ctrl, data.origin), nil)
		mockDbCollection.
			EXPECT().
			UpdateOne(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&mongo.UpdateResult{}, nil)
		mockDbClient := mock_mongo.NewMockDbClient(ctrl)
		mockDbClient.EXPECT().Collection(libmongo.SessionStatsMongoCollection).Return(mockDbCollection)
		m := NewManager(mockDbClient, time.Minute)

		s, err := m.Ping(data.session, data.path)

		if err != nil {
			t.Errorf("%s: expected no error but got %v", name, err)
		}

		assertStats(t, data.expected, *s, name)
	}
}

type pingDataSet struct {
	origin   *Stats
	session  SessionData
	path     PathData
	expected Stats
}

func getPingDataSets() map[string]pingDataSet {
	now := types.CpsTime{Time: time.Now()}
	minuteAgo := types.CpsTime{Time: time.Now().Add(-time.Minute)}
	halfminuteAgo := types.CpsTime{Time: time.Now().Add(-time.Minute / 2)}
	id := "test-id"
	userID := "test-user-id"
	sessionID := "test-session-id"
	viewID := "test-view-id"
	tabID := "test-tab-id"
	return map[string]pingDataSet{
		"Given no stats": {
			origin: nil,
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  "",
				TabID:   "",
				Visible: false,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    now,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{},
				LastVisiblePing: now,
				VisibleDuration: 0,
				TabDuration:     map[string]interface{}{},
			},
		},
		"Given start stats and path without tab id": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        minuteAgo,
				LastVisiblePath: []string{},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 0,
				TabDuration:     map[string]interface{}{},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   "",
				Visible: true,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{viewID},
				LastVisiblePing: now,
				VisibleDuration: 60,
				TabDuration: map[string]interface{}{
					viewID: 60,
				},
			},
		},
		"Given start stats and path with tab id": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        minuteAgo,
				LastVisiblePath: []string{},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 0,
				TabDuration:     map[string]interface{}{},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   tabID,
				Visible: true,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{viewID, tabID},
				LastVisiblePing: now,
				VisibleDuration: 60,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID: 60,
					},
				},
			},
		},
		"Given last visible ping stats and path without tab id": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        halfminuteAgo,
				LastVisiblePath: []string{viewID},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 40,
				TabDuration: map[string]interface{}{
					viewID:         40,
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   "",
				Visible: true,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{viewID},
				LastVisiblePing: now,
				VisibleDuration: 70,
				TabDuration: map[string]interface{}{
					viewID:         100,
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
		},
		"Given last visible ping stats and path with tab id": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        halfminuteAgo,
				LastVisiblePath: []string{viewID, tabID},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 40,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID:        40,
						"anotherTab": 15,
					},
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   tabID,
				Visible: true,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{viewID, tabID},
				LastVisiblePing: now,
				VisibleDuration: 70,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID:        100,
						"anotherTab": 15,
					},
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
		},
		"Given not last visible ping stats and path without tab id": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        halfminuteAgo,
				LastVisiblePath: []string{"anotherView1"},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 40,
				TabDuration: map[string]interface{}{
					viewID:         40,
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   "",
				Visible: true,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{viewID},
				LastVisiblePing: now,
				VisibleDuration: 70,
				TabDuration: map[string]interface{}{
					viewID:         70,
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
		},
		"Given not last visible ping stats and path with tab id": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        halfminuteAgo,
				LastVisiblePath: []string{"anotherView1"},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 40,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID:        40,
						"anotherTab": 15,
					},
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   tabID,
				Visible: true,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{viewID, tabID},
				LastVisiblePing: now,
				VisibleDuration: 70,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID:        70,
						"anotherTab": 15,
					},
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
		},
		"Given not visible path": {
			origin: &Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        halfminuteAgo,
				LastVisiblePath: []string{"anotherView1"},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 40,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID:        40,
						"anotherTab": 15,
					},
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
			session: SessionData{
				SessionID: sessionID,
				UserID:    userID,
			},
			path: PathData{
				ViewID:  viewID,
				TabID:   "",
				Visible: false,
			},
			expected: Stats{
				ID:              id,
				SessionID:       sessionID,
				SessionStart:    minuteAgo,
				UserID:          userID,
				LastPing:        now,
				LastVisiblePath: []string{"anotherView1"},
				LastVisiblePing: minuteAgo,
				VisibleDuration: 40,
				TabDuration: map[string]interface{}{
					viewID: map[string]interface{}{
						tabID:        40,
						"anotherTab": 15,
					},
					"anotherView1": 30,
					"anotherView2": map[string]interface{}{
						"anotherTab2": 10,
					},
				},
			},
		},
	}
}

func assertStats(t *testing.T, expected Stats, s Stats, testName string) {
	errorf := func(field string, expected, value interface{}) {
		if testName == "" {
			t.Errorf("expected %s: %v but got %v", field, expected, value)
		} else {
			t.Errorf("%s: expected %s: %v but got %v", testName, field, expected, value)
		}
	}

	if s.ID == "" {
		errorf("Stats.ID", "", s.ID)
	}

	if s.SessionID != expected.SessionID {
		errorf("Stats.ID", expected.SessionID, s.SessionID)
	}

	if s.UserID != expected.UserID {
		errorf("Stats.UserID", expected.UserID, s.UserID)
	}

	if !isEqualDuration(s.VisibleDuration, expected.VisibleDuration) {
		errorf("Stats.VisibleDuration", expected.VisibleDuration, s.VisibleDuration)
	}

	if !isEqualTabDuration(s.TabDuration, expected.TabDuration) {
		errorf("Stats.TabDuration", expected.TabDuration, s.TabDuration)
	}

	if !reflect.DeepEqual(s.LastVisiblePath, expected.LastVisiblePath) {
		errorf("Stats.LastVisiblePath", expected.LastVisiblePath, s.LastVisiblePath)
	}

	if !isEqualTime(s.SessionStart, expected.SessionStart) {
		errorf("Stats.SessionStart", expected.SessionStart.Unix(), s.SessionStart.Unix())
	}

	if !isEqualTime(s.LastPing, expected.LastPing) {
		errorf("Stats.LastPing", expected.LastPing.Unix(), s.LastPing.Unix())
	}

	if !isEqualTime(s.LastVisiblePing, expected.LastVisiblePing) {
		errorf("Stats.LastVisiblePing", expected.LastVisiblePing.Unix(), s.LastVisiblePing.Unix())
	}
}

func isEqualTime(l, r types.CpsTime) bool {
	return math.Abs(float64(l.Unix()-r.Unix())) < 2
}

func isEqualDuration(l, r int64) bool {
	return math.Abs(float64(l-r)) < 2
}

func isEqualTabDuration(l, r map[string]interface{}) bool {
	if len(l) != len(r) {
		return false
	}

	for k := range l {
		lv := l[k]
		rv, ok := r[k]
		if !ok {
			return false
		}

		if lvi := toInt(lv); lvi != 0 {
			if rvi := toInt(rv); rvi != 0 {
				if isEqualDuration(lvi, rvi) {
					continue
				}
			}

			return false
		}

		if lvm, ok := lv.(map[string]interface{}); ok {
			if rvm, ok := rv.(map[string]interface{}); ok {
				if len(lvm) != len(rvm) {
					return false
				}

				for dk := range lvm {
					ld := lvm[dk]
					rd, ok := rvm[dk]
					if !ok {
						return false
					}

					if ldi := toInt(ld); ldi != 0 {
						if rdi := toInt(rd); rdi != 0 {
							if isEqualDuration(ldi, rdi) {
								continue
							}
						}
					}

					return false
				}

				continue
			}

			return false
		}

		return false
	}

	return true
}

func mockCursor(ctrl *gomock.Controller, data *Stats) libmongo.Cursor {
	mockCursor := mock_mongo.NewMockCursor(ctrl)

	if data != nil {
		mockCursor.EXPECT().Next(gomock.Any()).Return(true)
		mockCursor.
			EXPECT().
			Decode(gomock.Any()).
			Do(func(val interface{}) {
				if u, ok := val.(*Stats); ok {
					*u = *data
				}
			}).
			Return(nil)
	} else {
		mockCursor.EXPECT().Next(gomock.Any()).Return(false)
	}

	mockCursor.EXPECT().Close(gomock.Any())

	return mockCursor
}
