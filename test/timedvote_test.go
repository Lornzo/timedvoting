package timedvoting_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Lornzo/timedvoting"
)

type testFeedBack struct {
	locker sync.Mutex
	status string
}

func (tfb *testFeedBack) getStatus() string {
	tfb.locker.Lock()
	defer tfb.locker.Unlock()
	return tfb.status
}

func (tfb *testFeedBack) setStatus(s string) {
	tfb.locker.Lock()
	defer tfb.locker.Unlock()
	tfb.status = s
}

func (tfb *testFeedBack) Pass() {
	tfb.setStatus("pass")
}

func (tfb *testFeedBack) UnPass() {
	tfb.setStatus("unpass")
}

func TestTimeOutUnpass(t *testing.T) {

	var feedbacks *testFeedBack = &testFeedBack{}
	var tVote timedvoting.TimedVote = timedvoting.TimedVote{
		Duration:  1 * time.Second,
		Vote:      timedvoting.NewBasicVote(10),
		FeedBacks: feedbacks,
	}

	tVote.Start()

	for i := 0; i < 5; i++ {
		go func() {
			tVote.Agree(1)
		}()
	}

	time.Sleep(2 * time.Second)

	if feedbacks.getStatus() == "pass" {
		t.Error("timedvote should not pass")
	}

}

func TestAgreePass(t *testing.T) {

	var feedbacks *testFeedBack = &testFeedBack{}
	var tVote timedvoting.TimedVote = timedvoting.TimedVote{
		Duration:  10 * time.Second,
		Vote:      timedvoting.NewBasicVote(10),
		FeedBacks: feedbacks,
	}

	tVote.Start()

	var wg sync.WaitGroup

	for i := 0; i < 11; i++ {
		wg.Add(1)

		go func() {
			tVote.Agree(1)
			wg.Done()
		}()

	}
	wg.Wait()

	if feedbacks.getStatus() != "pass" {
		t.Error("timedvote should be pass")
	}
}

func TestCloseUnPass(t *testing.T) {

	var feedbacks *testFeedBack = &testFeedBack{}
	var tVote timedvoting.TimedVote = timedvoting.TimedVote{
		Duration:  10 * time.Second,
		Vote:      timedvoting.NewBasicVote(10),
		FeedBacks: feedbacks,
	}

	tVote.Start()

	var wg sync.WaitGroup

	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup, v *timedvoting.TimedVote) {
			tVote.Agree(1)
			w.Done()
		}(&wg, &tVote)
	}
	wg.Wait()

	if err := tVote.Close(); err != nil {
		t.Error(err)
	}

	if feedbacks.getStatus() == "pass" {
		t.Error("timedvote should be pass")
	}
}

func TestClosePass(t *testing.T) {

}
