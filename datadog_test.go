package datadog

import (
	"bufio"
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	go func() {
		a, err := net.ResolveUDPAddr("udp", ":5000")
		check(err)

		c, err := net.ListenUDP("udp", a)
		check(err)

		buf := make([]byte, 1024)

		for {
			_, _, err := c.ReadFromUDP(buf)
			check(err)
		}
	}()
}

func TestWrites(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.Gauge("memory", 512)
	client.Incr("count")
	client.Decr("count")
	client.Flush()

	assert.Equal(t, "\n"+buf.String(), `
memory:512|g
count:1|c
count:-1|c`)
}

func TestWritesWithTags(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.Gauge("memory", 512, "tag:a")
	client.Incr("count", "tag:b")
	client.Decr("count", "tag:c")
	client.Flush()

	assert.Equal(t, "\n"+buf.String(), `
memory:512|g|#tag:a
count:1|c|#tag:b
count:-1|c|#tag:c`)
}

func TestWritesGlobalTags(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.SetTags("global:tag")
	client.Gauge("memory", 512, "tag:a")
	client.Incr("count", "tag:b")
	client.Decr("count", "tag:c")
	client.Flush()

	assert.Equal(t, "\n"+buf.String(), `
memory:512|g|#global:tag,tag:a
count:1|c|#global:tag,tag:b
count:-1|c|#global:tag,tag:c`)
}

func TestPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.SetPrefix("myprogram")
	client.Gauge("memory", 512)
	client.Incr("count")
	client.Decr("count")
	client.Flush()

	assert.Equal(t, "\n"+buf.String(), `
myprogram.memory:512|g
myprogram.count:1|c
myprogram.count:-1|c`)
}

func TestHistogram(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.Histogram("size", 512)
	client.Flush()

	assert.Equal(t, buf.String(), `size:512|h`)
}

func TestDuration(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.Duration("duration", time.Second)
	client.Flush()

	assert.Equal(t, buf.String(), `duration:1000|h`)
}

func TestUnique(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.Unique("users", "68b260d7c0")
	client.Flush()

	assert.Equal(t, buf.String(), `users:68b260d7c0|s`)
}

func TestRate(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.rand = func() float64 { return .4 }
	client.SetPrefix("program")
	client.SetTags("env:stage")
	client.Increment("count", 1, 0.5)
	client.Flush()

	assert.Equal(t, buf.String(), `program.count:1|c|@0.5|#env:stage`)
}

func TestRateIgnored(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)

	client.rand = func() float64 { return .6 }
	client.Increment("count", 1, 0.5)
	client.Flush()

	assert.Equal(t, buf.String(), ``)
}

func TestFlushWrite(t *testing.T) {
	buf := new(bytes.Buffer)
	client := New(buf)
	client.buf = bufio.NewWriterSize(buf, 20)

	client.SetPrefix("program")
	client.SetTags("env:foo")
	client.Incr("counter-a")

	assert.Equal(t, buf.String(), `program.counter-a:1|c|#env:foo`)
}

func TestDialClose(t *testing.T) {
	c, err := Dial(":5000")
	assert.Equal(t, err, nil)
	assert.Equal(t, c.Close(), nil)
}
