package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

const (
	latestMetricsSymlink = "latest"
)

type Persistence struct {
	DataDir string
}

func NewPersistence(dataDir string) (*Persistence, error) {
	if err := createDataDir(dataDir); err != nil {
		return nil, err
	}
	return &Persistence{DataDir: dataDir}, nil
}

func (p *Persistence) Save(message proto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	fileName := p.generateBinaryFileName()
	latestSymlinkFileName := p.generateLatestSymlinkFileName()
	if err := ioutil.WriteFile(fileName, data, 0600); err != nil {
		return fmt.Errorf("cannot write binary data to file: %w", err)
	}
	if err := linkToLatestSymlink(latestSymlinkFileName, fileName); err != nil {
		log.Fatal("cannot link to latest symlink: ", err)
		return err
	}

	return nil
}

func (p *Persistence) Read(message proto.Message) error {
	data, err := ioutil.ReadFile(p.generateLatestSymlinkFileName())
	if err != nil {
		return fmt.Errorf("cannot read binary data from file: %w", err)
	}

	if err = proto.Unmarshal(data, message); err != nil {
		return fmt.Errorf("cannot unmarshal binary to proto message: %w", err)
	}

	return nil
}

func (p *Persistence) generateBinaryFileName() string {
	return fmt.Sprintf("%s/%d.bin", p.DataDir, time.Now().Unix())
}

func (p *Persistence) generateLatestSymlinkFileName() string {
	return fmt.Sprintf("%s/%s", p.DataDir, latestMetricsSymlink)
}

func linkToLatestSymlink(latestSymlinkFileName string, fileName string) error {
	if _, err := os.Lstat(latestSymlinkFileName); err == nil {
		if err := os.Remove(latestSymlinkFileName); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	if err := os.Symlink(fileName, latestSymlinkFileName); err != nil {
		return fmt.Errorf("cannot create latest symlink: %w", err)
	}
	return nil
}

func createDataDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
