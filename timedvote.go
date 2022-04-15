package timedvoting

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type TimedVote struct {
	Vote      IVote
	Duration  time.Duration
	FeedBacks IVoteFeedBacks

	locker       sync.Mutex
	agreeChannel chan uint
	duration     time.Duration
	isVoting     bool
}

func (tm *TimedVote) Start() error {

	tm.locker.Lock()
	defer tm.locker.Unlock()

	if tm.isVoting {
		return fmt.Errorf("object TimedVote err : %w", errors.New("vote has already start"))
	}

	if tm.FeedBacks == nil {
		return fmt.Errorf("object TimedVote err : %w", errors.New("feedbacks should not empty"))
	}

	tm.isVoting = true
	tm.agreeChannel = make(chan uint)
	tm.duration = tm.Duration
	go tm.votingHandler()

	return nil
}

func (tm *TimedVote) Agree(agreeNum uint) error {
	tm.locker.Lock()
	defer tm.locker.Unlock()

	if !tm.isVoting {
		return fmt.Errorf("object TimedVote err : %w", errors.New("vote is finished or not start yet"))
	}

	tm.agreeChannel <- agreeNum

	return nil

}

func (tm *TimedVote) Close() error {
	var isVoting bool

	tm.locker.Lock()
	isVoting = tm.isVoting
	tm.locker.Unlock()

	if !isVoting {
		return fmt.Errorf("vote is finished or not start yet")
	}

	tm.closeVoting()
	return nil
}

func (tm *TimedVote) IsVoting() bool {
	tm.locker.Lock()
	defer tm.locker.Unlock()
	return tm.isVoting
}

func (tm *TimedVote) votingHandler() {

	var timmer *time.Timer = time.AfterFunc(tm.duration, tm.closeVoting)

	for {
		if agreeNum, ok := <-tm.agreeChannel; ok {
			tm.Vote.Agree(agreeNum)
			if tm.Vote.IsPass() {
				tm.closeVoting()
			}
		} else {
			break
		}
	}

	if !timmer.Stop() {
		<-timmer.C
	}

	if tm.Vote.IsPass() {
		tm.FeedBacks.Pass()
	} else {
		tm.FeedBacks.UnPass()
	}

}

func (tm *TimedVote) closeVoting() {
	tm.locker.Lock()
	defer tm.locker.Unlock()
	if tm.isVoting {
		tm.isVoting = false
		close(tm.agreeChannel)
	}
}
