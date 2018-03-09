// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.
/*
modification history
--------------------
2014/3/27, by Zhang Miao, modify, according to 
                                  https://code.google.com/p/log4go/issues/detail?id=15
*/
package log4go

import (
	"fmt"
	"bytes"
	"io"
	"sync"
)

const (
	FORMAT_DEFAULT = "[%D %T] [%L] (%S) %M"
	FORMAT_SHORT   = "[%t %d] [%L] %M"
	FORMAT_ABBREV  = "[%L] %M"
)

type formatCacheType struct {
	LastUpdateMilliseconds    int64
	shortTime, shortDate string
	longTime, longDate   string
}

var formatCache = &formatCacheType{}
var formatMutex sync.Mutex

// Known format codes:
// %T - Time (15:04:05 MST)
// %t - Time (15:04)
// %D - Date (2006/01/02)
// %d - Date (01/02/06)
// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
// %S - Source
// %M - Message
// Ignores unknown formats
// Recommended: "[%D %T] [%L] (%S) %M"
func FormatLogRecord(format string, rec *LogRecord) string {
	if rec == nil {
		return "<nil>"
	}
	if len(format) == 0 {
		return ""
	}

	out := bytes.NewBuffer(make([]byte, 0, 64))
	msec := rec.Created.UnixNano() / 1e6

    formatMutex.Lock()
	cache := *formatCache
	formatMutex.Unlock()
	if cache.LastUpdateMilliseconds != msec {
		month, day, year := rec.Created.Month(), rec.Created.Day(), rec.Created.Year()
		hour, minute, second := rec.Created.Hour(), rec.Created.Minute(), rec.Created.Second()
        millisecond := rec.Created.Nanosecond() / 1e6
		//zone, _ := rec.Created.Zone()
		updated := &formatCacheType{
			LastUpdateMilliseconds: msec,
            shortTime:         fmt.Sprintf("%02d:%02d:%02d", hour, minute, second),
			shortDate:         fmt.Sprintf("%02d/%02d/%02d", month, day, year%100),
			longTime:          fmt.Sprintf("%02d:%02d:%02d.%d", hour, minute, second, millisecond),
			longDate:          fmt.Sprintf("%04d-%02d-%02d", year, month, day),
		}
		formatMutex.Lock()
		cache = *updated
		formatCache = updated
		formatMutex.Unlock()
	}

	// Split the string into pieces by % signs
	pieces := bytes.Split([]byte(format), []byte{'%'})

	// Iterate over the pieces, replacing known formats
	for i, piece := range pieces {
		if i > 0 && len(piece) > 0 {
			switch piece[0] {
			case 'T':
				out.WriteString(cache.longTime)
			case 't':
				out.WriteString(cache.shortTime)
			case 'D':
				out.WriteString(cache.longDate)
			case 'd':
				out.WriteString(cache.shortDate)
			case 'L':
				out.WriteString(levelStrings[rec.Level])
			case 'S':
				out.WriteString(rec.Source)
			case 'M':
				out.WriteString(rec.Message)
			}
			if len(piece) > 1 {
				out.Write(piece[1:])
			}
		} else if len(piece) > 0 {
			out.Write(piece)
		}
	}
	out.WriteByte('\n')

	return out.String()
}

// This is the standard writer that prints to standard output.
type FormatLogWriter chan *LogRecord

// This creates a new FormatLogWriter
func NewFormatLogWriter(out io.Writer, format string) FormatLogWriter {
	records := make(FormatLogWriter, LogBufferLength)
	go records.run(out, format)
	return records
}

func (w FormatLogWriter) run(out io.Writer, format string) {
	for rec := range w {
		fmt.Fprint(out, FormatLogRecord(format, rec))
	}
}

// This is the FormatLogWriter's output method.  This will block if the output
// buffer is full.
func (w FormatLogWriter) LogWrite(rec *LogRecord) {
    if !LogWithBlocking {
        if len(w) >= LogBufferLength {
//            if WithModuleState {
//                log4goState.Inc("ERR_PATT_LOG_OVERFLOW", 1)
//            }            
            
            return
        }
    }
    
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w FormatLogWriter) Close() {
	close(w)
}

// This is the FormatLogWriter's output method implement io.Writer
func (w FormatLogWriter) Write(p []byte) (n int, err error){
    return 0, nil
}