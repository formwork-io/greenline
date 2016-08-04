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
import "time"

var logging = false
var printing = true

// DisableLogging disables logging.
func DisableLogging() { logging = false }

// EnableLogging enables logging.
func EnableLogging() { logging = true }

// DisablePrinting disables printing.
func DisablePrinting() { printing = false }

// EnablePrinting enables printing.
func EnablePrinting() { printing = true }

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
}

// Out calls Print and Log with the arguments and calls os.Exit(1).
func Out(msg string, args ...interface{}) {
	Print(msg, args...)
	Log(msg, args...)
}

// Die calls Print and Log with the arguments and calls os.Exit(1).
func Die(msg string, args ...interface{}) {
	Print(msg, args...)
	Log(msg, args...)
	os.Exit(1)
}
