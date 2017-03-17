package bc

import "context"

type (
	TimeRangeBody struct {
		MinTimeMS, MaxTimeMS uint64
		ExtHash              Hash
	}

	// TimeRange denotes a time range. It satisfies the Entry interface.
	TimeRange struct {
		TimeRangeBody
	}
)

func (TimeRange) Type() string          { return "timerange1" }
func (tr *TimeRange) Body() interface{} { return tr.TimeRangeBody }

func (TimeRange) Ordinal() int { return -1 }

// NewTimeRange creates a new TimeRange.
func NewTimeRange(minTimeMS, maxTimeMS uint64) *TimeRange {
	tr := new(TimeRange)
	tr.MinTimeMS = minTimeMS
	tr.MaxTimeMS = maxTimeMS
	return tr
}

func (tr *TimeRange) CheckValid(_ context.Context) error {
	// xxx check MinTimeMS <= MaxTimeMS?
	// xxx check ExtHash is all zeroes?
	return nil
}
