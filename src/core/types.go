// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.
//
// See http://formwork-io.github.io/ for more.

package grnlcore

import "fmt"
import "math"
import "time"
import "goczmq-4.0.1"

// Executable is the complete path to the process executable.
//
// path: absolute path to process executable
// base: last element of path
// dir: all but the last element of the path
type Executable struct {
	path string
	base string
	dir  string
}

// Rail identifies a Pull/Push pair. The rail's number serves to identify
// it in logs.
type Rail struct {
	Nr                      int
	Pull                    *goczmq.Sock
	Push                    *goczmq.Sock
	Relayed                 int
	MinimumDeliveryNS       time.Duration
	MaximumDeliveryNS       time.Duration
	LastReceivedTime        time.Time
	LastSentTime            time.Time
	MomentaryInterarrivalNS time.Duration
}

// MakeRail creates a rail.
func MakeRail(nr int, pull *goczmq.Sock, push *goczmq.Sock) Rail {
	zt := time.Time{}
	rail := Rail{nr, pull, push, 0, 0, 0, zt, zt, 0}
	return rail
}

// RelayMessage pulls a message and pushes it to its next destination.
func (r *Rail) RelayMessage() error {
	frame, flag, err := r.Pull.RecvFrameNoWait()
	received := time.Now()
	if err != nil {
		return err
	}
	r.Push.SendFrame(frame, flag)
	sent := time.Now()
	r.Relayed++

	// time-to-send duration
	ttsDuration := sent.Sub(received)

	if ttsDuration < r.MinimumDeliveryNS || r.MinimumDeliveryNS == 0 {
		r.MinimumDeliveryNS = ttsDuration
	}
	if ttsDuration > r.MaximumDeliveryNS {
		r.MaximumDeliveryNS = ttsDuration
	}

	r.MomentaryInterarrivalNS = received.Sub(r.LastReceivedTime)
	r.LastReceivedTime = received
	r.LastSentTime = sent
	return nil
}

// DumpStats prints the rail's stats to the console and log.
func (r *Rail) DumpStats() {
	pref := fmt.Sprintf("rail #%d: ", r.Nr)
	if r.Relayed == 0 {
		Outerr(fmt.Sprintf("%s: no activity", pref))
		return
	}
	now := time.Now()
	Outerr(fmt.Sprintf("%s: relayed %d messages", pref, r.Relayed))
	lastSent := now.Sub(r.LastSentTime)
	lastSentSec := math.Abs(lastSent.Seconds())
	lastReceived := now.Sub(r.LastReceivedTime)
	lastReceivedSec := math.Abs(lastReceived.Seconds())
	Outerr(fmt.Sprintf("%s: sent last message %0.2fs ago", pref, lastSentSec))
	Outerr(fmt.Sprintf("%s: received last message %0.2fs ago", pref, lastReceivedSec))
	Outerr(fmt.Sprintf("%s: momentary interarrival time %0.2fs", pref, r.MomentaryInterarrivalNS.Seconds()))
	Outerr(fmt.Sprintf("%s: best time-to-delivery %0.2fs", pref, r.MinimumDeliveryNS.Seconds()))
	Outerr(fmt.Sprintf("%s: worst time-to-delivery %0.2fs", pref, r.MaximumDeliveryNS.Seconds()))
}

// Destroy closes the Pull/Push sockets used by this rail.
func (r Rail) Destroy() {
	r.Pull.Destroy()
	r.Push.Destroy()
}
