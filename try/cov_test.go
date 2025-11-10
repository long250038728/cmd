package try

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func FuzzX(f *testing.F) {
	f.Add(1)
	f.Fuzz(func(t *testing.T, s int) {
		t.Log(s)
		v := Add(s, s)
		if v == 2 {
			t.Error(v)
		}

	})
}

type LoggerConfig struct {
	bufferSize  int32
	waitTimeout time.Duration
	file        string
}
type Opt func(config *LoggerConfig)

func SetBufferSize(bufferSize int32) Opt {
	return func(config *LoggerConfig) {
		config.bufferSize = bufferSize
	}
}

func SetWaitTimeout(waitTimeout time.Duration) Opt {
	return func(config *LoggerConfig) {
		config.waitTimeout = waitTimeout
	}
}

func SetFile(file string) Opt {
	return func(config *LoggerConfig) {
		config.file = file
	}
}

type Logger struct {
	ctx    context.Context
	cancel context.CancelFunc

	buffer  []string
	ch      chan string
	writeCh chan []string
	locker  atomic.Bool

	config *LoggerConfig
}

func NewLogger(opts ...Opt) *Logger {
	ctx, cancel := context.WithCancel(context.Background())
	config := &LoggerConfig{
		bufferSize:  100,
		waitTimeout: time.Minute,
		file:        "./log.log",
	}

	for _, opt := range opts {
		opt(config)
	}

	return &Logger{
		ctx:     ctx,
		cancel:  cancel,
		ch:      make(chan string, config.bufferSize),
		buffer:  make([]string, 0, config.bufferSize),
		writeCh: make(chan []string, config.bufferSize),
		config:  config,
	}
}

// AddLog 写入日志，当队列满了丢弃数据
func (l *Logger) AddLog(ctx context.Context, data string) {
	select {
	case <-l.ctx.Done():
		return
	case <-ctx.Done():
		return
	case l.ch <- data:
	default: // 丢弃满批次
	}
}

func (l *Logger) Write() error {
	if !l.locker.CompareAndSwap(false, true) {
		return errors.New("logging already started")
	}

	f, err := os.OpenFile(l.config.file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	writer := bufio.NewWriter(f)

	go func() {
		timer := time.NewTimer(l.config.waitTimeout)
		defer timer.Stop()

		for {
			select {
			case <-l.ctx.Done():
				return
			case log := <-l.ch:
				l.buffer = append(l.buffer, log)
				if len(l.buffer) >= int(l.config.bufferSize) {
					select {
					case l.writeCh <- l.buffer:
					default: // 丢弃满批次
					}
					l.buffer = make([]string, 0, l.config.bufferSize)
				}
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(l.config.waitTimeout)
			case <-timer.C:
				if len(l.buffer) > 0 {
					select {
					case l.writeCh <- l.buffer:
					default: // 丢弃满批次
					}
					l.buffer = make([]string, 0, l.config.bufferSize)
				}
				timer.Reset(l.config.waitTimeout)
			}

		}
	}()

	go func() {
		defer func() {
			_ = writer.Flush()
			_ = f.Close()
		}()

		for {
			select {
			case <-l.ctx.Done():
				return
			case logs := <-l.writeCh:
				b, err := json.Marshal(logs)
				if err != nil {
					continue
				}
				_, _ = writer.Write(b)
				_, _ = writer.Write([]byte("\n"))
				_ = writer.Flush()
			}
		}
	}()

	return nil
}

func (l *Logger) Close() {
	l.cancel()
}
