package log

import (
	"code.google.com/p/log4go"
	conf "com/mc2/trade/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/* global logger    */
var AccessLogger log4go.Logger
var Logger log4go.Logger
var initialized bool = false

/* logDirCreate(): check and create dir if nonexist   */
func logDirCreate(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		/* create directory */
		err = os.MkdirAll(logDir, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

/* filenameGen(): generate filename    */
func filenameGen(progName, logDir string, isErrLog bool) string {
	/* remove the last '/'  */
	strings.TrimSuffix(logDir, "/")

	var fileName string
	if isErrLog {
		/* for log file of warning, error, critical  */
		fileName = filepath.Join(logDir, progName+".log.wf")
	} else {
		/* for log file of all log  */
		fileName = filepath.Join(logDir, progName+".log")
	}

	return fileName
}

/* convert level in string to log4go level  */
func stringToLevel(str string) log4go.LevelType {
	var level log4go.LevelType

	str = strings.ToUpper(str)

	switch str {
	case "DEBUG":
		level = log4go.DEBUG
	case "TRACE":
		level = log4go.TRACE
	case "INFO":
		level = log4go.INFO
	case "WARNING":
		level = log4go.WARNING
	case "ERROR":
		level = log4go.ERROR
	case "CRITICAL":
		level = log4go.CRITICAL
	default:
		level = log4go.INFO
	}
	return level
}

/*
* Init - initialize log lib
*
* PARAMS:
*   - progName: program name. Name of log file will be progName.log
*   - levelStr: "DEBUG", "TRACE", "INFO", "WARNING", "ERROR", "CRITICAL"
*   - logDir: directory for log. It will be created if noexist
*   - hasStdOut: whether to have stdout output
*   - when:
*       "M", minute
*       "H", hour
*       "D", day
*       "MIDNIGHT", roll over at midnight
*   - backupCount: If backupCount is > 0, when rollover is done, no more than
*       backupCount files are kept - the oldest ones are deleted.
*
* RETURNS:
*   nil, if succeed
*   error, if fail
 */
func InitLog4go(progName string, levelStr string, logDir string,
	hasStdOut bool, when string, backupCount int) error {
	if initialized {
		return errors.New("Initialized Already")
	}

	/* check, and create dir if nonexist    */
	if err := logDirCreate(logDir); err != nil {
		log4go.Error("Init(), in logDirCreate(%s)", logDir)
		return err
	}

	/* convert level from string to log4go level    */
	level := stringToLevel(levelStr)

	/* create logger    */
	AccessLogger = make(log4go.Logger)

	/* create writer for stdout */
	if hasStdOut {
		Logger.AddFilter("stdout", level, log4go.NewConsoleLogWriter())
	}

	/* create file writer for all log   */
	fileName := filenameGen(progName, logDir, false)
	logWriter := log4go.NewTimeFileLogWriter(fileName, when, backupCount)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileName)
	}
	logWriter.SetFormat("[NOTICE] %D %t %M")
	AccessLogger.AddFilter("log", log4go.INFO, logWriter)

	Logger = make(log4go.Logger)
	/* create file writer for warning and fatal log */
	fileNameWf := filenameGen(progName, logDir, true)
	logWriter = log4go.NewTimeFileLogWriter(fileNameWf, when, backupCount)
	if logWriter == nil {
		return fmt.Errorf("error in log4go.NewTimeFileLogWriter(%s)", fileNameWf)
	}
	logWriter.SetFormat("[%D %T] [%L] (%S) %M")
	Logger.AddFilter("log_wf", level, logWriter)

	initialized = true
	return nil
}

func Init(hasStdOut bool) error {
	cfg := conf.Read()

	/* initialize log   */
	/* set log buffer size  */
	log4go.SetLogBufferLength(10240)

	/* if blocking, log will be dropped */
	log4go.SetLogWithBlocking(false)
	err := InitLog4go("trade", cfg.Log.Level, cfg.Log.Path, hasStdOut, "H", cfg.Log.Save)
	if err != nil {
		//fmt.Printf("off_plat: err in log.Init():%s\n", err.Error())
		return err
	}

	return nil
}
