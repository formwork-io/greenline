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
import "sync"
import "goczmq-1.0"
import cr "core"

func main() {
	info := "greenline: notoriously unreliable\n" +
		"https://github.com/formwork-io/greenline\n" +
		"This is free software with ABSOLUTELY NO WARRANTY."
	cr.Out("%s\n--", info)
	cr.Out("greenline alive")

	cfg := cr.Configure()
	signalingCfg := cr.ConfigureSignaling()
	reloadingCfg := cr.ConfigureReloading(&cfg)

	cr.Out("configured %d rails", cr.Cfg.NrRails)
	for i := 0; i < cr.Cfg.NrRails; i++ {
		in := cr.Cfg.Incoming[i]
		out := cr.Cfg.Outgoing[i]
		cr.Out("rail #%d: %s ==> %s", i, in, out)
	}
	cr.ConfigureLogging(cr.Cfg)
	cr.Out("greenline ready")

	var wg sync.WaitGroup
	var railSentinel = false
	var rails = make([]*cr.Rail, 0)

	for i := 0; i < cr.Cfg.NrRails; i++ {
		in := cr.Cfg.Incoming[i]
		out := cr.Cfg.Outgoing[i]
		inSock := cr.FirstOrDie(goczmq.NewPull(in)).(*goczmq.Sock)
		inSock.SetRcvhwm(1)
		outSock := cr.FirstOrDie(goczmq.NewPush(out)).(*goczmq.Sock)
		outSock.SetSndhwm(1)
		rail := cr.MakeRail(i, inSock, outSock)
		rails = append(rails, &rail)
		wg.Add(1)
		go handleRail(&wg, &railSentinel, &rail)
	}

	var restarting = false
	var forever = true

	for {
		if !forever {
			break
		}
		select {
		case sig := <-signalingCfg.SignalChannel:
			if sig == nil {
				continue
			}
			if sig == cr.SIGQUIT {
				cr.Out("exiting immediately on signal (%s)", sig.String())
				os.Exit(1)
			} else if sig == cr.SIGTERM || sig == cr.SIGINT {
				cr.Out("initiating graceful shutdown on signal (%s)", sig.String())
				forever = false
				// break For
			} else if sig == cr.SIGUSR1 {
				for _, rail := range rails {
					rail.DumpStats()
				}
			} else if sig == cr.SIGABRT {
				cr.Out("ABORT")
				os.Exit(1)
			}
		case _ = <-reloadingCfg.ReloadChannel:
			cr.Out("new binary available, restarting greenline")
			restarting = true
			forever = false
			// break For
		}
	}

	railSentinel = true
	cr.Out("bringing rails offline")
	wg.Wait()
	for _, rail := range rails {
		rail.Destroy()
	}

	cr.ShutdownReloading(reloadingCfg)
	goczmq.Shutdown()
	cr.ShutdownLogging()
	cr.ShutdownSignaling(signalingCfg)

	if restarting {
		// sleep a moment (avoids text file busy)
		cr.SleepMs(250)
		cr.Restart()
	}
	cr.Out("and the rest, after a sudden wet thud, was silence")
	os.Exit(0)
}

func handleRail(wg *sync.WaitGroup, sentinel *bool, rail *cr.Rail) {
	poller := cr.FirstOrDie(goczmq.NewPoller()).(*goczmq.Poller)
	poller.Add(rail.Pull)
	defer poller.Destroy()
	for {
		if *sentinel {
			break
		}
		s := poller.Wait(500)
		if s == nil {
			continue
		}
		err := rail.RelayMessage()
		if err != nil {
			cr.Out("%s", err.Error())
			break
		}
	}
	wg.Done()
}
