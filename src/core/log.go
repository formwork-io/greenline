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

package grnlcore

import "fmt"
import "os"
import "log"
import "time"
import "strconv"

var printing = true
var logging = false
var logFile *os.File

// DisableLogging disables logging.
func DisableLogging() { logging = false }

// EnableLogging enables logging.
func EnableLogging(path string) {
	flags := os.O_RDWR | os.O_CREATE | os.O_APPEND
	mode := os.FileMode(0660)
	logFile = FirstOrDie(os.OpenFile(path, flags, mode)).(*os.File)
	log.SetOutput(logFile)
	logging = true
}

// ConfigureLogging configures logging.
func ConfigureLogging(cfg GreenlineConfig) {
	if cfg.LogFile != "" {
		Out("enabling logging to %s", Cfg.LogFile)
		EnableLogging(cfg.LogFile)
	}
	if !cfg.Print {
		Out("disabling printing, no other output will be shown")
		DisablePrinting()
	}
}

// ShutdownLogging shuts down logging.
func ShutdownLogging() {
	if logging {
		DisableLogging()
		logFile.Close()
	}
}

// DisablePrinting disables printing.
func DisablePrinting() { printing = false }

// EnablePrinting enables printing.
func EnablePrinting() { printing = true }

func makeMsg(msg string, args ...interface{}) string {
	const layout = "grnl[%d]: %02d:%02d:%02d.%s: %s\n"
	now := time.Now()
	ns := (fmt.Sprintf("%4s", strconv.Itoa(now.Nanosecond())))[0:4]
	h, m, s := now.Hour(), now.Minute(), now.Second()
	pid := os.Getpid()
	arg := fmt.Sprintf(msg, args...)
	ret := fmt.Sprintf(layout, pid, h, m, s, ns, arg)
	return ret
}

// Print prints the arguments similar to fmt.Printf.
// If printing is disabled nothing will be shown.
func Print(msg string, args ...interface{}) {
	if !printing {
		return
	}
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stdout, msg)
	os.Stdout.Sync()
}

// Log logs the arguments similar to fmt.Printf.
// If logging is disabled nothing will be shown.
func Log(msg string, args ...interface{}) {
	if !logging {
		return
	}
	msg = makeMsg(msg, args...)
	log.Printf(msg)
}

// Out calls Print and Log with the arguments and calls os.Exit(1).
func Out(msg string, args ...interface{}) {
	Print(msg, args...)
}

// Outerr prints the arguments to stderr similar to fmt.Printf.
func Outerr(msg string, args ...interface{}) {
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stderr, msg)
	os.Stderr.Sync()
}

// Die calls Print with the arguments and calls os.Exit(1).
func Die(msg string, args ...interface{}) {
	Log(msg, args...)
	msg = makeMsg(msg, args...)
	fmt.Fprintf(os.Stderr, msg)
	os.Stderr.Sync()
	os.Exit(1)
}
