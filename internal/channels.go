package internal

import (
	"fmt"
	"sync"
)

var vChannels map[string]VoteChannels
var vLock sync.RWMutex
var vInit sync.Once

func init() {
	vInit.Do(func() {
		vLock.Lock()
		defer vLock.Unlock()
		vChannels = make(map[string]VoteChannels)
	})
}

type VoteChannels struct {
	VoteAgree chan uint
	CloseMsg  chan string
}

func GetVoteChannels(voteName string) (VoteChannels, error) {

	vLock.RLock()
	defer vLock.RUnlock()

	var ok bool
	var channels VoteChannels

	if channels, ok = vChannels[voteName]; !ok {
		return VoteChannels{}, fmt.Errorf("vote %s isn't exist", voteName)
	}

	return channels, nil
}

func CreateVoteChannels(voteName string) error {

	if _, err := GetVoteChannels(voteName); err == nil {
		return fmt.Errorf("vote %s has already exist ", voteName)
	}

	vLock.Lock()
	defer vLock.Unlock()

	vChannels[voteName] = VoteChannels{
		VoteAgree: make(chan uint),
		CloseMsg:  make(chan string),
	}

	return nil
}
