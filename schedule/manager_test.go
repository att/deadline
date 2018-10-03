package schedule

import (
	"testing"
	"time"

	"github.com/att/deadline/config"
	"github.com/stretchr/testify/assert"
)

var (
	testManager = GetManagerInstance(&config.DefaultConfig)
)

func TestAliases(test *testing.T) {
	hourly, err := timingToDuration("hourly")
	assert.Nil(test, err)
	assert.True(test, hourly == time.Hour)

	daily, err := timingToDuration("daily")
	assert.Nil(test, err)
	assert.True(test, daily == time.Hour*24)

	weekly, err := timingToDuration("weekly")
	assert.Nil(test, err)
	assert.True(test, weekly == time.Hour*24*7)
}

func TestTiming(test *testing.T) {
	hrs3, err := timingToDuration("3h")
	assert.Nil(test, err)
	assert.True(test, hrs3 == time.Hour*3)

	hrs7min30, err := timingToDuration("7h30m")
	assert.Nil(test, err)
	assert.True(test, hrs7min30 == time.Hour*7+time.Minute*30)
}

// func TestNormalizeGood(test *testing.T) {
// 	now := time.Now()
// 	t, _, err := normailizeTime("2010-03-14T12:38:05+00:00", time.Hour*12)
// 	assert.Nil(test, err)

// 	assert.Nil(test, err)
// 	// assert.True(test, now.Year() == t.Year())
// 	// assert.True(test, now.Month() == t.Month())
// 	assert.True(test, t.Unix() < now.Unix())
// 	assert.True(test, t.Add(time.Hour*3).Unix() > now.Unix())
// }

func TestManagerInitSched(test *testing.T) {
	yr, m, d := time.Now().In(time.UTC).Date()
	past := time.Date(2018, 5, 8, 0, 0, 0, 0, time.UTC)
	midnight := time.Date(yr, m, d, 0, 0, 0, 0, time.UTC)

	simpleBlueprint.StartsAt = past.Format(time.RFC3339)
	err := testManager.AddSchedule(*simpleBlueprint)

	assert.Nil(test, err)

	simpleSched := testManager.GetSchedule(simpleBlueprint.Name)
	assert.NotNil(test, simpleSched)

	// assert the first node exists, and is correctly time shifted
	firstNode := simpleSched.nodes["firstEvent"]
	assert.NotNil(test, firstNode)
	eventNode, ok := firstNode.value.(*EventNode)
	assert.True(test, ok)
	assert.NotNil(test, eventNode)

	assert.Equal(test, midnight.Add(4*time.Hour).Unix(), eventNode.constraints.ReceiveBy)
	assert.Equal(test, midnight, simpleSched.StartTime)

	// assert the second node exists, and is correctly time shifted
	secondNode := simpleSched.nodes["secondEvent"]
	assert.NotNil(test, firstNode)
	eventNode, ok = secondNode.value.(*EventNode)
	assert.True(test, ok)
	assert.NotNil(test, eventNode)

	assert.Equal(test, midnight.Add(12*time.Hour).Unix(), eventNode.constraints.ReceiveBy)
}
