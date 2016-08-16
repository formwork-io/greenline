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

import "os"
import "strconv"
import "fmt"
import "strings"
import "gopkg.in/redis.v4"

// NoPrint is the setting "GREENLINE_NO_PRINT".
const NoPrint = "GREENLINE_NO_PRINT"

// LogFile is the setting "GREENLINE_LOG_FILE".
const LogFile = "GREENLINE_LOG_FILE"

// Reload is the setting "GREENLINE_RELOAD".
const Reload = "GREENLINE_RELOAD"

// Rails is the setting "GREENLINE_RAILS".
const Rails = "GREENLINE_RAILS"

// UseRedis is the setting "GREENLINE_USE_REDIS".
const UseRedis = "GREENLINE_USE_REDIS"

// RedisHost is the setting "GREENLINE_REDIS_HOST".
const RedisHost = "GREENLINE_REDIS_HOST"

// RedisPort is the setting "GREENLINE_REDIS_PORT".
const RedisPort = "GREENLINE_REDIS_PORT"

// RedisPassword is the setting "GREENLINE_REDIS_PASSWORD".
const RedisPassword = "GREENLINE_REDIS_PASSWORD"

// RedisDB is the setting "GREENLINE_REDIS_DB".
const RedisDB = "GREENLINE_REDIS_DB"

// GreenlineConfig contains the configuration.
type GreenlineConfig struct {
	Reload   bool
	Print    bool
	LogFile  string
	NrRails  int
	Incoming []string
	Outgoing []string
}

// Cfg contains the greenline configuration.
var Cfg GreenlineConfig

// Configure ...
func Configure() {

	// Apply the default configuration
	Cfg.Print = true
	Cfg.Reload = true
	Cfg.LogFile = ""
	Cfg.Incoming = make([]string, 0)
	Cfg.Outgoing = make([]string, 0)

	// Change configuration through env or redis
	useRedis := getenv(UseRedis, "", true)
	if useRedis == "" {
		Out("configuring via environment")
		configureViaEnv()
	} else {
		Out("configuring via redis")
		configureViaRedis()
	}
}

func configureViaEnv() {
	noPrint := getenv(NoPrint, "", true)
	logFile := getenv(LogFile, "", true)
	reload := getenv(Reload, "", true)

	rails := getenv(Rails, "", false)
	configureRails(rails)

	if noPrint == "1" {
		Cfg.Print = false
	}
	if logFile != "" {
		Cfg.LogFile = logFile

	}
	if reload == "0" {
		Cfg.Reload = false
	}
}

func configureViaRedis() {
	host := getenv(RedisHost, "", false)
	port := asPort(getenv(RedisPort, "", false))
	password := getenv(RedisPassword, "", true)
	db := asInt(getenv(RedisDB, "0", true))
	addr := fmt.Sprintf("%s:%d", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Ping redis or die trying
	FirstOrDie(client.Ping().Result())

	var cmd *redis.StringCmd
	cmd = client.Get(NoPrint)
	noPrint := FirstOrZStr(cmd.Result()).(string)
	cmd = client.Get(LogFile)
	logFile := FirstOrZStr(cmd.Result()).(string)
	cmd = client.Get(Reload)
	reload := FirstOrZStr(cmd.Result()).(string)
	cmd = client.Get(Rails)

	rails := FirstOrZStr(cmd.Result()).(string)
	configureRails(rails)

	if noPrint == "1" {
		Cfg.Print = false
	}
	if logFile != "" {
		Cfg.LogFile = logFile
	}
	if reload == "0" {
		Cfg.Reload = false
	}
}

func configureRails(setting string) {
	if setting == "" {
		Die("no rails configured")
	}
	tokens := strings.Split(setting, ",")
	if len(tokens)%2 != 0 {
		Die("Odd number (%d) of rails specified ('%s')", len(tokens), setting)
	}

	var in string
	var out string

	Cfg.NrRails = len(tokens) / 2
	for i := 0; i < len(tokens); i += 2 {
		in = tokens[i]
		in = fmt.Sprintf("tcp://%s", in)
		Cfg.Incoming = append(Cfg.Incoming, in)
		out = tokens[i+1]
		out = fmt.Sprintf("tcp://%s", out)
		Cfg.Outgoing = append(Cfg.Outgoing, out)
	}
}

func getenv(env string, dflt string, useDefault bool) string {
	_env := os.Getenv(env)
	if len(_env) == 0 && useDefault {
		return dflt
	} else if len(_env) != 0 {
		return _env
	}
	Die("%s is not set", env)
	panic("Die()")
}

func asPort(env string) int {
	port, err := strconv.Atoi(env)
	if err != nil {
		Die("invalid port: %s", env)
		panic("Die()")
	} else if port < 1 || port > 65535 {
		Die("invalid port: %s", env)
		panic("Die()")
	}
	return port
}

func asInt(env string) int {
	ret, err := strconv.Atoi(env)
	if err != nil {
		Die("invalid int: %s", env)
		panic("Die()")
	}
	return ret
}
