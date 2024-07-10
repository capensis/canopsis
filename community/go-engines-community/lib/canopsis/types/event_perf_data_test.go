package types_test

import (
	"testing"

	"git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/canopsis/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestEvent_GetPerfData(t *testing.T) {
	dataSet := []struct {
		Input    string
		Expected []types.PerfData
	}{
		{
			Input: "cpu=20%;80;90;0;100",
			Expected: []types.PerfData{
				{
					Name:  "cpu_%",
					Value: 20,
					Unit:  "%",
				},
			},
		},
		{
			Input: "'memory'=2GB;1;2.5;0;3",
			Expected: []types.PerfData{
				{
					Name:  "memory_GB",
					Value: 2,
					Unit:  "GB",
				},
			},
		},
		{
			Input: "'disk i/o'=1.5TB",
			Expected: []types.PerfData{
				{
					Name:  "disk i/o_TB",
					Value: 1.5,
					Unit:  "TB",
				},
			},
		},
		{
			Input: "temp=57.5",
			Expected: []types.PerfData{
				{
					Name:  "temp",
					Value: 57.5,
					Unit:  "",
				},
			},
		},
		{
			Input:    "temp=57.",
			Expected: nil,
		},
		{
			Input: "'temp ''C'''=57.5",
			Expected: []types.PerfData{
				{
					Name:  "temp 'C'",
					Value: 57.5,
					Unit:  "",
				},
			},
		},
		{
			Input: "cpu=20%;80;90;0;100 'memory'=2GB;1;2.5;0;3 'disk i/o'=1.5TB temp=57.5",
			Expected: []types.PerfData{
				{
					Name:  "cpu_%",
					Value: 20,
					Unit:  "%",
				},
				{
					Name:  "memory_GB",
					Value: 2,
					Unit:  "GB",
				},
				{
					Name:  "disk i/o_TB",
					Value: 1.5,
					Unit:  "TB",
				},
				{
					Name:  "temp",
					Value: 57.5,
					Unit:  "",
				},
			},
		},
		{
			Input: "",
		},
		{
			Input: "=",
		},
		{
			Input: "temp=57.5 cpu=",
		},
		{
			Input: "temp=57.5 cpu=%",
		},
		{
			Input: "temp=57.5 'cpu=20%",
		},
		{
			Input: "temp=57.5 'cpu''=20%",
		},
		{
			Input: "cp''u=20%",
			Expected: []types.PerfData{
				{
					Name:  "cp'u_%",
					Value: 20,
					Unit:  "%",
				},
			},
		},
		{
			Input: "temp=57.5 cpu =20%",
		},
		{
			Input: "temp=U",
		},
		{
			Input: "temp=U cpu=20%",
			Expected: []types.PerfData{
				{
					Name:  "cpu_%",
					Value: 20,
					Unit:  "%",
				},
			},
		},
		{
			Input: "memory=UGB cpu=20%",
			Expected: []types.PerfData{
				{
					Name:  "cpu_%",
					Value: 20,
					Unit:  "%",
				},
			},
		},
		{
			Input: "cpu=20.2.2%",
		},
		{
			Input: "cpu=20%%",
		},
		{
			Input: "cpu=20&",
		},
		{
			Input: "temp=20°C",
			Expected: []types.PerfData{
				{
					Name:  "temp_°C",
					Value: 20,
					Unit:  "°C",
				},
			},
		},
	}

	event := types.Event{}
	for _, data := range dataSet {
		event.PerfData = data.Input
		result := event.GetPerfData()
		if diff := pretty.Compare(data.Expected, result); diff != "" {
			t.Errorf("%q: %s", data.Input, diff)
		}
	}
}
