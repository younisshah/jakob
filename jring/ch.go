package jring

import (
	"github.com/younisshah/jakob/jfs"
	"stathat.com/c/consistent"
)

/**
*  Created by Galileo on 17/6/17.
 */

type jRing struct {
	pType jfs.Peer
}

// New returns a new JRing with the specified peer type
func New(peer jfs.Peer) *jRing {
	return &jRing{pType: peer}
}

// Get gets the peer on which the "id" is likely to be found
func (j *jRing) Get(id string) (string, error) {
	jyml := jfs.NewJYaml()
	if j.pType == jfs.SETTER {
		jyml.Type = jfs.SETTER
	} else {
		jyml.Type = jfs.GETTER
	}
	peers, err := jyml.Peers()
	if err != nil {
		return "", err
	}
	ring := consistent.New()
	for _, p := range peers {
		ring.Add(p)
	}
	return ring.Get(id)
}
