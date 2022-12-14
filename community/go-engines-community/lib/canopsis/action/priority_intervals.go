package action

import (
	"context"
	"sort"
	"sync"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const inf = 0

type priorityIntervals struct {
	intervals  map[int]int // key - lower bound, value - upper bound
	sortedKeys []int       // map may be not sorted, so we have an additional slice to keep keys sorted.
	mx         sync.Mutex
}

type PriorityIntervals interface {
	Recalculate(ctx context.Context, collection mongo.DbCollection) error
	Take(priority int)
	RestoreAndTake(oldPriority, newPriority int)
	Restore(priority int)
	GetIntervals() map[int]int
	GetSortedKeys() []int
	GetMinimal() int
}

func NewPriorityIntervals() PriorityIntervals {
	return &priorityIntervals{
		intervals:  map[int]int{1: inf},
		sortedKeys: []int{1},
		mx:         sync.Mutex{},
	}
}

func (pi *priorityIntervals) Recalculate(ctx context.Context, collection mongo.DbCollection) error {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	pi.intervals = map[int]int{1: inf}
	pi.sortedKeys = []int{1}

	var objPriority struct {
		Priority int `bson:"priority"`
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		err := cursor.Decode(&objPriority)
		if err != nil {
			// if err, don't save what was calculated, since it won't be valid
			pi.intervals = map[int]int{1: inf}
			pi.sortedKeys = []int{1}

			return err
		}

		pi.takePriority(objPriority.Priority)
	}

	return nil
}

func (pi *priorityIntervals) Take(priority int) {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	pi.takePriority(priority)
}

func (pi *priorityIntervals) RestoreAndTake(oldPriority, newPriority int) {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	pi.restorePriority(oldPriority)
	pi.takePriority(newPriority)
}

func (pi *priorityIntervals) Restore(priority int) {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	pi.restorePriority(priority)
}

func (pi *priorityIntervals) takePriority(priority int) {
	if priority < 1 {
		return
	}

	rightIntervalIndex := sort.Search(len(pi.sortedKeys), func(i int) bool {
		return pi.sortedKeys[i] >= priority
	})

	hitUpperBound, ok := pi.intervals[priority]
	if ok {
		/**
		if key exists in the map that means that we've hit the interval's lower bound, then we just increment the lower bound,
		if it's not an interval but a point(lowerBound = upperBound), then we just delete an interval from a map.
		*/

		delete(pi.intervals, priority)
		if priority == hitUpperBound {
			// if it's a point, remove from keys
			pi.sortedKeys = append(pi.sortedKeys[:rightIntervalIndex], pi.sortedKeys[rightIntervalIndex+1:]...)
			return
		}

		// if it's an interval, increase lower bound
		pi.intervals[priority+1] = hitUpperBound
		pi.sortedKeys[rightIntervalIndex] = priority + 1
		return
	}

	/**
	if key doesn't exist in the map that means that we've hit an interval or an upper bound
	so we need to divide the interval into two intervals with exclusion of the point that we've hit
	*/

	// Check if we've already take this point: only right interval case
	if rightIntervalIndex-1 < 0 && priority < pi.sortedKeys[rightIntervalIndex] {
		// can't take an already taken point
		return
	}

	leftLowerBound := pi.sortedKeys[rightIntervalIndex-1]
	leftUpperBound := pi.intervals[leftLowerBound]

	// Check if we've already take this point: between intervals case
	if leftUpperBound != inf && leftUpperBound < priority {
		// can't take an already taken point
		return
	}

	// if we've hit an upperbound, then we just decrement the upper bound
	if leftUpperBound == priority {
		pi.intervals[leftLowerBound]--
		return
	}

	// If we've hit inside the interval, then we need to divide it to two new intervals with exclusion of priority point(p)
	// ex: (1; p; 10) = (1; p-1) and (p+1; 10)

	pi.intervals[priority+1] = leftUpperBound
	pi.intervals[leftLowerBound] = priority - 1
	pi.sortedKeys = append(pi.sortedKeys[:rightIntervalIndex], append([]int{priority + 1}, pi.sortedKeys[rightIntervalIndex:]...)...)
}

func (pi *priorityIntervals) restorePriority(priority int) {
	if priority < 1 {
		return
	}

	_, ok := pi.intervals[priority]
	if ok {
		// can't restore an existing point
		return
	}

	// find the right interval
	rightIntervalIndex := sort.Search(len(pi.sortedKeys), func(i int) bool {
		return pi.sortedKeys[i] >= priority
	})
	leftIntervalIndex := rightIntervalIndex - 1

	if leftIntervalIndex < 0 {
		/**
		it means that left interval doesn't exist,
		so we have 2 options:
		1. Create a point. ex: {[1, 1]; [3, +inf)}
		2. Merge the point with the right interval, if distance between them = 1. ex: {[1, 1]; [2, +inf)} = [1, +inf)
		*/

		rightLowerBound := pi.sortedKeys[rightIntervalIndex]
		if rightLowerBound-1 == priority {
			pi.intervals[priority] = pi.intervals[rightLowerBound]
			delete(pi.intervals, rightLowerBound)
			pi.sortedKeys[rightIntervalIndex] = priority

			return
		}

		pi.intervals[priority] = priority
		pi.sortedKeys = append([]int{priority}, pi.sortedKeys[rightIntervalIndex:]...)

		return
	}

	leftLowerBound := pi.sortedKeys[leftIntervalIndex]
	leftUpperBound, ok := pi.intervals[leftLowerBound]
	if !ok {
		return
	}

	if leftUpperBound == inf || leftUpperBound > priority {
		//we're already inside the interval, just return
		return
	}

	rightLowerBound := pi.sortedKeys[rightIntervalIndex]
	rightUpperBound, ok := pi.intervals[rightLowerBound]
	if !ok {
		return
	}

	mergeLeft := leftUpperBound+1 == priority
	mergeRight := rightLowerBound-1 == priority

	// if distances between left and right intervals = 1, then merge both intervals
	//ex: p = 6, {[3, 5]; [p, p]; [7, +inf)} => [3, +inf)
	if mergeLeft && mergeRight {
		//merge two intervals
		pi.intervals[leftLowerBound] = rightUpperBound
		delete(pi.intervals, rightLowerBound)
		pi.sortedKeys = append(pi.sortedKeys[:rightIntervalIndex], pi.sortedKeys[rightIntervalIndex+1:]...)

		return
	}

	// if distances between left = 1, then merge with left
	//ex: p = 6, {[3, 5]; [p, p]; [8, +inf)} = {[3, p]; [8, +inf)}
	if mergeLeft {
		//increment left upper bound
		pi.intervals[leftLowerBound] = priority

		return
	}

	//otherwise with right
	//ex: p = 6, {[3, 4]; [p, p]; [7, +inf)} = {[3, 4], [p, +inf)}
	if mergeRight {
		//decrement right lower bound
		pi.intervals[priority] = rightUpperBound
		delete(pi.intervals, rightLowerBound)
		pi.sortedKeys[rightIntervalIndex] = priority

		return
	}

	//if both distances > 1 => don't merge, create a point
	pi.intervals[priority] = priority

	pi.sortedKeys = append(pi.sortedKeys, 0)
	copy(pi.sortedKeys[rightIntervalIndex+1:], pi.sortedKeys[rightIntervalIndex:])
	pi.sortedKeys[rightIntervalIndex] = priority
}

func (pi *priorityIntervals) GetIntervals() map[int]int {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	return pi.intervals
}

func (pi *priorityIntervals) GetSortedKeys() []int {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	return pi.sortedKeys
}

func (pi *priorityIntervals) GetMinimal() int {
	pi.mx.Lock()
	defer pi.mx.Unlock()

	return pi.sortedKeys[0]
}
