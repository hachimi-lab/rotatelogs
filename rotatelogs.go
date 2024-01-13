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
	logLinkPath   string
	logPathFormat string
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
	ins.logLinkPath = ins.logPath
	ins.logPathFormat = strings.Join([]string{ins.logPath, ins.rotateType.TimeFormat()}, ".")
	ins.fileWildcard = fmt.Sprintf("%s.*", ins.logPath)

	if err := os.Mkdir(filepath.Dir(ins.logPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if err := ins.rotate(time.Now().Local()); err != nil {
		return nil, err
	}

	if ins.rotateType.isValid() {
		go ins.listen()
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

func (slf *RotateLogs) listen() {
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
	if slf.rotateType.isValid() {
		duration := slf.rotateType.UntilNextTime(now)
		slf.rotateCh = time.After(duration)
	}

	filePath := now.Format(slf.logPathFormat)
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
		go slf.expire(now)
	}

	return nil
}

func (slf *RotateLogs) expire(now time.Time) {
	cutoffTime := now.Add(-slf.maxAge)
	matches, err := filepath.Glob(slf.fileWildcard)
	if err != nil {
		return
	}

	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		stat, err := os.Stat(path)
		if err != nil {
			continue
		}
		if slf.maxAge > 0 && stat.ModTime().After(cutoffTime) {
			continue
		}
		if stat.Name() == filepath.Base(slf.logPath) {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	for _, path := range toUnlink {
		_ = os.Remove(path)
	}
}
