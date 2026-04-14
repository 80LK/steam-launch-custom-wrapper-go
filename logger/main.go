package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type logger struct {
	Prefix string
}

var logFile *os.File = nil

func NewLogger() *logger {
	return NewLoggerP("Logger")
}

func NewLoggerP(prefix string) *logger {
	if err := initLogFile(); err != nil {
		panic(err)
	}
	return &logger{
		Prefix: prefix,
	}
}

func initLogFile() error {
	if logFile != nil {
		return nil
	}
	path := getAppDataFilePath("logs/log.wrapper.txt")

	if _, err := os.Stat(path); err == nil {
		bckPath := path + ".bck"
		err = os.Remove(bckPath)
		if err != nil {
			return err
		}
		err = os.Rename(path, bckPath)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	logFile, _ = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	return nil
}

func (l *logger) getPrefixLine(level string) string {
	return fmt.Sprintf("[%s][%s][%s]",
		time.Now().Format("2006-01-02 15:04:05"),
		strings.ToUpper(l.Prefix),
		strings.ToUpper(level),
	)
}

func (l *logger) write(prefix, content string) {
	fmt.Print(prefix, content)
	logFile.Write([]byte(prefix + " " + content))
}

func (l *logger) log(level string, args ...any) {
	l.write(level, fmt.Sprintln(args...))
}

func (l *logger) logF(level, format string, args ...any) {
	l.write(level, fmt.Sprintf(format, args...))
}

func (l *logger) Log(a ...any) {
	l.log("LOG", a...)
}

func (l *logger) Error(a ...any) {
	l.log("Error", a...)
}

func (l *logger) Warn(a ...any) {
	l.log("Warn", a...)
}

func (l *logger) LogF(format string, a ...any) {
	l.logF("log", format, a...)
}

func (l *logger) ErrorF(format string, a ...any) {
	l.logF("error", format, a...)
}

func (l *logger) WarnF(format string, a ...any) {
	l.logF("Warn", format, a...)
}

func (l *logger) LogR(a ...any) any {
	l.Log(a...)
	return a[len(a)-1]
}

func (l *logger) WarnR(a ...any) any {
	l.Warn(a...)
	return a[len(a)-1]
}

func (l *logger) ErrorR(a ...any) any {
	l.Error(a...)
	return a[len(a)-1]
}

func (l *logger) LogFR(format string, a ...any) any {
	l.LogF(format, a...)
	return a[len(a)-1]
}

func (l *logger) WarnFR(format string, a ...any) any {
	l.WarnF(format, a...)
	return a[len(a)-1]
}

func (l *logger) ErrorFR(format string, a ...any) any {
	l.ErrorF(format, a...)
	return a[len(a)-1]
}

func getAppDataPath(file string) string {
	base := os.Getenv("APPDATA")
	appData := filepath.Join(base, "80LK", "steam-launch-custom", file)

	if filepath.IsAbs(file) {
		file = "." + file
	}

	return filepath.Join(appData, file)
}

func getAppDataFilePath(file string) string {
	full := getAppDataPath(file)
	dir := filepath.Dir(full)
	_ = os.MkdirAll(dir, 0755)
	return full
}

func getAppDataDirPath(file string) string {
	full := getAppDataPath(file)
	_ = os.MkdirAll(full, 0755)
	return full
}
