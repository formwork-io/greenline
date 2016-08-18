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

import "os"
import "os/signal"
import "syscall"

// import "os/signal"

// SIGINT is an alias for syscall.SIGINT.
const SIGINT = syscall.SIGINT

// SIGTERM is an alias for syscall.SIGTERM.
const SIGTERM = syscall.SIGTERM

// SIGQUIT is an alias for syscall.SIGQUIT.
const SIGQUIT = syscall.SIGQUIT

// SIGUSR1 is an alias for syscall.SIGUSR1.
const SIGUSR1 = syscall.SIGUSR1

// SIGABRT is an alias for syscall.SIGABRT.
const SIGABRT = syscall.SIGABRT

// SignalingConfiguration contains the signaling configuration.
type SignalingConfiguration struct {
	SignalChannel chan os.Signal
}

// ConfigureSignaling configures signaling.
func ConfigureSignaling() *SignalingConfiguration {
	channel := make(chan os.Signal, 0)
	cfg := SignalingConfiguration{channel}
	signal.Notify(channel, SIGINT, SIGTERM, SIGQUIT, SIGUSR1, syscall.SIGABRT)
	return &cfg
}

// ShutdownSignaling shuts down signaling.
func ShutdownSignaling(cfg *SignalingConfiguration) {
	signal.Stop(cfg.SignalChannel)
	close(cfg.SignalChannel)
}
