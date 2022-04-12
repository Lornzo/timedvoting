package timedvoting

type IVoteModel interface {
	AddVote(voteName string) error
}
