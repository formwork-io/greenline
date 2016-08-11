/*
Permission is hereby granted, free of charge, to any person
obtaining a copy of this software and associated documentation
files (the "Software"), to deal in the Software without
restriction, including without limitation the rights to use,
copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following
conditions:
The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

See http://formwork-io.github.io/ for more.
*/

package main

import "os"
import "os/signal"
import "fmt"
import "syscall"
import "goczmq-1.0"
import cr "core"

func main() {
	info := "greenline: notoriously unreliable\n" +
		"https://github.com/formwork-io/greenline\n" +
		"This is free software with ABSOLUTELY NO WARRANTY."
	fmt.Printf("%s\n--\n", info)
	cr.Out("greenline alive")
	cr.Configure()

	// migrate signal/handler out and into core
	exitchan := make(chan os.Signal, 0)
	signal.Notify(exitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	reloadchan := make(chan int)
	if cr.Cfg.Reload {
		cr.Out("restart enabled; don't panic")
		go cr.Reloader(reloadchan)
	}
	cr.Out("greenline ready")
	router := cr.FirstOrDie(goczmq.NewRouter("tcp://*:5555")).(*goczmq.Sock)
	defer router.Destroy()

	cr.Out("router created and bound")

	var restarting = false

For:
	for {
		select {
		case sig := <-exitchan:
			if sig == nil {
				continue
			}
			if sig == syscall.SIGQUIT {
				cr.Out("exiting immediately on signal (%s)", sig.String())
				os.Exit(1)
			}
			cr.Out("initiating graceful shutdown on signal (%s)", sig.String())
			break For
		case _ = <-reloadchan:
			cr.Out("new binary available, restarting greenline")
			restarting = true
			break For
		}
	}

	signal.Stop(exitchan)
	router.Destroy()
	goczmq.Shutdown()
	close(exitchan)
	close(reloadchan)

	if restarting {
		// sleep a moment before restarting
		cr.SleepMs(250)
		cr.Restart()
	}
	cr.Out("and the rest, after a sudden wet thud, was silence")
	os.Exit(0)
}
