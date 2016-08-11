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

import "runtime"
import "time"

// Prlvalue calls Out with the %v value format for arg.
func Prlvalue(arg interface{}) {
	Out("%v\n", arg)
}

// FirstOrDie returns the first argument or calls Die if err is non-nil.
func FirstOrDie(first interface{}, err error) interface{} {
	if err == nil {
		return first
	}
	_, file, line, _ := runtime.Caller(0)
	Die("Die: %s:%d %s\n", file, line, err.Error())
	panic("Die()")
}

// SecondOrDie returns the second argument or calls Die if err is non-nil.
func SecondOrDie(err error, second interface{}) interface{} {
	if err == nil {
		return second
	}
	_, file, line, _ := runtime.Caller(0)
	Die("Die: %s:%d %s\n", file, line, err.Error())
	panic("Die()")
}

// DieOnErr calls Die if err is non-nil.
func DieOnErr(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		Die("Die: %s:%d %s\n", file, line, err.Error())
	}
}

// SleepMs sleeps for the specified number of milliseconds.
func SleepMs(ms time.Duration) {
	time.Sleep(time.Millisecond * ms)
}
