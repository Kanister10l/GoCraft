package templater

import (
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kanister10l/GoCraft/logger"

	"github.com/kanister10l/GoCraft/eventmanager"
	"github.com/ztrue/tracerr"
)

// Templater ...
type Templater struct {
	BaseLocation  string
	ScanPeriod    int
	InputFileList []string
	HashMap       map[string]string
	Hasher        hash.Hash
	ToProcess     []string
}

// NewTemplateEngine ...
func NewTemplateEngine(location string, period int) *Templater {
	t := &Templater{
		BaseLocation:  location,
		ScanPeriod:    period,
		InputFileList: []string{},
		HashMap:       make(map[string]string, 10),
		Hasher:        sha512.New(),
		ToProcess:     []string{},
	}

	go t.PeriodicalScan()

	return t
}

// PeriodicalScan ...
func (t *Templater) PeriodicalScan() {
	stop := make(chan bool)
	eventmanager.Master.Register("stop", func() {
		logger.Logger.Info("Stopping Template Engine")
		stop <- true
	})

	for {
		if err := t.CheckFilesForChanges(); err != nil {
			continue
		}

		delay := time.After(time.Second * time.Duration(t.ScanPeriod))
		logger.Logger.Info(fmt.Sprintf("Template Engine asleep for: %d seconds", t.ScanPeriod))
		select {
		case <-delay:
		case <-stop:
			return
		}
	}
}

// CheckFilesForChanges ...
func (t *Templater) CheckFilesForChanges() error {
	t.InputFileList = []string{}

	err := filepath.Walk(t.BaseLocation, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			t.InputFileList = append(t.InputFileList, path)
		}
		return nil
	})

	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return err
	}

	t.ToProcess = []string{}
	for _, v := range t.InputFileList {
		if t.HashFile(v) {
			t.ToProcess = append(t.ToProcess, v)
		}
	}

	logger.Logger.Infow("new files to process", "files", t.ToProcess)

	return nil
}

// HashFile ...
func (t *Templater) HashFile(filename string) bool {
	f, err := os.Open(filename)
	if err != nil {
		tracerr.PrintSourceColor(tracerr.Wrap(err))
		return false
	}
	defer f.Close()

	content, err := t.ReadFile(f)
	if err != nil {
		return false
	}

	t.Hasher.Reset()
	hash := string(t.Hasher.Sum(content))
	if v, ok := t.HashMap[filename]; ok && v == hash {
		return false
	}
	t.HashMap[filename] = hash

	return true
}

// ReadFile ...
func (t *Templater) ReadFile(f *os.File) ([]byte, error) {
	fileContent := []byte{}
	for {
		buffer := make([]byte, 8192)
		_, err := f.Read(buffer)
		if err == io.EOF {
			fileContent = append(fileContent, buffer...)
			break
		} else if err != nil {
			tracerr.PrintSourceColor(tracerr.Wrap(err))
			return nil, err
		}
		fileContent = append(fileContent, buffer...)
	}

	return fileContent, nil
}
