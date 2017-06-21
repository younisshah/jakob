package jfs

import (
	yaml "gopkg.in/yaml.v2"
)

/**
*  Created by Galileo on 18/6/17.
 */

// PeerType is the common interface implemented by SetterPeers and GetterPeers
type PeerType interface {
	// YAMLString returns the peers as YAML string
	YAMLString(peers []string) (string, error)
	// Peers returns the array of setter/getter peers
	Peers(yml string) ([]string, error)
}

// SetterPeers is the setter peers
type SetterPeers struct {
	Setters []string `yaml:"setters"`
}

func (s SetterPeers) YAMLString(peers []string) (string, error) {
	s.Setters = peers
	b, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func (s SetterPeers) Peers(yml string) ([]string, error) {
	var sp SetterPeers
	if err := yaml.Unmarshal([]byte(yml), &sp); err != nil {
		return nil, err
	}
	return sp.Setters, nil
}

// GetterPeers is the getter peers
type GetterPeers struct {
	Getters []string `yaml:"getters"`
}

func (s GetterPeers) YAMLString(peers []string) (string, error) {
	s.Getters = peers
	b, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func (g GetterPeers) Peers(yml string) ([]string, error) {
	var gp GetterPeers
	if err := yaml.Unmarshal([]byte(yml), &gp); err != nil {
		return nil, err
	}
	return gp.Getters, nil
}
