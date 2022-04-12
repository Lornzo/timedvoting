package timedvoting

import (
	"errors"
	"sync"
	"time"

	"github.com/Lornzo/timedvoting/internal"
)

var model IVoteModel

func SetModel() {}

type VotePassFeedBack func()

type VoteFaildFeedBack func()

type CreateVoteParams struct {
	VoteName      string
	AgreeLimit    uint
	VoteExpired   time.Duration
	PassFeedBack  VotePassFeedBack
	FaildFeedBack VoteFaildFeedBack
}

type voteChannels struct {
	Agree chan uint
	Close chan string
}

var vChannels map[string]voteChannels = make(map[string]voteChannels, 0)
var vLock sync.RWMutex

func CreateVote(params CreateVoteParams) error {

	if model == nil {
		return errors.New("model is nil, please call SetModel() first")
	}

	if err := internal.CreateVoteChannels(params.VoteName); err != nil {
		return err
	}

	go timmerHandler(params.VoteName)

	return nil
}

func timmerHandler(voteName string) {

}

func voteHandler() {}
