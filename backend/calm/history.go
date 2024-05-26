package calm

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/anda-ai/anda/models"
	logger "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
)

type MessageHistory struct {
	dataDir  string
	maxSize  int
	userFile sync.Map
	locker   sync.Locker
}

func NewMessageHistory(dataDir string, maxSize int) *MessageHistory {
	return &MessageHistory{
		dataDir: dataDir,
		maxSize: maxSize,
	}
}

// iterator dir and find the max name of file
func maxFileByDir(dir string) (*os.File, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var maxName string
	for _, name := range names {
		if maxName < name {
			maxName = name
		}
	}

	return os.OpenFile(dir+"/"+maxName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (h *MessageHistory) AddMessage(ctx context.Context, user string, message *models.Message) error {
	uh, ok := h.userFile.Load(user)

	if !ok {
		h.locker.Lock()
		defer h.locker.Unlock()
		var err error

		uh, err = newUserHistory(h.dataDir, user)
		if err != nil {
			return err
		}
	}
	return uh.(*userHistory).Write(message)
}

type userHistory struct {
	user    string
	dataDir string
	file    *os.File
	locker  sync.RWMutex
}

func newUserHistory(dataDir, user string) (*userHistory, error) {

	file, err := maxFileByDir(dataDir + "/" + user)
	if err != nil {
		return nil, err
	}

	if file == nil {
		file, err = os.OpenFile(dataDir+"/"+user+"/"+time.Now().Format("20060102150405.000"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	return &userHistory{
		user:    user,
		dataDir: dataDir,
		file:    file,
	}, nil
}

func (h *userHistory) Write(msg *models.Message) error {
	h.locker.Lock()
	defer h.locker.Unlock()
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if _, err = h.file.Write([]byte{uint8('\n')}); err != nil {
		return err
	}

	if _, err := h.file.Write(data); err != nil {
		return err
	}

	return err
}

func (h *userHistory) fileSize() (int64, error) {
	fileInfo, err := h.file.Stat()
	if err != nil {
		return 0, fmt.Errorf("error getting file information: %v", err)
	}
	return fileInfo.Size(), nil
}

func (h *userHistory) Close() {
	h.locker.Lock()
	defer h.locker.Unlock()

	if err := h.file.Sync(); err != nil {
		logger.Errorf("user: %s sync file error: %v", h.user, err)
		return
	}

	err := h.file.Close()
	if err != nil {
		logger.Errorf("user: %s close file error: %v", h.user, err)
	}
}

func (h *userHistory) reset() error {
	h.locker.Lock()
	defer h.locker.Unlock()

	if err := h.file.Sync(); err != nil {
		return fmt.Errorf("user: %s sync file error: %v", h.user, err)
	}

	err := h.file.Close()
	if err != nil {
		return fmt.Errorf("user: %s close file error: %v", h.user, err)
	}

	file, err := os.OpenFile(h.dataDir+"/"+h.user+"/"+time.Now().Format("20060102150405.000"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("user: %s open file error: %v", h.user, err)
	}

	h.file = file

	return nil
}
