package hook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// RFC5424 log message levels.
const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

// LoggerInterface Logger接口
type LoggerInterface interface {
	Init(config string) error
	WriteMsg(msg string, level int) error
	Destroy()
	Flush()
}

// LogWriter implements LoggerInterface.
// It writes messages by lines limit, file size limit, or time frequency.
type LogWriter struct {
	*log.Logger
	mw *MuxWriter
	// The opened file
	Filename string `json:"filename"`

	Maxlines         int `json:"maxlines"`
	maxlinesCurlines int

	// Rotate at size
	Maxsize        int `json:"maxsize"`
	maxsizeCursize int

	// Rotate daily
	Daily         bool  `json:"daily"`
	Maxdays       int64 `json:"maxdays"`
	dailyOpendate int

	Rotate bool `json:"rotate"`

	startLock sync.Mutex // Only one log can write to the file

	Level int `json:"level"`
}

// MuxWriter an *os.File writer with locker.
type MuxWriter struct {
	sync.Mutex
	fd *os.File
}

// write to os.File.
func (l *MuxWriter) Write(b []byte) (int, error) {
	l.Lock()
	defer l.Unlock()
	return l.fd.Write(b)
}

// SetFd set os.File in writer.
func (l *MuxWriter) SetFd(fd *os.File) {
	if l.fd != nil {
		_ = l.fd.Close()
	}
	l.fd = fd
}

// NewFileWriter create a FileLogWriter returning as LoggerInterface.
func NewFileWriter() LoggerInterface {
	w := &LogWriter{
		Filename: "",
		Maxlines: 1000000,
		Maxsize:  1 << 28, //256 MB
		Daily:    true,
		Maxdays:  7,
		Rotate:   true,
		Level:    LevelDebug,
	}
	// use MuxWriter instead direct use os.File for lock write when rotate
	w.mw = new(MuxWriter)
	// set MuxWriter as Logger's io.Writer
	w.Logger = log.New(w.mw, "", log.Ldate|log.Ltime)
	return w
}

// Init file logger with json config.
// jsonconfig like:
//	{
//	"filename":"logs/sample.log",
//	"maxlines":10000,
//	"maxsize":1<<30,
//	"daily":true,
//	"maxdays":15,
//	"rotate":true
//	}
func (w *LogWriter) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), w)
	if err != nil {
		return err
	}
	if len(w.Filename) == 0 {
		return errors.New("jsonconfig must have filename")
	}
	err = w.startLogger()
	return err
}

// start file logger. create log file and set to locker-inside file writer.
func (w *LogWriter) startLogger() error {
	fd, err := w.createLogFile()
	if err != nil {
		return err
	}
	w.mw.SetFd(fd)
	err = w.initFd()
	if err != nil {
		return err
	}
	return nil
}

func (w *LogWriter) docheck(size int) {
	w.startLock.Lock()
	defer w.startLock.Unlock()
	if w.Rotate && ((w.Maxlines > 0 && w.maxlinesCurlines >= w.Maxlines) ||
		(w.Maxsize > 0 && w.maxsizeCursize >= w.Maxsize) ||
		(w.Daily && time.Now().Day() != w.dailyOpendate)) {
		if err := w.DoRotate(); err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.Filename, err)
			return
		}
	}
	w.maxlinesCurlines++
	w.maxsizeCursize += size
}

// WriteMsg write logger message into file.
func (w *LogWriter) WriteMsg(msg string, level int) error {
	if level > w.Level {
		return nil
	}
	n := 24 + len(msg) // 24 stand for the length "2013/06/23 21:00:22 [T] "
	w.docheck(n)
	w.Logger.Print(msg)
	return nil
}

func (w *LogWriter) createLogFile() (*os.File, error) {
	// Open the log file
	fd, err := os.OpenFile(w.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return fd, err
}

func (w *LogWriter) initFd() error {
	fd := w.mw.fd
	finfo, err := fd.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s", err)
	}
	w.maxsizeCursize = int(finfo.Size())
	w.dailyOpendate = time.Now().Day()
	if finfo.Size() > 0 {
		content, err := ioutil.ReadFile(w.Filename)
		if err != nil {
			return err
		}
		w.maxlinesCurlines = len(strings.Split(string(content), "\n"))
	} else {
		w.maxlinesCurlines = 0
	}
	return nil
}

// DoRotate means it need to write file in new file.
// new file name like xx.log.2013-01-01.2
func (w *LogWriter) DoRotate() error {
	_, err := os.Lstat(w.Filename)
	if err == nil { // file exists
		// Find the next available number
		num := 1
		fname := ""
		for ; err == nil && num <= 999; num++ {
			fname = w.Filename + fmt.Sprintf(".%s.%03d", time.Now().Format("2006-01-02"), num)
			_, err = os.Lstat(fname)
		}
		// return error if the last file checked still existed
		if err == nil {
			return fmt.Errorf("Rotate: Cannot find free log number to rename %s", w.Filename)
		}

		// block Logger's io.Writer
		w.mw.Lock()
		defer w.mw.Unlock()

		fd := w.mw.fd
		_ = fd.Close()

		// close fd before rename
		// Rename the file to its newfound home
		err = os.Rename(w.Filename, fname)
		if err != nil {
			return fmt.Errorf("Rotate: %s", err)
		}

		// re-start logger
		err = w.startLogger()
		if err != nil {
			return fmt.Errorf("Rotate StartLogger: %s", err)
		}

		go w.deleteOldLog()
	}

	return nil
}

func (w *LogWriter) deleteOldLog() {
	dir := filepath.Dir(w.Filename)
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				returnErr = fmt.Errorf("Unable to delete old log '%s', error: %+v", path, r)
				fmt.Println(returnErr)
			}
		}()

		if !info.IsDir() && info.ModTime().Unix() < (time.Now().Unix()-60*60*24*w.Maxdays) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(w.Filename)) {
				_ = os.Remove(path)
			}
		}
		return
	})
}

// Destroy destroy file logger, close file writer.
func (w *LogWriter) Destroy() {
	_ = w.mw.fd.Close()
}

// Flush file logger.
// there are no buffering messages in file logger in memory.
// flush file means sync file from disk.
func (w *LogWriter) Flush() {
	_ = w.mw.fd.Sync()
}
