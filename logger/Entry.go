package logger

import (

)

/**
 * @desc Entry
 * @author zhaojiangwei
 * @date 2018-12-13 18:26
 */

const defaultEntryFieldsSize = 4

type EntryFields map[string]interface{}

type Entry struct {
	logger *GothicLogger
	fields EntryFields
}

func NewEntry(logger *GothicLogger, fields EntryFields) *Entry{
	return &Entry{
		logger: logger,
		fields: fields,
	}
}

func (entry *Entry) Debug() {
	if entry.logger.logLevel > LevelDebug {
		return
	}

	logName := entry.logger.prefix + getLevelName(LevelDebug) + entry.logger.suffix
	entry.logger.writeFields(logName, entry.fields)
}

func (entry *Entry) Access() {
	if entry.logger.logLevel > LevelAccess {
		return
	}

	logName := entry.logger.prefix + getLevelName(LevelAccess) + entry.logger.suffix
	entry.logger.writeFields(logName, entry.fields)
}

func (entry *Entry) Warn() {
	if entry.logger.logLevel > LevelWarn {
		return
	}

	logName := entry.logger.prefix + getLevelName(LevelWarn) + entry.logger.suffix
	entry.logger.writeFields(logName, entry.fields)
}

func (entry *Entry) Error() {
	if entry.logger.logLevel > LevelError {
		return
	}

	logName := entry.logger.prefix + getLevelName(LevelError) + entry.logger.suffix
	entry.logger.writeFields(logName, entry.fields)
}

/**
 * @desc 扩展log，如play.log
 * @author zhaojiangwei
 * @date 18:29 2018/12/12
 * @param logName:日志名，如play
 * @return
 **/
func (entry *Entry) Extend(logName string) {
	logName = entry.logger.prefix + logName + entry.logger.suffix
	entry.logger.writeFields(logName, entry.fields)
}