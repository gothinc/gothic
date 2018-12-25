package logger

/**
 * @desc Formatter
 * @author zhaojiangwei
 * @date 2018-12-14 16:55
 */

type GothicLogFormatter interface {
	Format(v ...interface{}) (string, error)
	FormatFields(v EntryFieldsAny) (string, error)
	SetTimestampFormat(format string)
}