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

import "flag"
import "os"
import "strconv"
import "fmt"

// GreenlineConfig contains the configuration.
type GreenlineConfig struct {
	Reload  bool
	Print   bool
	Logfile string
}

// Cfg contains the greenline configuration.
var Cfg GreenlineConfig

// Configure ...
func Configure() {
	flag.BoolVar(&Cfg.Reload, "restart", false, "enable restart on updates")
	flag.BoolVar(&Cfg.Print, "print", true, "enable console output")
	flag.StringVar(&Cfg.Logfile, "logfile", "", "enable logging to file")
	flag.Parse()

	if Cfg.Reload {
		Out("restart enabled; don't panic")
	} else {
		Out("restart disabled; buckle up Mr. Safety")
	}
	if Cfg.Print {
		Out("console output enabled; no turning back now")
	} else {
		Out("console output disabled; applying tape to mouth")
	}
	if Cfg.Logfile != "" {
		Out("logging enabled; plant a tree in reparations (%s)", Cfg.Logfile)
	} else {
		Out("logging disabled; enjoy the trees")
	}
}

func getenv(env string, dflt interface{}) (string, error) {
	_env := os.Getenv(env)
	if len(_env) == 0 {
		return "", fmt.Errorf("no %s is set", env)
	}
	return _env, nil
}

func asPort(env string) (int, error) {
	port, err := strconv.Atoi(env)
	if err != nil {
		Die("invalid port: %s", env)
		return -1, fmt.Errorf("invalid port: %v - %s", env, err.Error())
	} else if port < 1 || port > 65535 {
		Die("invalid port: %s", env)
		return -1, fmt.Errorf("invalid port: %v - %s", env, err.Error())
	}
	return port, nil
}
