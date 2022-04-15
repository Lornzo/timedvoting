package timedvoting

import "sync"

func NewBasicVote(passLimit uint) *basicVote {
	return &basicVote{
		pass: passLimit,
	}
}

type basicVote struct {
	pass   uint
	agree  uint
	locker sync.RWMutex
}

func (bv *basicVote) Agree(agreeNum uint) {
	bv.locker.Lock()
	defer bv.locker.Unlock()
	bv.agree += agreeNum
}

func (bv *basicVote) IsPass() bool {
	bv.locker.RLock()
	defer bv.locker.RUnlock()
	return bv.agree >= bv.pass
}
