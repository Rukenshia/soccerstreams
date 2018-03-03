package logging

import raven "github.com/getsentry/raven-go"

func CreatePacket(level raven.Severity, format string, args ...interface{}) *raven.Packet {
	p := &raven.Packet{}
	p.Interfaces = []raven.Interface{
		&raven.Message{
			Message: format,
			Params:  args,
		},
	}
	p.Message = format
	p.Level = level

	return p
}
