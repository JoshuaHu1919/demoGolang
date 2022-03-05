package logformat

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

func InitLogger() zerolog.Logger {

	//init console
	consoleOutput := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	consoleOutput.FormatLevel = func(i interface{}) string {
		return fmt.Sprintf("\x1b[%dm%v\x1b[0m", consoleDefaultFormatLevel(i), strings.ToUpper(fmt.Sprintf("| %-6s|", i)))

	}
	consoleOutput.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf(" %s ", i)
	}
	consoleOutput.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	consoleOutput.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	fileIO := getFileIO()

	multi := zerolog.MultiLevelWriter(consoleOutput, fileIO)

	ZLog := zerolog.New(multi).With().Timestamp().Logger()

	return ZLog.With().Caller().Logger()
}

func consoleDefaultFormatLevel(i interface{}) int {
	if ll, ok := i.(string); ok {
		switch ll {
		case "trace":
			return colorMagenta
		case "debug":
			return colorYellow
		case "info":
			return colorGreen
		case "warn":
			return colorCyan
		case "error":
			return colorRed
		case "fatal":
			return colorBold
		case "panic":
			return colorBold
		default:
			return colorBold
		}
	}
	return colorBold
}

func getFileIO() io.Writer {
	err := mkdirLog("./log")
	if err != nil {
		log.Fatal().Msgf("error mkdir: %v", err)
	}
	projectName := projectDir()
	dateStr := time.Now().Format("150405-2006-01-02")
	logfile := fmt.Sprintf("./log/%s_%s.log", projectName, dateStr)
	f, err2 := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err2 != nil {
		log.Fatal().Msgf("error opening file: %v", err)
	}
	return f
}

func mkdirLog(dir string) (e error) {
	_, er := os.Stat(dir)
	b := er == nil || os.IsExist(er)
	if !b {
		if err := os.MkdirAll(dir, 0775); err != nil {
			if os.IsPermission(err) {
				e = err
			}
		}
	}
	return
}

func projectDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b), "../..")
	folderNames := strings.Split(d, "/")
	return folderNames[len(folderNames)-1]
}
