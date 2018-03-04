package logging

import (
	"testing"

	raven "github.com/getsentry/raven-go"
	"github.com/stretchr/testify/assert"
)

func TestCreatePacketSetsMessage(t *testing.T) {
	level := raven.DEBUG
	format := "abcd"
	args := []interface{}{
		"foo",
	}

	p := CreatePacket(level, format, args...)

	assert.Equal(t, format, p.Message)
}

func TestCreatePacketSetsMessageInterface(t *testing.T) {
	level := raven.DEBUG
	format := "abcd"
	args := []interface{}{
		"foo",
	}

	p := CreatePacket(level, format, args...)

	assert.Equal(t, 1, len(p.Interfaces))

	if m, ok := p.Interfaces[0].(*raven.Message); ok {
		assert.Equal(t, format, m.Message)
		assert.Equal(t, args, m.Params)
	} else {
		t.Error("Packet interface is not of type raven.Message")
	}
}
