package datadog

import (
	"io/ioutil"
	"testing"
	"time"
)

func BenchmarkGauge(b *testing.B) {
	c := New(ioutil.Discard)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			c.Gauge("free.memory", 512, "program:counters")
		}
	})
}

func BenchmarkIncr(b *testing.B) {
	c := New(ioutil.Discard)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			c.Incr("counter", "program:counters")
		}
	})
}

func BenchmarkDecr(b *testing.B) {
	c := New(ioutil.Discard)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			c.Decr("counter", "program:counters")
		}
	})
}

func BenchmarkHistogram(b *testing.B) {
	c := New(ioutil.Discard)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			c.Histogram("response.size", 512, "program:server")
		}
	})
}

func BenchmarkDuration(b *testing.B) {
	c := New(ioutil.Discard)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			c.Duration("response.time", time.Second, "program:counters")
		}
	})
}

func BenchmarkUnique(b *testing.B) {
	c := New(ioutil.Discard)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			c.Unique("users", "1", "program:counters")
		}
	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
