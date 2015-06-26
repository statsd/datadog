# datadog
--
    import "github.com/statsd/datadog"

Package datadog provides a wrapper around the DataDog-specific statsd client
which augments the protocol to support tags.

## Usage

#### type Client

```go
type Client struct {
  DataDog *dogstatsd.Client
}
```

Client.

#### func  New

```go
func New(addr string) (*Client, error)
```
New statsd client.

#### func (*Client) SetPrefix

```go
func (c *Client) SetPrefix(prefix string)
```
Sets the DataDog namespace to the prefix provided.


#### func (*Client) Annotate

```go
func (c *Client) Annotate(name string, value string, args ...interface{}) error
```
Annotate sends an annotation.

#### func (*Client) Decr

```go
func (c *Client) Decr(name string, tags ...string) error
```
Decr decrements the counter for the given bucket by 1 at a rate of 1.

#### func (*Client) DecrBy

```go
func (c *Client) DecrBy(name string, value int, tags ...string) error
```
DecrBy decrements the counter for the given bucket by N at a rate of 1.

#### func (*Client) Decrement

```go
func (c *Client) Decrement(name string, count int, rate float64, tags ...string) error
```
Decrement decrements the counter for the given bucket.

#### func (*Client) Duration

```go
func (c *Client) Duration(name string, duration time.Duration, tags ...string) error
```
Duration records time spent for the given bucket with time.Duration.

#### func (*Client) Flush

```go
func (c *Client) Flush() error
```
Flush flushes writes any buffered data to the network.

Unsupported.

#### func (*Client) Gauge

```go
func (c *Client) Gauge(name string, value int, tags ...string) error
```
Gauge records arbitrary values for the given bucket.

#### func (*Client) Histogram

```go
func (c *Client) Histogram(name string, value int, tags ...string) error
```
Histogram is an alias of .Duration() until the statsd protocol figures its shit
out.

#### func (*Client) Incr

```go
func (c *Client) Incr(name string, tags ...string) error
```
Incr increments the counter for the given bucket by 1 at a rate of 1.

#### func (*Client) IncrBy

```go
func (c *Client) IncrBy(name string, n int, tags ...string) error
```
IncrBy increments the counter for the given bucket by N at a rate of 1.

#### func (*Client) Increment

```go
func (c *Client) Increment(name string, count int, rate float64, tags ...string) error
```
Increment increments the counter for the given bucket.

#### func (*Client) Unique

```go
func (c *Client) Unique(name string, value int, rate float64) error
```
Unique records unique occurences of events.

Unsupported.

