package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAliases(test *testing.T) {
	hourly, err := timingToDuration("hourly")
	assert.Nil(test, err, "")
	assert.True(test, hourly == time.Hour, "")

	daily, err := timingToDuration("daily")
	assert.Nil(test, err, "")
	assert.True(test, daily == time.Hour*24, "")

	weekly, err := timingToDuration("weekly")
	assert.Nil(test, err, "")
	assert.True(test, weekly == time.Hour*24*7, "")
}

func TestTiming(test *testing.T) {
	hrs3, err := timingToDuration("3h")
	assert.Nil(test, err, "")
	assert.True(test, hrs3 == time.Hour*3, "")

	hrs7min30, err := timingToDuration("7h30m")
	assert.Nil(test, err, "")
	assert.True(test, hrs7min30 == time.Hour*7+time.Minute*30, "")
}

func TestNormalizeGood(test *testing.T) {
	now := time.Now()
	t, err := normailizeTime("2010-03-14T15:38:05+00:00", time.Hour*3)
	assert.Nil(test, err, "")

	//t, err := time.Parse(ExpectedTimeLayout, timeString)
	assert.Nil(test, err, "")
	assert.True(test, now.Year() == t.Year(), "")
	assert.True(test, now.Month() == t.Month(), "")
	assert.True(test, t.Unix() < now.Unix(), "")
	assert.True(test, t.Add(time.Hour*3).Unix() > now.Unix(), "")
}
