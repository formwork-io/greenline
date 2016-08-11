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
import fsn "fsnotify-1.3.1"
import "syscall"

// Reloader ...
func Reloader(reload chan int) {
	watcher := (FirstOrDie(fsn.NewWatcher())).(*fsn.Watcher)
	exe := BinPath()
	Print("monitoring %s", exe.dir)
	DieOnErr(watcher.Add(exe.dir))
For:
	for {
		select {
		case event := <-watcher.Events:
			if event.Name != exe.path {
				continue
			}
			if event.Op&fsn.Remove == fsn.Remove ||
				event.Op&fsn.Rename == fsn.Rename {
				continue
			}
			Out("I/O on %s (op %d)", event.Name, event.Op)
			break For
		case err := <-watcher.Errors:
			Die("failed gettiing events (%s)", err.Error())
		}
	}
	watcher.Close()
	reload <- 0
}

// Restart ...
func Restart() {
	argv0 := os.Args[0]
	argv := os.Args
	env := os.Environ()
	DieOnErr(syscall.Exec(argv0, argv, env))
}

// vim: ts=4 noexpandtab
