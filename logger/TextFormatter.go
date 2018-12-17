package logger

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

/**
 * @desc TextFormatter
 * @author zhaojiangwei
 * @date 2018-12-14 17:02
 */

const defaultTextLogSeparator = ","

type TextFormatter struct {
	TimestampFormat string

	//是否用引号包起空值
	QuoteEmptyFields bool
	Separator string
}

func NewTextFormatter(timestampFormat string) *TextFormatter{
	return &TextFormatter{
		TimestampFormat: timestampFormat,
		Separator: defaultTextLogSeparator,
	}
}

func (this *TextFormatter) SetTimestampFormat(format string){
	this.TimestampFormat = format
}

func (this *TextFormatter) SetSeparator(s string) *TextFormatter{
	this.Separator = s
	return this
}

func (this *TextFormatter) FormatFields(v EntryFields) (string, error) {
	msgstr := ""
	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		msgstr = msgstr + fmt.Sprintf("%s=%s%s", k, this.appendValue(v[k]), this.Separator)
	}

	msgstr = strings.TrimRight(msgstr, this.Separator)
	if this.TimestampFormat != "" {
		msgstr = fmt.Sprintf("[%s] %s", time.Now().Format(this.TimestampFormat), msgstr)
	}

	return msgstr, nil
}

func (this *TextFormatter) Format(v ...interface{}) (string, error)  {
	msgstr := ""
	for _, msg := range v {
		if msg1, ok := msg.(map[string]interface{}); ok {
			var keys []string
			for k := range msg1 {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				msgstr = msgstr + fmt.Sprintf("%s=%s%s", k, this.appendValue(msg1[k]), this.Separator)
			}
		} else {
			msgstr = msgstr + fmt.Sprintf("%+v%s", msg, this.Separator)
		}
	}

	msgstr = strings.TrimRight(msgstr, this.Separator)
	if this.TimestampFormat != "" {
		msgstr = fmt.Sprintf("[%s] %s", time.Now().Format(this.TimestampFormat), msgstr)
	}

	return msgstr, nil
}

func (this *TextFormatter) appendValue(value interface{}) string{
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !this.needsQuoting(stringVal) {
		return stringVal
	} else {
		return fmt.Sprintf("%q", stringVal)
	}
}

func (this *TextFormatter) needsQuoting(text string) bool {
	if this.QuoteEmptyFields && len(text) == 0 {
		return true
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}