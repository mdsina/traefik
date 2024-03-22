package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"strings"
)

type FieldsDecoratingJSONFormatter struct {
	Delegate *logrus.JSONFormatter
	Fields   *logrus.Fields
}

func CreateJsonFormatter(fields *map[string]string, fieldsMap *map[string]string) *FieldsDecoratingJSONFormatter {
	return &FieldsDecoratingJSONFormatter{
		Delegate: &logrus.JSONFormatter{FieldMap: *CreateFieldsMap(fieldsMap)},
		Fields:   CreateFields(fields),
	}
}

func (f *FieldsDecoratingJSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	e := entry.WithFields(*f.Fields)
	e.Level = entry.Level
	e.Message = entry.Message
	e.Caller = entry.Caller

	return f.Delegate.Format(e)
}

func CreateFields(fields *map[string]string) *logrus.Fields {
	f := make(logrus.Fields, len(*fields))
	for k, v := range *fields {
		v = strings.TrimSpace(v)
		if strings.HasPrefix(v, "$") {
			env, _ := strings.CutPrefix(v, "$")
			env = os.Getenv(env)
			if env != "" {
				v = env
			}
		}
		f[k] = v
	}
	return &f
}

func CreateFieldsMap(fieldsMap *map[string]string) *logrus.FieldMap {
	f := make(logrus.FieldMap, len(*fieldsMap))
	ps := reflect.ValueOf(&f).Elem()
	t := reflect.TypeOf(f).Key()
	for k, v := range *fieldsMap {
		ps.SetMapIndex(reflect.ValueOf(k).Convert(t), reflect.ValueOf(v))
	}
	return &f
}
