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
	exitchan := make(chan os.Signal, 0)
	signal.Notify(exitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-exitchan
        if sig == nil {
            return
        }
		cr.Out("got %s, exiting", sig.String())
		os.Exit(0)
	}()

	reloadchan := make(chan int)
	go cr.Reloader(reloadchan)
	cr.Out("greenline ready")
	router := cr.FirstOrDie(goczmq.NewRouter("tcp://*:5555")).(*goczmq.Sock)
	defer router.Destroy()

	cr.Out("router created and bound")
	for {
		select {
		case _ = <-reloadchan:
			cr.Out("new binary available, restarting greenline")
			cr.Out("RESART NOW")
	        signal.Stop(exitchan)
			router.Destroy()
			close(exitchan)
			close(reloadchan)
			cr.Restart()
		}
	}
}
