package p2p

import (
	"fmt"
	"strings"
	"time"
)

type Metrics struct {
	topic string
	addr  string

	tm time.Time

	handshake time.Duration
	resume    time.Duration
	read      time.Duration
	handle    time.Duration
	write     time.Duration

	stat []string
}

func newMetrics(addr string) (m *Metrics) {
	return &Metrics{
		addr: addr,
		tm:   time.Now(),

		handshake: -1,
		resume:    -1,
		read:      -1,
		handle:    -1,
		write:     -1,

		stat: []string{fmt.Sprintf("addr (%s)", addr)},
	}
}

func (m *Metrics) reset() {
	m.tm = time.Now()
}

func (m *Metrics) setTopic(topic string) {
	m.topic = topic
}

const statPattern = "%s (%s)"

func (m *Metrics) fixHandshake() {
	m.handshake = time.Since(m.tm)
	m.stat = append(m.stat, fmt.Sprintf(statPattern, "handshake", prepareValue(m.handshake)))
	m.reset()
}

func (m *Metrics) fixResume() {
	m.resume = time.Since(m.tm)
	m.stat = append(m.stat, fmt.Sprintf(statPattern, "resume", prepareValue(m.resume)))
	m.reset()
}

func (m *Metrics) fixReadDuration() {
	m.read = time.Since(m.tm)
	m.stat = append(m.stat, fmt.Sprintf(statPattern, "read", prepareValue(m.read)))
	m.reset()
}

func (m *Metrics) fixHandleDuration() {
	m.handle = time.Since(m.tm)
	m.stat = append(m.stat, fmt.Sprintf(statPattern, "handle", prepareValue(m.handle)))
	m.reset()
}

func (m *Metrics) fixWriteDuration() {
	m.write = time.Since(m.tm)
	m.stat = append(m.stat, fmt.Sprintf(statPattern, "write", prepareValue(m.write)))
	m.reset()
}

func (m *Metrics) string() (line string) {
	var total time.Duration

	if m.handshake >= 0 {
		total += m.handshake
	}

	if m.resume >= 0 {
		total += m.resume
	}

	if m.read >= 0 {
		total += m.read
	}

	if m.handle >= 0 {
		total += m.handle
	}

	if m.write >= 0 {
		total += m.write
	}

	if total > 0 {
		m.stat = append(m.stat, fmt.Sprintf(statPattern, "total", prepareValue(total)))
	}

	line = fmt.Sprintf("%s: %s", m.topic, strings.Join(m.stat, ", "))

	return line
}

func prepareValue(dur time.Duration) (size string) {
	const valuePattern = "%d %ss"

	ns := dur.Nanoseconds()
	if ns < 1000 {
		return fmt.Sprintf(valuePattern, ns, "n")
	}

	mcs := dur.Microseconds()
	if mcs < 1000 {
		return fmt.Sprintf(valuePattern, mcs, "µ")
	}

	ms := dur.Milliseconds()
	if ms < 1000 {
		return fmt.Sprintf(valuePattern, ms, "m")
	}

	return fmt.Sprintf("%.2f s", dur.Seconds())
}
