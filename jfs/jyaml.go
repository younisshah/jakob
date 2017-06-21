// Package jfs deals with initializing the setter and getter YAML files.
// It provides for creating, truncating, reading, appending peer YAML files.
package jfs

import (
	"path/filepath"

	"os"

	"bufio"

	"io/ioutil"

	"golang.org/x/sync/errgroup"
)

const (
	SETTER Peer = iota
	GETTER
	_base = "/jring/"
)

type Peer int

var (
	pwd, _         = os.Getwd()
	setterFilePath = filepath.Join(pwd, _base, "setters.yaml")
	getterFilePath = filepath.Join(pwd, _base, "getters.yaml")
)

type JYaml struct {
	Type    Peer
	Address string
}

func NewJYaml() *JYaml {
	return &JYaml{}
}

// Init initializes/resets the setters and getter YAML files
// to store the first setter and getter peer address.
func (j *JYaml) Init(setterAddress, getterAddress string) error {
	var g errgroup.Group

	g.Go(func() error {
		s := new(SetterPeers)
		return j.create(setterFilePath, setterAddress, s)
	})

	g.Go(func() error {
		g := new(GetterPeers)
		return j.create(getterFilePath, getterAddress, g)
	})

	if err := g.Wait(); err != nil {
		return err
	} else {
		return nil
	}
}

// Append appends to the YAML file
func (j *JYaml) Append() error {
	file, err := j.open(os.O_APPEND | os.O_WRONLY)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString("- " + j.Address)
	return err
}

// Peers returns the array of setter/getter peers
func (j *JYaml) Peers() ([]string, error) {
	file, err := j.open(os.O_RDONLY)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if j.Type == GETTER {
		peers, err := GetterPeers{}.Peers(string(b))
		if err != nil {
			return nil, err
		}
		return peers, nil
	}
	peers, err := SetterPeers{}.Peers(string(b))
	if err != nil {
		return nil, err
	}
	return peers, nil
}

// create creates/truncates the file at path and
// write a YAML array to the file
func (j *JYaml) create(path, data string, p PeerType) error {
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	peers, err := p.YAMLString([]string{data})
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)
	if _, err := w.WriteString(peers); err != nil {
		return err
	}
	w.Flush()
	return nil
}

// open opens a YAML file depending on the PeerType set - getters or setters file
func (j *JYaml) open(flag int) (*os.File, error) {
	if j.Type == SETTER {
		return os.OpenFile(setterFilePath, flag, 0600)
	}
	return os.OpenFile(getterFilePath, flag, 0600)
}
