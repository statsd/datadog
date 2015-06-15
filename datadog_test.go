package datadog

import (
	"github.com/statsd/client-interface"
)

var _ statsd.Client = (*Client)(nil)
