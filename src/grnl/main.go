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

import "log"
import "os"
import "os/signal"
import "fmt"
import "syscall"
import "time"
import "goczmq-1.0"

func main() {
	info := "greenline: notoriously unreliable\n" +
		"https://github.com/formwork-io/greenline\n" +
		"This is free software with ABSOLUTELY NO WARRANTY."
	fmt.Printf("%s\n--\n", info)
	/*
		var rails []Rail
		var err error
		rails, err = ReadEnvironment()
		if err != nil {
			die(err.Error())
		}
		pprint("configuring %d rails", len(rails))
	*/

	pprint("greenline alive")
	exitchan := make(chan os.Signal, 0)
	signal.Notify(exitchan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-exitchan
		out("received %s signal, exiting.\n", sig.String())
		os.Exit(0)
	}()

	reloadchan := make(chan int)
	/*
		go reloader(reloadchan)
	*/
	readychan := make(chan bool)
	pollchan := make(chan bool)
	go func() {
		for {
			readychan <- false
			break
		}
	}()
	pprint("greenline ready")
	router, err := goczmq.NewRouter("tcp://*:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer router.Destroy()

	log.Println("router created and bound")
	for {
		select {
		case reloadOp := <-reloadchan:
			/*
				if reloadOp&BinReload == BinReload {
					pprint("new binary available, restarting greenline")
					// exec or die
					restart()
				} else if reloadOp&ConfigReload == ConfigReload {
					pprint("new configuration available, restarting greenline")
					// exec or die
					restart()
				}
			*/
			fmt.Printf("%s\n", reloadOp)
		case ready := <-readychan:
			if !ready {
				die("ready set fail")
			}
			// ready set go
			pollchan <- true
		}
	}
}

func makeMsg(msg string, args ...interface{}) string {
	const layout = "%d%02d%02d-%02d-%02d-%02d greenline[%d]: %s"
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	seconds := now.Second()
	pid := os.Getpid()
	arg := fmt.Sprintf(msg, args...)
	ret := fmt.Sprintf(layout, year, month, day, hour, minute, seconds, pid, arg)
	return ret
}

func pprint(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg+"\n")
}

func out(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg)
	os.Stdout.Sync()
}

func die(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stderr, msg+"\n")
	os.Exit(1)
}

// vim: ts=4 noexpandtab
