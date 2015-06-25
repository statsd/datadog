// Package datadog provides a wrapper around the DataDog-specific
// statsd client which augments the protocol to support tags.
package datadog

import (
	"time"

	"github.com/ooyala/go-dogstatsd"
)

// Client.
type Client struct {
	DataDog *dogstatsd.Client
}

// New statsd client.
func New(addr string) (*Client, error) {
	c, err := dogstatsd.New(addr)

	if err != nil {
		return nil, err
	}

	return &Client{
		DataDog: c,
	}, nil
}

// Increment increments the counter for the given bucket.
func (c *Client) Increment(name string, count int, rate float64, tags ...string) error {
	return c.DataDog.Count(name, int64(count), tags, rate)
}

// Incr increments the counter for the given bucket by 1 at a rate of 1.
func (c *Client) Incr(name string, tags ...string) error {
	return c.Increment(name, 1, 1, tags...)
}

// IncrBy increments the counter for the given bucket by N at a rate of 1.
func (c *Client) IncrBy(name string, n int, tags ...string) error {
	return c.Increment(name, n, 1, tags...)
}

// Decrement decrements the counter for the given bucket.
func (c *Client) Decrement(name string, count int, rate float64, tags ...string) error {
	return c.Increment(name, -count, rate, tags...)
}

// Decr decrements the counter for the given bucket by 1 at a rate of 1.
func (c *Client) Decr(name string, tags ...string) error {
	return c.Increment(name, -1, 1, tags...)
}

// DecrBy decrements the counter for the given bucket by N at a rate of 1.
func (c *Client) DecrBy(name string, value int, tags ...string) error {
	return c.Increment(name, -value, 1, tags...)
}

// Duration records time spent for the given bucket with time.Duration.
func (c *Client) Duration(name string, duration time.Duration, tags ...string) error {
	return c.Histogram(name, millisecond(duration), tags...)
}

// Histogram is an alias of .Duration() until the statsd protocol figures its shit out.
func (c *Client) Histogram(name string, value int, tags ...string) error {
	return c.DataDog.Histogram(name, float64(value), tags, 1)
}

// Gauge records arbitrary values for the given bucket.
func (c *Client) Gauge(name string, value int, tags ...string) error {
	return c.DataDog.Gauge(name, float64(value), tags, 1)
}

// Annotate sends an annotation.
func (c *Client) Annotate(name string, value string, args ...interface{}) error {
	return c.DataDog.Event(name, value, nil)
}

// Unique records unique occurences of events.
//
// Unsupported.
func (c *Client) Unique(name string, value int, rate float64) error {
	return nil // Unsupported
}

// Flush flushes writes any buffered data to the network.
//
// Unsupported.
func (c *Client) Flush() error {
	return nil
}

func millisecond(d time.Duration) int {
	return int(d.Seconds() * 1000)
}
