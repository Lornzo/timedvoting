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

	locker   sync.Mutex
	isVoting bool
	timmer   *time.Timer
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
	tm.timmer = time.AfterFunc(tm.Duration, tm.closeVoting)

	return nil
}

func (tm *TimedVote) Agree(agreeNum uint) error {

	if tm.IsVoting() {
		return fmt.Errorf("object TimedVote err : %w", errors.New("vote is finished or not start yet"))
	}
	fmt.Println("123")
	tm.Vote.Agree(agreeNum)
	if tm.Vote.IsPass() {
		fmt.Println("pass and close")
		tm.closeVoting()
	} else {
		fmt.Println("unpass and not close")
	}

	return nil
}

func (tm *TimedVote) Close() error {

	if !tm.IsVoting() {
		return fmt.Errorf("vote is finished or not start yet")
	}

	tm.closeVoting()
	return nil
}

func (tm *TimedVote) IsVoting() bool {
	tm.locker.Lock()
	defer tm.locker.Unlock()
	fmt.Println(tm.isVoting)
	return tm.isVoting
}

func (tm *TimedVote) closeVoting() {

	if !tm.IsVoting() {
		return
	}

	tm.locker.Lock()
	defer tm.locker.Unlock()

	tm.isVoting = false
	if !tm.timmer.Stop() {
		<-tm.timmer.C
	}

	if tm.Vote.IsPass() {
		fmt.Println("pass")
		tm.FeedBacks.Pass()
	} else {
		fmt.Println("unpass")
		tm.FeedBacks.UnPass()
	}
}
