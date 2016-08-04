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

import "fmt"
import "log"
import "os"
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

	pprint("grnlctl alive")
	router, err := goczmq.NewRouter("tcp://*:5555")
	if err != nil {
		log.Fatal(err)
	}
	defer router.Destroy()
	log.Println("router created and bound")
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

// vim: ts=4 noexpandtab
