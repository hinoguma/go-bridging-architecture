package timer

import (
	"time"
)

var globalTimeGenerator TimeGenerator = NewStdTimeGenerator()

func SetGlobalTimeGenerator(tg TimeGenerator) {
	globalTimeGenerator = tg
}

type TimeGenerator interface {
	GetTz() time.Location
	Now() time.Time
	NowTs() UnixTimestamp
	NowTsMills() UnixTimestampMillis
	TimeFromInt64(i int64) time.Time
}

func NowTs() UnixTimestamp {
	return globalTimeGenerator.NowTs()
}

func Now() time.Time {
	return globalTimeGenerator.Now()
}

func NowTsMills() UnixTimestampMillis {
	return globalTimeGenerator.NowTsMills()
}

func NowYmd_Hms_ms() Ymd_Hms_ms {
	return NowTsMills().Ymd_Hms_ms()
}

func TimeFromInt64(i int64) time.Time {
	return globalTimeGenerator.TimeFromInt64(i)
}

type StdTimeGenerator struct {
	tz time.Location
}

func NewStdTimeGenerator() TimeGenerator {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return &StdTimeGenerator{tz: *jst}
}

func (t StdTimeGenerator) GetTz() time.Location {
	return t.tz
}

func (t StdTimeGenerator) Now() time.Time {
	return time.Now().In(&t.tz)
}

func (t StdTimeGenerator) NowTs() UnixTimestamp {
	return UnixTimestamp(time.Now().In(&t.tz).Unix())
}

func (t StdTimeGenerator) NowTsMills() UnixTimestampMillis {
	return UnixTimestampMillis(time.Now().In(&t.tz).UnixMilli())
}

func (t StdTimeGenerator) TimeFromInt64(i int64) time.Time {
	return time.Unix(i, 0).In(&t.tz)
}
