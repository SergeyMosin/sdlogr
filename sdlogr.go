// Package sdlogr: another implementation of logr interface with addition of systemd specific prefixes (severity levels). This logger is meant for apps/services that are started via systemd and want to send their logs to system journal.
package sdlogr

import (
	"bytes"
	"fmt"
	"github.com/go-logr/logr"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

const (
	// from systemd/sd-daemon.h
	//
	// #define SD_EMERG   "<0>"  /* system is unusable */
	// #define SD_ALERT   "<1>"  /* action must be taken immediately */
	// #define SD_CRIT    "<2>"  /* critical conditions */
	// #define SD_ERR     "<3>"  /* error conditions */
	// #define SD_WARNING "<4>"  /* warning conditions */
	// #define SD_NOTICE  "<5>"  /* normal but significant condition */
	// #define SD_INFO    "<6>"  /* informational */
	// #define SD_DEBUG   "<7>"  /* debug-level messages */

	sdLvlError = "<3>"
	sdLvlInfo  = "<6>"
)
const emptyStringPlaceholder = "\"\""

var mergeSeparators = [2]string{": ", ", "}

type Options struct {
	// Depth biases the assumed number of call frames to the "true" caller.
	// Values less than zero will be treated as zero.
	Depth int

	// Verbosity tells sdlogr which V logs to produce.  Higher values enable
	// more logs.  Info logs at or below this level will be written, while logs
	// above this level will be discarded.
	Verbosity int

	// LogCallerInfo if this is false caller file name and line number will not be logged in logr.Info. Default is true, caller file name and line number are always logged in logr.Error.
	LogCallerInfo bool

	// Out where to send logs. Defaults to os.Stdout
	Out io.Writer
}

// New returns a logr.Logger
func New() logr.Logger {
	return NewWithOptions(Options{LogCallerInfo: true})
}

// NewWithOptions returns a logr.Logger
func NewWithOptions(opts Options) logr.Logger {

	if opts.Depth < 0 {
		opts.Depth = 0
	}
	if opts.Verbosity < 0 {
		opts.Verbosity = 0
	}
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	l := sdLogr{
		level:         opts.Verbosity,
		depth:         opts.Depth + 1,
		logCallerInfo: opts.LogCallerInfo,
		out:           opts.Out,
		valuesMap:     make(map[string]interface{}, 1),
	}
	return logr.New(&l)
}

type sdLogr struct {
	level         int
	depth         int
	prefix        string
	valuesMap     map[string]interface{}
	valuesStr     string
	out           io.Writer
	logCallerInfo bool
}

func (l *sdLogr) Init(info logr.RuntimeInfo) {
	if info.CallDepth < 0 {
		return
	}
	l.depth += info.CallDepth
}

func (l *sdLogr) Enabled(level int) bool {
	return level <= l.level
}

func (l *sdLogr) Info(_ int, msg string, kvList ...interface{}) {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.WriteString(sdLvlInfo)

	// we do npt log if there is no message or values
	hasData := false

	if len(l.prefix) > 0 {
		buf.WriteString(l.prefix + ": ")
	}

	if l.logCallerInfo {
		if _, file, line, ok := runtime.Caller(l.depth); ok {
			buf.WriteString(file[getLastSlashPos(file)+1:] + ": " + strconv.Itoa(line) + ", ")
		} else {
			buf.WriteString("???: 0, ")
		}
	}

	if msg != "" {
		buf.WriteString(msg + ", ")
		hasData = true
	}

	if l.valuesStr != "" {
		buf.WriteString(l.valuesStr)
		hasData = true
	}

	kvLen := len(kvList)
	if kvLen > 0 {
		bufferKv(buf, kvList, kvLen)
		hasData = true
	}

	if hasData {
		write(l.out, buf)
	}
}

func (l *sdLogr) Error(err error, msg string, kvList ...interface{}) {

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.WriteString(sdLvlError)

	if len(l.prefix) > 0 {
		buf.WriteString(l.prefix + ": ")
	}

	if _, file, line, ok := runtime.Caller(l.depth); ok {
		buf.WriteString(file[getLastSlashPos(file)+1:] + ": " + strconv.Itoa(line) + ", ")
	} else {
		buf.WriteString("???: 0, ")
	}

	if err != nil {
		buf.WriteString(err.Error() + ", ")
	} else {
		buf.WriteString("<nil>, ")
	}

	if msg != "" {
		buf.WriteString(msg + ", ")
	}

	if l.valuesStr != "" {
		buf.WriteString(l.valuesStr)
	}

	kvLen := len(kvList)
	if kvLen > 0 {
		bufferKv(buf, kvList, kvLen)
	}

	write(l.out, buf)
}

// WithName returns a new logr.Logger with the specified name appended, '/' character is used to separate multiple names.
func (l sdLogr) WithName(name string) logr.LogSink {
	if len(l.prefix) > 0 {
		l.prefix += "/"
	}
	l.prefix += name
	return &l
}

func (l sdLogr) WithValues(kvList ...interface{}) logr.LogSink {
	// because we have data before the values string...
	n := len(kvList)
	if n == 0 {
		return &l
	}

	if (n & 1) != 0 {
		// for some reason a value is missing
		// let's add an empty string
		kvList = append(kvList, emptyStringPlaceholder)
		n++
	}

	for i := 0; i < n; i++ {
		k := kvList[i]
		if k == "" {
			k = emptyStringPlaceholder
		}
		i++
		l.valuesMap[fmt.Sprintf("%v", deref(k))] = kvList[i]
	}

	// rebuild the string
	l.valuesStr = ""
	for k, v := range l.valuesMap {
		if v == "" {
			v = emptyStringPlaceholder
		}
		l.valuesStr += k + ": " + fmt.Sprintf("%v", deref(v)) + ", "
	}
	return &l
}

func (l sdLogr) WithCallDepth(depth int) logr.LogSink {
	l.depth += depth
	return &l
}

// UnmarshalStruct converts struct to string (including keys) using "%+v" format
func UnmarshalStruct(i interface{}) string {
	return fmt.Sprintf("%+v", deref(i))
}

func write(out io.Writer, buf *bytes.Buffer) {
	// cut off extra separator and add LF(10)
	newLen := buf.Len() - 1
	finalBytes := buf.Bytes()[:newLen]
	newLen -= 1
	if finalBytes[newLen] != 10 {
		finalBytes[newLen] = 10
	}
	_, _ = out.Write(finalBytes)
}

func bufferKv(buf *bytes.Buffer, kvList []interface{}, kvLen int) {
	for i := 0; i < kvLen; i++ {
		v := kvList[i]
		if v == "" {
			v = emptyStringPlaceholder
		}
		_, _ = fmt.Fprintf(buf, "%v", deref(v))
		buf.WriteString(mergeSeparators[i&1])
	}
}

func deref(i interface{}) interface{} {
	// try Type switch first
	switch i.(type) {

	case string, int, bool, uint, int8, int16, int32, int64, uint8, uint16, uint32, uint64, uintptr, []string, []int, []bool, []uint, []int8, []int16, []int32, []int64, []uint8, []uint16, []uint32, []uint64, []uintptr, error:
		return i
	default:
		// use reflect
		v := reflect.ValueOf(i)
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil
			} else {
				return v.Elem().Interface()
			}
		}
		return i
	}
}

func getLastSlashPos(str string) int {
	i := len(str) - 1
	for ; i > 0; i-- {
		if str[i] == '/' {
			break
		}
	}
	return i
}

var _ logr.LogSink = &sdLogr{}
var _ logr.CallDepthLogSink = &sdLogr{}
