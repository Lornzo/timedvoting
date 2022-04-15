package timedvoting

type IVote interface {
	Agree(agreeNum uint)
	IsPass() bool
}
