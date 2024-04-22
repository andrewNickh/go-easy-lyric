package log

import (
	"database/sql/driver"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const LstdFlags = log.LstdFlags
const Lshortfile = log.Lshortfile

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug]"
	printReleaseLevel = "[info]"
	printErrorLevel   = "[error]"
	printFatalLevel   = "[fatal]"
)

const defaultCallDepth = 3

type Logger struct {
	logDir     string
	level      int
	prefix     string
	baseLogger *log.Logger
}

func getLevel(strLevel string) int {
	var level = releaseLevel
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
		break
	case "release":
		level = releaseLevel
		break
	case "error":
		level = errorLevel
		break
	case "fatal":
		level = fatalLevel
		break
	default:
		break
	}

	return level
}

func New(strLevel string, dir string, flag int, prefix string, logCompress bool, env string, fileNameSuffix string) (*Logger, error) {
	// level
	level := getLevel(strLevel)

	fileName := "game.log"
	fileNameSuffix = strings.TrimSpace(fileNameSuffix)
	if len(fileNameSuffix) > 0 {
		fileName = "game_" + fileNameSuffix + ".log"
	}

	// logger
	var baseLogger *log.Logger
	if dir != "" {
		hook := &lumberjack.Logger{
			Filename:   dir + fileName,
			MaxSize:    10,
			MaxBackups: 30,
			MaxAge:     7,
			Compress:   logCompress,
			LocalTime:  true,
		}

		if env == "dev" { // log to console and file
			mw := io.MultiWriter(hook, os.Stdout)
			baseLogger = log.New(mw, "", flag)
		} else { // only log to file
			baseLogger = log.New(hook, "", flag)
		}
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.logDir = dir

	if prefix != "" {
		logger.prefix = fmt.Sprintf("[%s]", prefix)
	}

	return logger, nil
}

func (logger *Logger) GetLevel() int {
	return logger.level
}

func (logger *Logger) SetLevel(strLevel string) {
	level := getLevel(strLevel)
	logger.level = level
}

func (logger *Logger) GetWriter() io.Writer {
	return logger.baseLogger.Writer()
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	logger.baseLogger = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, callDepth int, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = logger.prefix + printLevel + format
	if logger.logDir == "" {
		_ = logger.baseLogger.Output(callDepth, fmt.Sprintf(SetMsgColor(level, format), a...)) // log content is colored
	} else {
		_ = logger.baseLogger.Output(callDepth, fmt.Sprintf(format, a...))
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) doPrint(level int, printLevel string, callDepth int, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	p := make([]interface{}, 0, 8)
	for _, name := range a {
		p = append(p, name, " ")
	}

	msg := logger.prefix + printLevel + fmt.Sprint(p...)
	if logger.logDir == "" {
		_ = logger.baseLogger.Output(callDepth, SetMsgColor(level, msg)) // log content is colored
	} else {
		_ = logger.baseLogger.Output(callDepth, msg)
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Print(v ...interface{}) {
	logger.doPrint(debugLevel, printDebugLevel, defaultCallDepth, gormLogFormatterNoColor(v...))
}

//func (logger *Logger) Debug(format string, a ...interface{}) {
//	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
//}
//
//func (logger *Logger) Release(format string, a ...interface{}) {
//	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
//}
//
//func (logger *Logger) Error(format string, a ...interface{}) {
//	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
//}
//
//func (logger *Logger) Fatal(format string, a ...interface{}) {
//	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
//}

var gLogger, _ = New("debug", "", Lshortfile|LstdFlags, "", false, "dev", "")

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if gLogger != nil {
		gLogger.Close()
		gLogger = nil
	}

	if logger != nil {
		gLogger = logger
	}
}

func GetLogger() *Logger {
	return gLogger
}

// doPrintf
func Debugf(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, defaultCallDepth, format, a...)
}

func Infof(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, defaultCallDepth, format, a...)
}

func Errorf(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, defaultCallDepth, format, a...)
}

func Fatalf(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, defaultCallDepth, format, a...)
}

// doPrint
func Debug(a ...interface{}) {
	gLogger.doPrint(debugLevel, printDebugLevel, defaultCallDepth, a...)
}

func Info(a ...interface{}) {
	gLogger.doPrint(releaseLevel, printReleaseLevel, defaultCallDepth, a...)
}

func Error(a ...interface{}) {
	gLogger.doPrint(errorLevel, printErrorLevel, defaultCallDepth, a...)
}

func Fatal(a ...interface{}) {
	gLogger.doPrint(fatalLevel, printFatalLevel, defaultCallDepth, a...)
}

// customized call level output file name
func DebugfCallDepth(callDepth int, format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, callDepth+defaultCallDepth, format, a...)
}

func DebugCallDepth(callDepth int, a ...interface{}) {
	gLogger.doPrint(debugLevel, printDebugLevel, callDepth+defaultCallDepth, a...)
}

func Close() {
	gLogger.Close()
}

//gorm override output format

var sqlRegexp = regexp.MustCompile(`\?`)
var numericPlaceHolderRegexp = regexp.MustCompile(`\$\d+`)

var NowFunc = func() time.Time {
	return time.Now()
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

var gormLogFormatterNoColor = func(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			currentTime     = "\n[" + NowFunc().Format("2006-01-02 15:04:05") + "]"
			source          = fmt.Sprintf("(%v)", values[1])
		)

		messages = []interface{}{source, currentTime}

		if len(values) == 2 {
			//remove the line break
			currentTime = currentTime[1:]
			//remove the brackets
			source = fmt.Sprintf("%v", values[1])

			messages = []interface{}{currentTime, source}
		}

		if level == "sql" {
			// duration
			messages = append(messages, fmt.Sprintf(" [%.2fms] ", float64(values[2].(time.Duration).Nanoseconds()/1e4)/100.0))
			// sql

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						if t.IsZero() {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", "0000-00-00 00:00:00"))
						} else {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
						}
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						switch value.(type) {
						case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
							formattedValues = append(formattedValues, fmt.Sprintf("%v", value))
						default:
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						}
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if numericPlaceHolderRegexp.MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range sqlRegexp.Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages = append(messages, sql)
			messages = append(messages, fmt.Sprintf(" \n[%v] ", strconv.FormatInt(values[5].(int64), 10)+" rows affected or returned "))
		} else {
			//messages = append(messages, "\033[31;1m")
			messages = append(messages, values[2:]...)
			//messages = append(messages, "\033[0m")
		}
	}

	return
}
