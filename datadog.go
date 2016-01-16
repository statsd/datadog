// Package datadog implements a simple dogstatsd client
//
// 			c, err := Dial(":5000")
// 			c.SetPrefix("myprogram")
// 			c.SetTags("env:stage", "program:myprogram")
// 			c.Incr("count")
//
package datadog

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

// maxBufSize is the default buffer size
// when that size is reached the packet is flushed.
const maxBufSize = 512

// Client represents a datadog client.
type Client struct {
	conn   net.Conn
	buf    *bufio.Writer
	prefix string
	tags   []string
	rand   func() float64
	sync.Mutex
}

// Dial connects to `addr` and returns a client.
func Dial(addr string) (*Client, error) {
	return DialSize(addr, maxBufSize)
}

// DialSize connects with the given buffer `size`.
// see https://git.io/vzC0D.
func DialSize(addr string, size int) (*Client, error) {
	c, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		buf:  bufio.NewWriterSize(c, size),
		rand: rand.Float64,
		conn: c,
	}, nil
}

// New returns a new Client with writer `w`, useful for testing.
func New(w io.Writer) *Client {
	return &Client{
		buf: bufio.NewWriterSize(w, 512),
	}
}

// SetPrefix sets global prefix `name`.
func (c *Client) SetPrefix(name string) {
	c.Lock()
	defer c.Unlock()

	if name[len(name)-1] != '.' {
		name += "."
	}

	c.prefix = name
}

// SetTags sets global tags `tags...`.
func (c *Client) SetTags(tags ...string) {
	c.Lock()
	defer c.Unlock()
	c.tags = tags
}

// Increment incremenets the given stat `name`
// with the given `count`, `rate` and `tags...`.
func (c *Client) Increment(name string, count int, rate float64, tags ...string) error {
	value := strconv.Itoa(count) + "|c"
	return c.send(name, value, rate, tags)
}

// IncrBy increments counter `name` by `n` with optional `tags`.
func (c *Client) IncrBy(name string, n int, tags ...string) error {
	return c.Increment(name, n, 1, tags...)
}

// DecrBy decrements counter `name` by `n` with optional `tags`.
func (c *Client) DecrBy(name string, n int, tags ...string) error {
	return c.Increment(name, -n, 1, tags...)
}

// Incr increments counter `name` with `tags`.
func (c *Client) Incr(name string, tags ...string) error {
	return c.IncrBy(name, 1, tags...)
}

// Decr decrements counter `name` with `tags`.
func (c *Client) Decr(name string, tags ...string) error {
	return c.DecrBy(name, 1, tags...)
}

// Gauge sets the metric `name` to `n` at a given time.
func (c *Client) Gauge(name string, n int, tags ...string) error {
	value := strconv.Itoa(n) + "|g"
	return c.send(name, value, 1, tags)
}

// Histogram measures the statistical distribution of a metric `name`
// with the given `v` value, `rate` and `tags`.
func (c *Client) Histogram(name string, v int, tags ...string) error {
	value := strconv.Itoa(v) + "|h"
	return c.send(name, value, 1, tags)
}

// Duration uses `Histogram()` to send the given `d` duration.
func (c *Client) Duration(name string, d time.Duration, tags ...string) error {
	return c.Histogram(name, int(d.Seconds()*1000), tags...)
}

// Unique records a unique occurence of events.
func (c *Client) Unique(name, value string, tags ...string) error {
	return c.send(name, value+"|s", 1, tags)
}

// Flush will flush the underlying buffer.
func (c *Client) Flush() error {
	return c.buf.Flush()
}

// Close will flush and close the connection.
func (c *Client) Close() error {
	var err error

	defer func() {
		if e := c.conn.Close(); err == nil {
			err = e
		}
	}()

	err = c.Flush()
	return err
}

// send sends the given `stat` with `format`, `rate`
// `tags` and the given `args...`.
func (c *Client) send(stat, value string, rate float64, tags []string) error {
	var buf bytes.Buffer

	c.Lock()
	defer c.Unlock()

	if c.buf.Buffered() > 0 {
		buf.WriteString("\n")
	}

	buf.WriteString(c.prefix)
	buf.WriteString(stat + ":" + value)

	// rate
	if rate < 1 {
		if c.rand() < rate {
			buf.WriteString("|@" + strconv.FormatFloat(rate, 'f', 1, 64))
		} else {
			return nil
		}
	}

	// tags
	if len(c.tags)+len(tags) > 0 {
		tags = append(c.tags, tags...)
		buf.WriteString("|#" + strings.Join(tags, ","))
	}

	// flush if there's no space left.
	if c.buf.Available() < buf.Len() {
		err := c.buf.Flush()
		if err != nil {
			return err
		}
	}

	_, err := c.buf.Write(buf.Bytes())
	return err
}
