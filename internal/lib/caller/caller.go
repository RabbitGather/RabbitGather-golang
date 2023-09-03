package caller

import (
	"fmt"
	"path/filepath"
	"runtime"
)

type formatSetting int

const (
	CALLER_FORMAT_SHORT formatSetting = iota
	CALLER_FORMAT_LONG
)

// 取得呼叫的文件與行號
func Caller(skip int, formatSet formatSetting) string {
	_, file, line, ok := runtime.Caller(1 + skip)
	if !ok {
		return "[fail to get caller]"
	}
	switch formatSet {
	case CALLER_FORMAT_SHORT:
		dir, f := filepath.Split(file)
		return fmt.Sprintf("%s/%s:%d", filepath.Base(dir), f, line)
	case CALLER_FORMAT_LONG:
		return fmt.Sprintf("%s:%d", file, line)
	default:
		panic("unknown formatSetting")
	}
}
