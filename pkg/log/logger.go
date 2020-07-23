/*
MIT License

Copyright (c) 2017 Kolide

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// A version of this logger was originally created at Kolide. Borrowed and modified from
// https://github.com/kolide/kit/tree/8cde91971ef08747188adf1f0673c2565598aa73/logutil

package log

import (
	"io"
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Exit logs an an error and message, terminating the program.
// Exit should only be used in package main where it's okay to terminate the program.
// When err is nil, the program exits successfully.
func Exit(logger Logger, err error, msg string) {
	if err == nil {
		Debug(logger).Log("err", err, "msg", msg)
		return
	}
	Info(logger).Log("err", err, "msg", msg)
	os.Exit(1)
}

// Option sets configuration for the logger.
type Option func(*config)

// SwapSignal specifies a os.Signal to swap between Debug and Info levels in real time.
// The default is SIGUSR2.
func SwapSignal(sig os.Signal) Option {
	return func(c *config) {
		c.sig = sig
	}
}

// JSON configures the logger format to JSON.
// The default is logfmt (https://brandur.org/logfmt)
func JSON() Option {
	return func(c *config) {
		c.format = log.NewJSONLogger
	}
}

// StartDebug creates a logger configured to allow debug level logs from the start.
func StartDebug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Output configures the log output. Stderr is default.
func Output(w io.Writer) Option {
	return func(c *config) {
		c.w = w
	}
}

type config struct {
	w      io.Writer
	format func(io.Writer) log.Logger
	sig    os.Signal
	debug  bool
}

// New creates a Logger.
func New(opts ...Option) *log.SwapLogger {
	c := config{
		w:      os.Stderr,
		format: log.NewLogfmtLogger,
		sig:    defaultSwapSignal,
	}

	for _, optFn := range opts {
		optFn(&c)
	}

	base := c.format(log.NewSyncWriter(c.w))
	base = log.With(base, "ts", log.DefaultTimestampUTC)
	base = level.NewInjector(base, level.InfoValue())
	lev := level.AllowInfo()
	if c.debug {
		lev = level.AllowDebug()
	}

	var swapLogger log.SwapLogger
	swapLogger.Swap(level.NewFilter(base, lev))

	go c.swapLevelHandler(base, &swapLogger, c.debug)
	return &swapLogger
}

func (c *config) swapLevelHandler(base Logger, swapLogger *log.SwapLogger, debug bool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, c.sig)
	for {
		<-sigChan
		if debug {
			newLogger := level.NewFilter(base, level.AllowInfo())
			swapLogger.Swap(newLogger)
		} else {
			newLogger := level.NewFilter(base, level.AllowDebug())
			swapLogger.Swap(newLogger)
		}
		Info(swapLogger).Log("msg", "swapping level", "debug", !debug)
		debug = !debug
	}
}
