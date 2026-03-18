package infra

import (
	"app/internal/crosscutting/timer"
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type basedOnPlaceTimeGenerator struct {
	tz time.Location
}

func (t basedOnPlaceTimeGenerator) Now() time.Time {
	return time.Now().In(&t.tz)
}

func (t basedOnPlaceTimeGenerator) GetTz() time.Location {
	return t.tz
}

func (t basedOnPlaceTimeGenerator) NowTs() timer.UnixTimestamp {
	return timer.UnixTimestamp(time.Now().In(&t.tz).Unix())
}

func (t basedOnPlaceTimeGenerator) NowTsMills() timer.UnixTimestampMillis {
	return timer.UnixTimestampMillis(time.Now().In(&t.tz).UnixMilli())
}

func (t basedOnPlaceTimeGenerator) TimeFromInt64(i int64) time.Time {
	return time.Unix(i, 0).In(&t.tz)
}

func NewTimeGenerator(tz time.Location) timer.TimeGenerator {
	return &basedOnPlaceTimeGenerator{tz: tz}
}

func NewJstTimeGenerator() timer.TimeGenerator {
	return &basedOnPlaceTimeGenerator{tz: *jst}
}
