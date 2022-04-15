package timedvoting_test

import (
	"sync"
	"testing"

	"github.com/Lornzo/timedvoting"
)

func TestMutiAgreePass(t *testing.T) {
	var vote timedvoting.IVote = timedvoting.NewBasicVote(10)
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup, v timedvoting.IVote) {
			vote.Agree(1)
			w.Done()
		}(&wg, vote)
	}

	wg.Wait()

	if !vote.IsPass() {
		t.Errorf("vote should pass for agree 10000 times / passLimit 10")
	}

}

func TestMutiAgreeUnPass(t *testing.T) {
	var vote timedvoting.IVote = timedvoting.NewBasicVote(100)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(w *sync.WaitGroup, v timedvoting.IVote) {
			vote.Agree(1)
			w.Done()
		}(&wg, vote)
	}

	wg.Wait()

	if vote.IsPass() {
		t.Errorf("vote should unpass for agree 10 times / passLimit 100")
	}

}
