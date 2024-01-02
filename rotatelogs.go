package rotatelogs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotateLogs struct {
	options
	file          *os.File
	fileWildcard  string
	logPath       string
	logPathFormat string
	logLinkPath   string
	mutex         *sync.Mutex
	rotateCh      <-chan time.Time
	closeCh       chan struct{}
}

func New(logPath string, opts ...Option) (*RotateLogs, error) {
	ins := &RotateLogs{
		logPath: logPath,
		mutex:   &sync.Mutex{},
		closeCh: make(chan struct{}, 1),
	}
	ins.options = defaultOpts
	for _, opt := range opts {
		opt.apply(&ins.options)
	}
	ins.fileWildcard = fmt.Sprintf("%s.*", ins.logPath)
	ins.logLinkPath = ins.logPath
	if ins.rotateTime/time.Hour == 24 {
		ins.logPathFormat = "20060102"
	} else if ins.rotateTime%time.Hour == 0 {
		ins.logPathFormat = "2006010215"
	} else {
		ins.logPathFormat = "200601021504"
	}

	if err := os.Mkdir(filepath.Dir(ins.logPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if err := ins.rotate(time.Now()); err != nil {
		return nil, err
	}

	if ins.rotateTime != 0 {
		go ins.handleEvent()
	}

	return ins, nil
}

func (slf *RotateLogs) Write(bytes []byte) (int, error) {
	slf.mutex.Lock()
	defer slf.mutex.Unlock()
	n, err := slf.file.Write(bytes)
	return n, err
}

func (slf *RotateLogs) Close() error {
	slf.closeCh <- struct{}{}
	return slf.file.Close()
}

func (slf *RotateLogs) handleEvent() {
	for {
		select {
		case <-slf.closeCh:
			return
		case now := <-slf.rotateCh:
			_ = slf.rotate(now)
		}
	}
}

func (slf *RotateLogs) rotate(now time.Time) error {
	if slf.rotateTime != 0 {
		duration := slf.untilNextTime(now, slf.rotateTime)
		slf.rotateCh = time.After(duration)
	}

	filePath := slf.latestFilePath(now)
	slf.mutex.Lock()
	defer slf.mutex.Unlock()
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if slf.file != nil {
		_ = slf.file.Close()
	}
	slf.file = file

	if len(slf.logLinkPath) > 0 {
		_ = os.Remove(slf.logLinkPath)
		_ = os.Link(filePath, slf.logLinkPath)
	}

	if slf.maxAge > 0 {
		go slf.deleteExpiredFile(now)
	}

	return nil
}

func (slf *RotateLogs) deleteExpiredFile(now time.Time) {
	cutoffTime := now.Add(-slf.maxAge)
	matches, err := filepath.Glob(slf.fileWildcard)
	if err != nil {
		return
	}

	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}
		if slf.maxAge > 0 && fileInfo.ModTime().After(cutoffTime) {
			continue
		}
		if fileInfo.Name() == filepath.Base(slf.logPath) {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	for _, path := range toUnlink {
		_ = os.Remove(path)
	}
}

func (slf *RotateLogs) latestFilePath(t time.Time) string {
	path := strings.Join([]string{slf.logPath, slf.logPathFormat}, ".")
	return t.Format(path)
}

func (slf *RotateLogs) untilNextTime(now time.Time, duration time.Duration) time.Duration {
	unixNano := now.UnixNano()
	nanoseconds := duration.Nanoseconds()
	next := nanoseconds - (unixNano % nanoseconds)
	return time.Duration(next)
}
