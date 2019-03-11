package schedule

import (
	"testing"
	"time"

	com "github.com/att/deadline/common"

	"github.com/stretchr/testify/assert"
)

var (
	simpleEventNode = &NodeInstance{
		value: &EventNode{
			okTo:    secondEventNode,
			errorTo: endNode,
			constraints: com.EventConstraints{
				ReceiveBy: time.Now().Unix() - int64(time.Hour.Seconds()*1000),
			},
		},
	}

	secondEventNode = &NodeInstance{
		value: &EventNode{
			okTo:    nil,
			errorTo: endNode,
			constraints: com.EventConstraints{
				ReceiveBy: time.Now().Unix() - int64(time.Hour.Seconds()*1000),
			},
		},
	}

	startNode = &StartNode{
		to: simpleEventNode,
	}

	endNode = &NodeInstance{
		value: &EndNode{},
	}
)

func TestStartNode(test *testing.T) {

	next, err := startNode.Next()
	assert.Nil(test, err, "")
	assert.Equal(test, len(next), 1)
	assert.Equal(test, next[0], simpleEventNode)
}

func TestEventOKTo(test *testing.T) {
	e := com.Event{
		ReceivedAt: time.Now().Unix() - int64(time.Hour.Seconds()*2*1000), // 2 hrs ago
	}

	if node, ok := simpleEventNode.value.(*EventNode); !ok {
		test.FailNow()
	} else {

		node.AddEvent(&e)
		next, ctx := node.Next()

		assert.NotNil(test, ctx, "")
		assert.Equal(test, true, ctx.Successful)
		assert.Equal(test, len(next), 1)
		assert.Equal(test, next[0], secondEventNode)
	}

}

func TestEventErrorTo(test *testing.T) {
	e := com.Event{
		ReceivedAt: time.Now().Unix(),
	}

	if node, ok := secondEventNode.value.(*EventNode); !ok {
		test.FailNow()
	} else {

		node.AddEvent(&e)
		next, ctx := node.Next()

		assert.Equal(test, len(next), 1)
		assert.Equal(test, next[0], endNode)
		assert.NotNil(test, ctx)
		assert.Equal(test, com.LateEvent, ctx.FailureContext.Reason)
	}
}
