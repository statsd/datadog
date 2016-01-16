[![Build Status](https://semaphoreci.com/api/v1/projects/78ca1ac8-9aec-4b33-bc9e-a4ff4fc5b70e/663600/badge.svg)](https://semaphoreci.com/yields/datadog)

# datadog
--
    import "github.com/statsd/datadog"

Package datadog implements a simple dogstatsd client

    c, err := Dial(":5000")
    c.SetPrefix("myprogram")
    c.SetTags("env:stage", "program:myprogram")
    c.Incr("count")

## Usage

#### type Client

```go
type Client struct {
	sync.Mutex
}
```

Client represents a datadog client.

#### func  Dial

```go
func Dial(addr string) (*Client, error)
```
Dial connects to `addr` and returns a client.

#### func  DialSize

```go
func DialSize(addr string, size int) (*Client, error)
```
DialSize connects with the given buffer `size`. see https://git.io/vzC0D.

#### func  New

```go
func New(w io.Writer) *Client
```
New returns a new Client with writer `w`, useful for testing.

#### func (*Client) Close

```go
func (c *Client) Close() error
```
Close will flush and close the connection.

#### func (*Client) Decr

```go
func (c *Client) Decr(name string, tags ...string) error
```
Decr decrements counter `name` with `tags`.

#### func (*Client) DecrBy

```go
func (c *Client) DecrBy(name string, n int, tags ...string) error
```
DecrBy decrements counter `name` by `n` with optional `tags`.

#### func (*Client) Duration

```go
func (c *Client) Duration(name string, d time.Duration, tags ...string) error
```
Duration uses `Histogram()` to send the given `d` duration.

#### func (*Client) Flush

```go
func (c *Client) Flush() error
```
Flush will flush the underlying buffer.

#### func (*Client) Gauge

```go
func (c *Client) Gauge(name string, n int, tags ...string) error
```
Gauge sets the metric `name` to `n` at a given time.

#### func (*Client) Histogram

```go
func (c *Client) Histogram(name string, v int, tags ...string) error
```
Histogram measures the statistical distribution of a metric `name` with the
given `v` value, `rate` and `tags`.

#### func (*Client) Incr

```go
func (c *Client) Incr(name string, tags ...string) error
```
Incr increments counter `name` with `tags`.

#### func (*Client) IncrBy

```go
func (c *Client) IncrBy(name string, n int, tags ...string) error
```
IncrBy increments counter `name` by `n` with optional `tags`.

#### func (*Client) Increment

```go
func (c *Client) Increment(name string, count int, rate float64, tags ...string) error
```
Increment incremenets the given stat `name` with the given `count`, `rate` and
`tags...`.

#### func (*Client) SetPrefix

```go
func (c *Client) SetPrefix(name string)
```
SetPrefix sets global prefix `name`.

#### func (*Client) SetTags

```go
func (c *Client) SetTags(tags ...string)
```
SetTags sets global tags `tags...`.

#### func (*Client) Unique

```go
func (c *Client) Unique(name, value string, tags ...string) error
```
Unique records a unique occurence of events.
